package idutil

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	mathrand "math/rand"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// GenerateUUID 生成UUID v4 (RFC 4122)
// 返回标准UUID字符串(36字符)和可能的错误
func UUID() (string, error) {
	var uuid [16]byte

	// 生成16字节随机数
	if _, err := rand.Read(uuid[:]); err != nil {
		return "", fmt.Errorf("UUID生成失败: %w", err)
	}

	// 设置UUID版本和变体
	uuid[6] = (uuid[6] & 0x0F) | 0x40 // 版本4 (随机)
	uuid[8] = (uuid[8] & 0x3F) | 0x80 // RFC 4122变体

	// 格式化UUID字符串
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// ObjectId相关变量与初始化
var (
	objectIDCounter  uint32
	machineID        [3]byte
	processID        [2]byte
	timestampCounter uint32
	lastTimestampSec int64
	randPool         = sync.Pool{
		New: func() interface{} {
			return mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
		},
	}
	randMutex sync.Mutex
)

func init() {
	// 初始化机器ID(优先使用MAC地址)
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range ifaces {
			if iface.Flags&net.FlagUp != 0 && len(iface.HardwareAddr) >= 6 {
				copy(machineID[:], iface.HardwareAddr[:3])
				break
			}
		}
	}
	// 若无法获取MAC地址则使用随机数
	if machineID == [3]byte{0, 0, 0} {
		_, _ = rand.Read(machineID[:])
	}

	// 初始化进程ID
	pid := os.Getpid()
	processID[0] = byte(pid >> 8)
	processID[1] = byte(pid)
}

// GenerateObjectID 生成MongoDB风格的ObjectId(24字符十六进制字符串)
// 结构: 4字节时间戳 + 3字节机器ID + 2字节进程ID + 3字节计数器
func ObjectID() string {
	b := make([]byte, 12)

	// 时间戳(秒级)
	binary.BigEndian.PutUint32(b[:4], uint32(time.Now().Unix()))
	// 机器ID
	copy(b[4:7], machineID[:])
	// 进程ID
	copy(b[7:9], processID[:])
	// 计数器(原子递增)
	counter := atomic.AddUint32(&objectIDCounter, 1)
	b[9] = byte(counter >> 16)
	b[10] = byte(counter >> 8)
	b[11] = byte(counter)

	return fmt.Sprintf("%x", b)
}

// Snowflake算法实现
const (
	snowflakeEpoch = 1609459200000 // 起始时间戳(2021-01-01 00:00:00 UTC)
	workerIDBits   = 5             // 机器ID位数
	processIDBits  = 5             // 进程ID位数
	sequenceBits   = 12            // 序列号位数

	maxWorkerID  = -1 ^ (-1 << workerIDBits)  // 最大机器ID(31)
	maxProcessID = -1 ^ (-1 << processIDBits) // 最大进程ID(31)
	maxSequence  = -1 ^ (-1 << sequenceBits)  // 最大序列号(4095)

	workerIDShift  = sequenceBits
	processIDShift = sequenceBits + workerIDBits
	timestampShift = sequenceBits + workerIDBits + processIDBits
)

// SnowflakeGenerator 雪花算法生成器
type SnowflakeGenerator struct {
	workerID      int64      // 机器ID(0-31)
	processID     int64      // 进程ID(0-31)
	lastTimestamp int64      // 上次生成ID的时间戳
	sequence      int64      // 当前序列号(0-4095)
	mu            sync.Mutex // 互斥锁，确保并发安全
}

// NewSnowflakeGenerator 创建雪花算法生成器
// workerID: 机器ID(0-31), processID: 进程ID(0-31)
func NewSnowflakeGenerator(workerID, processID int64) (*SnowflakeGenerator, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("workerID必须在[0, %d]范围内", maxWorkerID)
	}
	if processID < 0 || processID > maxProcessID {
		return nil, fmt.Errorf("processID必须在[0, %d]范围内", maxProcessID)
	}

	return &SnowflakeGenerator{
		workerID:      workerID,
		processID:     processID,
		lastTimestamp: 0,
		sequence:      0,
	}, nil
}

// NextID 生成下一个雪花ID
func (g *SnowflakeGenerator) NextID() (int64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for {
		// 获取当前时间戳(毫秒级)
		timestamp := time.Now().UnixMilli() - snowflakeEpoch

		// 处理时钟回拨
		if timestamp < g.lastTimestamp {
			return 0, errors.New("时钟回拨，无法生成ID")
		}

		// 同一毫秒内序列号递增
		if timestamp == g.lastTimestamp {
			g.sequence = (g.sequence + 1) & maxSequence
			// 序列号溢出，等待下一毫秒
			if g.sequence == 0 {
				// 等待直到时间戳递增
				for timestamp <= g.lastTimestamp {
					time.Sleep(time.Microsecond * 100)
					timestamp = time.Now().UnixMilli() - snowflakeEpoch
				}
				continue // 重新获取时间戳
			}
		} else {
			// 不同毫秒重置序列号
			g.sequence = 0
		}

		g.lastTimestamp = timestamp

		// 组合ID: 时间戳(41位) + 机器ID(5位) + 进程ID(5位) + 序列号(12位)
		return (timestamp<<timestampShift |
			g.workerID<<workerIDShift |
			g.processID<<processIDShift |
			g.sequence), nil
	}
}

