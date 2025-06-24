package idutil

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode"
)

// TestUUID 测试UUID生成功能
func TestUUID(t *testing.T) {
	// 测试基本生成功能
	uuid, err := UUID()
	if err != nil {
		t.Fatalf("UUID() failed: %v", err)
	}
	if len(uuid) != 36 {
		t.Errorf("UUID length should be 36, got %d", len(uuid))
	}
	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		t.Errorf("UUID format incorrect: %s", uuid)
	}

	// 测试并发安全性
	const concurrency = 100000
	results := make(chan string, concurrency)
	errChan := make(chan error, 1)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			uid, err := UUID()
			if err != nil {
				select {
				case errChan <- err:
				default:
				}
				return
			}
			results <- uid
		}()
	}

	// 等待所有goroutine完成并关闭通道
	go func() {
		wg.Wait()
		close(results)
	}()

	// 检查是否有错误
	select {
	case err := <-errChan:
		t.Fatalf("Concurrent UUID generation failed: %v", err)
	default:
	}

	// 检查唯一性
	seen := make(map[string]bool)
	for uid := range results {
		if seen[uid] {
			t.Error("Duplicate UUID generated")
		}
		seen[uid] = true
	}
}

// TestObjectID 测试ObjectID生成功能
func TestObjectID(t *testing.T) {
	// 测试基本生成功能
	oid := ObjectID()
	if len(oid) != 24 {
		t.Errorf("ObjectID length should be 24, got %d", len(oid))
	}

	// 解析时间戳部分并验证
	timestampHex := oid[:8]
	var timestamp uint32
	_, err := fmt.Sscanf(timestampHex, "%x", &timestamp)
	if err != nil {
		t.Fatalf("Failed to parse timestamp from ObjectID: %v", err)
	}

	// 时间戳不应早于当前时间10秒前
	currentTime := uint32(time.Now().Unix())
	if timestamp < currentTime-10 || timestamp > currentTime+1 {
		t.Errorf("ObjectID timestamp is not in valid range: %d (current: %d)", timestamp, currentTime)
	}

	// 测试并发唯一性
	const concurrency = 100000
	results := make(chan string, concurrency)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- ObjectID()
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	seen := make(map[string]bool)
	for id := range results {
		if seen[id] {
			t.Error("Duplicate ObjectID generated")
		}
		seen[id] = true
	}
}

// TestSnowflakeGenerator 测试雪花算法生成器
func TestSnowflakeGenerator(t *testing.T) {
	// 测试创建生成器
	generator, err := NewSnowflakeGenerator(1, 1)
	if err != nil {
		t.Fatalf("NewSnowflakeGenerator failed: %v", err)
	}

	// 测试基本生成功能
	id, err := generator.NextID()
	if err != nil {
		t.Fatalf("NextID failed: %v", err)
	}
	if id <= 0 {
		t.Errorf("Generated invalid Snowflake ID: %d", id)
	}

	// 测试ID结构
	extractedTimestamp := (id >> timestampShift) + snowflakeEpoch
	currentTimestamp := time.Now().UnixMilli()
	if extractedTimestamp < currentTimestamp-100 || extractedTimestamp > currentTimestamp+100 {
		t.Errorf("Snowflake ID timestamp is not in valid range: %d (current: %d)", extractedTimestamp, currentTimestamp)
	}

	extractedWorkerID := (id >> workerIDShift) & (1<<workerIDBits - 1)
	if extractedWorkerID != 1 {
		t.Errorf("Extracted workerID incorrect: got %d, want 1", extractedWorkerID)
	}

	extractedProcessID := (id >> processIDShift) & (1<<processIDBits - 1)
	if extractedProcessID != 1 {
		t.Errorf("Extracted processID incorrect: got %d, want 1", extractedProcessID)
	}

	// 测试并发生成
	const concurrency = 100000
	results := make(chan int64, concurrency)
	errChan := make(chan error, 1)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := generator.NextID()
			if err != nil {
				select {
				case errChan <- err:
				default:
				}
				return
			}
			results <- id
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	select {
	case err := <-errChan:
		t.Fatalf("Concurrent Snowflake generation failed: %v", err)
	default:
	}

	// 检查唯一性和递增性
	var ids []int64
	for id := range results {
		ids = append(ids, id)
	}

	// 检查唯一性
	seen := make(map[int64]bool)
	for _, id := range ids {
		if seen[id] {
			t.Error("Duplicate Snowflake ID generated")
		}
		seen[id] = true
	}

	// 检查递增性
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	for i := 1; i < len(ids); i++ {
		if ids[i] <= ids[i-1] {
			t.Errorf("Snowflake IDs are not increasing: prev %d, current %d", ids[i-1], ids[i])
		}
	}
}

// TestSnowflakeGenerator_InvalidParams 测试无效参数
func TestSnowflakeGenerator_InvalidParams(t *testing.T) {
	// 测试无效workerID
	_, err := NewSnowflakeGenerator(-1, 1)
	if err == nil {
		t.Error("Expected error for negative workerID, got nil")
	}

	_, err = NewSnowflakeGenerator(maxWorkerID+1, 1)
	if err == nil {
		t.Error("Expected error for workerID exceeding max, got nil")
	}

	// 测试无效processID
	_, err = NewSnowflakeGenerator(1, -1)
	if err == nil {
		t.Error("Expected error for negative processID, got nil")
	}

	_, err = NewSnowflakeGenerator(1, maxProcessID+1)
	if err == nil {
		t.Error("Expected error for processID exceeding max, got nil")
	}
}

