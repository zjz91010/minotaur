package survey

import "time"

// Option 选项
type Option func(logger *logger)

// WithLayout 设置日志文件名的时间戳格式
//   - 默认为 time.DateOnly
func WithLayout(layout string) Option {
	return func(logger *logger) {
		logger.layout = layout
		logger.layoutLen = len(layout)
	}
}

// WithFlushInterval 设置日志文件刷新间隔
//   - 默认为 3s，当日志文件刷新间隔 <= 0 时，将会在每次写入日志时刷新日志文件
func WithFlushInterval(interval time.Duration) Option {
	return func(logger *logger) {
		logger.interval = interval
	}
}
