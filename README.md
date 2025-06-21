# Golang工具类集合

一个实用的Golang工具类库，包含日期时间处理和字符串操作等常用功能模块。

## 功能模块

- **dateutil**: 日期时间处理工具，支持周计算、日期范围生成、月末判断等功能
- **strutil**: 字符串处理工具，包含字符串截取、替换、判断等功能

## 安装

```bash
go get github.com/luckxgo/go-utils
```

## 使用示例

### 日期范围生成

```go
package main

import (
  "fmt"
  "time"
  "github.com/luckxgo/go-utils/dateutil"
)

func main() {
  start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
  end := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
  
  // 生成每日日期范围
  dr := dateutil.Range(start, end, dateutil.DayUnit)
  dates := dr.Generate()
  
  for _, d := range dates {
    fmt.Println(d.Format("2006-01-02"))
  }
}
```

### 周结束日期计算

```go
// 获取指定日期的周结束日期（周一为一周开始）
date := time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC)
end := dateutil.EndOfWeek(date, false)
fmt.Println(end.Format("2006-01-02 15:04:05")) // 输出: 2023-10-15 23:59:59
```

## 测试

```bash
# 运行所有测试
cd dateutil
go test -v

# 快速失败模式
go test -failfast
```

## 许可证

[MIT](LICENSE)