// TestULID 测试ULID生成功能
func TestULID(t *testing.T) {
	// 测试基本功能
	ulid, err := ULID()
	if err != nil {
		t.Fatalf("ULID生成失败: %v", err)
	}
	if len(ulid) != 26 {
		t.Errorf("ULID长度应为26，实际为%d", len(ulid))
	}
	// 验证字符集
	for _, c := range ulid {
		if !unicode.IsUpper(c) && !unicode.IsDigit(c) {
			t.Errorf("ULID包含无效字符: %c", c)
		}
		// 排除Crockford Base32中不包含的字符
		if c == 'I' || c == 'L' || c == 'O' || c == 'U' {
			t.Errorf("ULID包含歧义字符: %c", c)
		}
	}

	// 测试排序性
	var prevULID string
	for i := 0; i < 100; i++ {
		currentULID, err := ULID()
		if err != nil {
			t.Fatalf("ULID生成失败: %v", err)
		}
		if i > 0 && currentULID <= prevULID {
			t.Errorf("ULID未按预期排序: 前一个=%s, 当前=%s", prevULID, currentULID)
		}
		prevULID = currentULID
		time.Sleep(1 * time.Millisecond) // 确保时间戳递增
	}
}

// TestULID_Concurrency 测试ULID并发生成唯一性
func TestULID_Concurrency(t *testing.T) {
	const concurrency = 1000
	const idsPerGoroutine = 100
	idChan := make(chan string, concurrency*idsPerGoroutine)
	errChan := make(chan error, concurrency)
	var wg sync.WaitGroup

	// 启动多个goroutine并发生成ULID
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				ulid, err := ULID()
				if err != nil {
					errChan <- fmt.Errorf("ULID生成失败: %w", err)
					return
				}
				idChan <- ulid
			}
		}()
	}

	// 等待所有goroutine完成
	go func() {
		wg.Wait()
		close(idChan)
		close(errChan)
	}()

	// 检查错误
	for err := range errChan {
		t.Error(err)
	}

	// 检查唯一性
	idSet := make(map[string]bool)
	count := 0
	for ulid := range idChan {
		count++
		if idSet[ulid] {
			t.Errorf("发现重复的ULID: %s", ulid)
		}
		idSet[ulid] = true
	}

	if count != concurrency*idsPerGoroutine {
		t.Errorf("生成的ULID数量不正确: 预期=%d, 实际=%d", concurrency*idsPerGoroutine, count)
	}
}

// TestNanoID 测试NanoID生成功能
func TestNanoID(t *testing.T) {
	// 测试默认参数
	defaultID, err := DefaultNanoID()
	if err != nil {
		t.Fatalf("默认NanoID生成失败: %v", err)
	}
	if len(defaultID) != 21 {
		t.Errorf("默认NanoID长度应为21，实际为%d", len(defaultID))
	}

	// 测试自定义长度
	customLenID, err := NanoID(10, "")
	if err != nil {
		t.Fatalf("自定义长度NanoID生成失败: %v", err)
	}
	if len(customLenID) != 10 {
		t.Errorf("自定义长度NanoID应为10，实际为%d", len(customLenID))
	}

	// 测试自定义字符集
	customAlphabet := "abc123"
	customID, err := NanoID(15, customAlphabet)
	if err != nil {
		t.Fatalf("自定义字符集NanoID生成失败: %v", err)
	}
	for _, c := range customID {
		if !strings.ContainsRune(customAlphabet, c) {
			t.Errorf("NanoID包含自定义字符集外的字符: %c", c)
		}
	}

	// 测试无效参数
	_, err = NanoID(0, "") // 无效长度，应使用默认值
	if err != nil {
		t.Error("无效长度测试失败:", err)
	}

	_, err = NanoID(10, "a") // 字符集长度不足
	if err == nil {
		t.Error("预期字符集长度不足错误，但未收到")
	}

	_, err = NanoID(10, strings.Repeat("a", 256)) // 字符集过长
	if err == nil {
		t.Error("预期字符集过长错误，但未收到")
	}
}

// TestNanoID_Concurrency 测试NanoID并发生成唯一性
func TestNanoID_Concurrency(t *testing.T) {
	const concurrency = 1000
	const idsPerGoroutine = 100
	idChan := make(chan string, concurrency*idsPerGoroutine)
	errChan := make(chan error, concurrency)
	var wg sync.WaitGroup

	// 启动多个goroutine并发生成NanoID
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				id, err := DefaultNanoID()
				if err != nil {
					errChan <- fmt.Errorf("NanoID生成失败: %w", err)
					return
				}
				idChan <- id
			}
		}()
	}

	// 等待所有goroutine完成
	go func() {
		wg.Wait()
		close(idChan)
		close(errChan)
	}()

	// 检查错误
	for err := range errChan {
		t.Error(err)
	}

	// 检查唯一性
	idSet := make(map[string]bool)
	count := 0
	for id := range idChan {
		count++
		if idSet[id] {
			t.Errorf("发现重复的NanoID: %s", id)
		}
		idSet[id] = true
	}

	if count != concurrency*idsPerGoroutine {
		t.Errorf("生成的NanoID数量不正确: 预期=%d, 实际=%d", concurrency*idsPerGoroutine, count)
	}
}
