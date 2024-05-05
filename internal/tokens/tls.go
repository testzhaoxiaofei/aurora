package tokens

import (
	"fmt"
	"sync"
	"time"
)

// 定义全局变量
var (
	counter     int        // 计数器
	startTime   time.Time  // 开始时间
	counterLock sync.Mutex // 互斥锁
)

// 初始化函数，设置开始时间和计数器初始值
func init() {
	startTime = time.Now()
	counter = 0
}

// AddCount  增加计数并检查是否达到条件
func AddCount() {
	counterLock.Lock() // 加锁以保护共享资源
	defer counterLock.Unlock()

	// 检查是否超过5分钟
	if time.Since(startTime) > 5*time.Minute {
		// 重置计数器和开始时间
		counter = 1
		Pingx = ""
		startTime = time.Now()
	} else {
		// 增加计数
		counter++
		fmt.Println("开始计数", counter)
		// 检查计数器是否达到100 次报错
		if counter >= 100 {
			// 触发相关函数
			fmt.Println("达到条件开始切换", counter)
			Pingx = "1"
			counter = 0
			startTime = time.Now()
		}
	}
}