// ULID 生成器接口
// ULID (Universally Unique Lexicographically Sortable Identifier) 是一种可排序的唯一标识符
// 格式: 128位 (16字节)，其中48位为时间戳(毫秒级)，80位为随机数
// 编码后为26个字符的Crockford Base32字符串

type ULIDGenerator struct {
	mu       sync.Mutex
	lastTime uint64
	random   [10]byte // 80位随机数
}

var defaultULIDGenerator = &ULIDGenerator{}

// NewULIDGenerator 创建新的ULID生成器
func NewULIDGenerator() *ULIDGenerator {
	return &ULIDGenerator{}
}

// ULID 生成一个新的ULID字符串
func (u *ULIDGenerator) ULID() (string, error) {
	now := uint64(time.Now().UnixMilli())
	u.mu.Lock()
	defer u.mu.Unlock()

	// 如果当前时间与上次相同，增加随机数
	if now == u.lastTime {
		// 递增随机数(大端序)
		for i := 9; i >= 0; i-- {
			u.random[i]++
			if u.random[i] != 0 {
				break
			}
			// 如果所有字节都溢出，则需要等待下一毫秒
			if i == 0 {
				time.Sleep(time.Millisecond - time.Duration(time.Now().UnixNano()%1e6)*time.Nanosecond)
				now = uint64(time.Now().UnixMilli())
				break
			}
		}
	} else {
		// 生成新的随机数
		if _, err := rand.Read(u.random[:]); err != nil {
			return "", fmt.Errorf("生成随机数失败: %w", err)
		}
	}

	u.lastTime = now

	// 组合ULID字节: 48位时间戳 + 80位随机数
	var ulidBytes [16]byte
	binary.BigEndian.PutUint64(ulidBytes[:8], now<<16) // 48位时间戳(左移16位对齐64位)
	copy(ulidBytes[6:], u.random[:])                   // 复制80位随机数到后10字节

	// 编码为Crockford Base32
	return encodeBase32(ulidBytes[:]), nil
}

// ULID 生成一个新的ULID字符串(使用默认生成器)
func ULID() (string, error) {
	return defaultULIDGenerator.ULID()
}

// base32编码表 (Crockford Base32)
const base32Alphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// encodeBase32 将字节数组编码为Crockford Base32字符串
func encodeBase32(data []byte) string {
	result := make([]byte, 0, 26)
	var buffer uint64
	bits := 0

	for _, b := range data {
		buffer = (buffer << 8) | uint64(b)
		bits += 8

		// 每次输出5位
		for bits >= 5 {
			bits -= 5
			result = append(result, base32Alphabet[(buffer>>bits)&0x1F])
		}
	}

	// 处理剩余的位
	if bits > 0 {
		buffer <<= (5 - bits)
		result = append(result, base32Alphabet[buffer&0x1F])
	}

	return string(result)
}

// NanoID 生成一个安全、紧凑、URL友好的唯一标识符
// length: ID长度，建议范围6-22，默认21
// alphabet: 自定义字符集，默认为"_-.0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
func NanoID(length int, alphabet string) (string, error) {
	if length <= 0 {
		length = 21
	}
	if alphabet == "" {
		alphabet = "_-.0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	alphabetLen := len(alphabet)
	if alphabetLen < 2 || alphabetLen > 255 {
		return "", errors.New("字符集长度必须在2-255之间")
	}

	// 计算每个随机字节需要的掩码
	logVal := math.Log2(float64(alphabetLen - 1))
	mask := (2 << int(math.Floor(logVal))) - 1
	// 计算预生成随机字节数
	step := int(math.Ceil(1.6 * float64(mask*length) / float64(alphabetLen)))

	result := make([]byte, 0, length)
	for {
		// 生成随机字节
		randomBytes := make([]byte, step)
		if _, err := rand.Read(randomBytes); err != nil {
			return "", fmt.Errorf("生成随机数失败: %w", err)
		}

		// 将随机字节转换为字符集中的字符
		for _, b := range randomBytes {
			idx := int(b) & mask
			if idx < alphabetLen {
				result = append(result, alphabet[idx])
				if len(result) == length {
					return string(result), nil
				}
			}
		}
	}
}

// DefaultNanoID 生成默认长度(21)的NanoID
func DefaultNanoID() (string, error) {
	return NanoID(21, "")
}
