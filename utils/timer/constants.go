package timer

import "time"

const (
	timingWheelTick = time.Millisecond * 10 // 时间轮一刻度长度
)

const (
	Forever   = -1 //  无限循环
	Once      = 1  // 一次
	Instantly = 0  // 立刻
)
