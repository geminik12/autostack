/**FileHeader
 * @Author: Liangkang Zhang
 * @Date: 2026/1/24 00:31:44
 * @LastEditors: Liangkang Zhang
 * @LastEditTime: 2026/1/25 00:25:00
 * @Description:
 * @Copyright: Copyright (©)}) 2026 Liangkang Zhang<lkzhang98@gmail.com>. All rights reserved. Use of this source code is governed by a MIT style license that can be found in the LICENSE file.. All rights reserved.
 * @Email: lkzhang98@gmail.com
 * @Repository: https://github.com/geminik12/autostack
 */
// Package log provides a simple logger.

package log

import (
	"context"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	gormlogger "gorm.io/gorm/logger"
)

// Logger defines the interface of a logger.
type Logger interface {
	Debugf(format string, args ...any)
	Debugw(msg string, keyvals ...any)
	Infof(format string, args ...any)
	Infow(msg string, keyvals ...any)
	Warnf(format string, args ...any)
	Warnw(msg string, keyvals ...any)
	Errorf(format string, args ...any)
	Errorw(err error, msg string, keyvals ...any)
	Panicf(format string, args ...any)
	Panicw(msg string, keyvals ...any)
	Fatalf(format string, args ...any)
	Fatalw(msg string, keyvals ...any)
	SetLevel(level string)
	W(ctx context.Context) Logger
	AddCallerSkip(skip int) Logger
	Sync()

	// integrate other loggers
	gormlogger.Interface
}

// logger is a wrapper of zap.Logger
type logger struct {
	z                 *zap.Logger
	opts              *Options
	atomicLevel       zap.AtomicLevel
	infoSync          *zapcore.BufferedWriteSyncer
	errorSync         *zapcore.BufferedWriteSyncer
	contextExtractors map[string]func(ctx context.Context) string
}

type Option func(*logger)

var _ Logger = (*logger)(nil)

var (
	mu  sync.Mutex
	std = NewLogger(NewOptions())
)

func Init(opts *Options, options ...Option) {
	mu.Lock()
	defer mu.Unlock()
	std = NewLogger(opts)
}

// NewLogger 创建一个新的Logger对象.
func NewLogger(opts *Options, options ...Option) *logger {
	if opts == nil {
		opts = NewOptions()
	}

	// 将文本格式的日志级别，例如 info 转换为 zapcore.Level 类型以供后面使用
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		// 如果指定了非法的日志级别，则默认使用 info 级别
		zapLevel = zapcore.InfoLevel
	}
	atomicLevel := zap.NewAtomicLevelAt(zapLevel)

	// 创建一个默认的 encoder 配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 自定义 MessageKey 为 message，message 语义更明确
	encoderConfig.MessageKey = "message"
	// 自定义 TimeKey 为 timestamp，timestamp 语义更明确
	encoderConfig.TimeKey = "timestamp"
	// 指定时间序列化函数，将时间序列化为 `2006-01-02 15:04:05.000` 格式，更易读
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	// 指定 time.Duration 序列化函数，将 time.Duration 序列化为经过的毫秒数的浮点数
	// 毫秒数比默认的秒数更精确
	encoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendFloat64(float64(d) / float64(time.Millisecond))
	}
	// when output to local path, with color is forbidden
	if opts.Format == "console" && opts.EnableColor {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	outputPaths := opts.OutputPaths
	if len(outputPaths) == 0 {
		outputPaths = []string{"stdout"}
	}

	// 创建构建 zap.Logger 需要的配置
	cfg := &zap.Config{
		// 是否在日志中显示调用日志所在的文件和行号，例如：`"caller":"onex/onex.go:75"`
		DisableCaller: opts.DisableCaller,
		// 是否禁止在 panic 及以上级别打印堆栈信息
		DisableStacktrace: opts.DisableStacktrace,
		// 指定日志级别
		Level: atomicLevel,
		// 指定日志显示格式，可选值：console, json
		Encoding:      opts.Format,
		EncoderConfig: encoderConfig,
		// 指定日志输出位置
		OutputPaths: outputPaths,
		// 设置 zap 内部错误输出位置
		ErrorOutputPaths: []string{"stderr"},
	}

	// 使用 cfg 创建 *zap.Logger 对象
	z, err := cfg.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(2))
	if err != nil {
		panic(err)
	}

	logger := &logger{atomicLevel: atomicLevel}

	// 如果开启了文件日志，则添加文件日志输出
	if opts.EnableFile {
		// 获取 encoder
		var encoder zapcore.Encoder
		if opts.Format == "json" {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		}

		// Info 日志 writer
		infoWriter := &lumberjack.Logger{
			Filename:   filepath.Join(opts.LogDir, "info.log"),
			MaxSize:    opts.MaxSize,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAge,
			Compress:   opts.Compress,
		}

		// Error 日志 writer
		errorWriter := &lumberjack.Logger{
			Filename:   filepath.Join(opts.LogDir, "error.log"),
			MaxSize:    opts.MaxSize,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAge,
			Compress:   opts.Compress,
		}

		// 使用 BufferedWriteSyncer 实现异步落盘
		logger.infoSync = &zapcore.BufferedWriteSyncer{
			WS:   zapcore.AddSync(infoWriter),
			Size: 256 * 1024, // 256KB 缓存
		}
		logger.errorSync = &zapcore.BufferedWriteSyncer{
			WS:   zapcore.AddSync(errorWriter),
			Size: 256 * 1024, // 256KB 缓存
		}

		// Info 核心：记录 Info 及以上级别日志（但不包含 Error 及以上，如果希望分开的话。通常 Info 包含所有，这里为了演示“分文件”，我们让 info.log 包含所有，error.log 包含 error）
		// 或者：info.log 只包含 Info 和 Warn，error.log 包含 Error, Panic, Fatal
		// 用户的需求是“按照日志级别分文件”，通常意味着不同级别去不同文件。
		// 这里实现：info.log 记录 Info及以上，error.log 记录 Error及以上。这样 error 也会出现在 info.log 中，这是最常见的做法。
		// 如果想要严格分离，可以使用 LevelEnablerFunc。
		// 这里采用：InfoWriter 记录 >= InfoLevel, ErrorWriter 记录 >= ErrorLevel.
		// 重新定义 infoCore 规则：使用 atomicLevel.Level() 保证动态生效
		infoCore := zapcore.NewCore(encoder, logger.infoSync, zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= atomicLevel.Level() && l < zapcore.ErrorLevel
		}))

		errorCore := zapcore.NewCore(encoder, logger.errorSync, zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= atomicLevel.Level() && l >= zapcore.ErrorLevel
		}))

		// 将文件日志核心与原有的核心（stdout）结合
		core := zapcore.NewTee(z.Core(), infoCore, errorCore)

		// 定时刷新异步缓存
		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for range ticker.C {
				_ = logger.infoSync.Sync()
				_ = logger.errorSync.Sync()
			}
		}()

		// 使用新的核心重建 Logger
		z = z.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return core
		}))
	}

	logger.z = z
	logger.opts = opts
	logger.contextExtractors = make(map[string]func(context.Context) string)

	// 应用所有传入的 Option
	for _, opt := range options {
		opt(logger)
	}

	return logger
}

// Default 返回全局 Logger.
func Default() Logger {
	return std
}

func SetLevel(level string) { std.SetLevel(level) }

func (l *logger) SetLevel(level string) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	l.atomicLevel.SetLevel(zapLevel)
}

func (l *logger) Sync() {
	_ = l.z.Sync()
	if l.infoSync != nil {
		_ = l.infoSync.Sync()
	}
	if l.errorSync != nil {
		_ = l.errorSync.Sync()
	}
}

func (l *logger) Options() *Options {
	return l.opts
}

func Debugf(format string, args ...any)            { std.Debugf(format, args...) }
func Debugw(msg string, keyvals ...any)            { std.Debugw(msg, keyvals...) }
func Infof(format string, args ...any)             { std.Infof(format, args...) }
func Infow(msg string, keyvals ...any)             { std.Infow(msg, keyvals...) }
func Warnf(format string, args ...any)             { std.Warnf(format, args...) }
func Warnw(msg string, keyvals ...any)             { std.Warnw(msg, keyvals...) }
func Errorf(format string, args ...any)            { std.Errorf(format, args...) }
func Errorw(err error, msg string, keyvals ...any) { std.Errorw(err, msg, keyvals...) }
func Panicf(format string, args ...any)            { std.Panicf(format, args...) }
func Panicw(msg string, keyvals ...any)            { std.Panicw(msg, keyvals...) }
func Fatalf(format string, args ...any)            { std.Fatalf(format, args...) }
func Fatalw(msg string, keyvals ...any)            { std.Fatalw(msg, keyvals...) }

func (l *logger) Debugf(format string, args ...any) { l.logf(zapcore.DebugLevel, format, args...) }
func (l *logger) Debugw(msg string, keyvals ...any) { l.logw(zapcore.DebugLevel, msg, keyvals...) }
func (l *logger) Infof(format string, args ...any)  { l.logf(zapcore.InfoLevel, format, args...) }
func (l *logger) Infow(msg string, keyvals ...any)  { l.logw(zapcore.InfoLevel, msg, keyvals...) }
func (l *logger) Warnf(format string, args ...any)  { l.logf(zapcore.WarnLevel, format, args...) }
func (l *logger) Warnw(msg string, keyvals ...any)  { l.logw(zapcore.WarnLevel, msg, keyvals...) }
func (l *logger) Errorf(format string, args ...any) { l.logf(zapcore.ErrorLevel, format, args...) }
func (l *logger) Errorw(err error, msg string, keyvals ...any) {
	l.logw(zapcore.ErrorLevel, msg, append(keyvals, "err", err)...)
}
func (l *logger) Panicf(format string, args ...any) { l.logf(zapcore.PanicLevel, format, args...) }
func (l *logger) Panicw(msg string, keyvals ...any) { l.logw(zapcore.PanicLevel, msg, keyvals...) }
func (l *logger) Fatalf(format string, args ...any) { l.logf(zapcore.FatalLevel, format, args...) }
func (l *logger) Fatalw(msg string, keyvals ...any) { l.logw(zapcore.FatalLevel, msg, keyvals...) }

func AddCallerSkip(skip int) Logger {
	return std.AddCallerSkip(skip)
}

func (l *logger) AddCallerSkip(skip int) Logger {
	lc := l.clone()
	lc.z = lc.z.WithOptions(zap.AddCallerSkip(skip))
	return lc
}

// W 解析传入的 context，尝试提取关注的键值，并添加到 zap.Logger 结构化日志中.
func W(ctx context.Context) Logger {
	return std.W(ctx)
}

// W 方法，根据 context 提取字段并添加到日志中
func (l *logger) W(ctx context.Context) Logger {
	lc := l.clone()

	for fieldName, extractor := range l.contextExtractors {
		if val := extractor(ctx); val != "" {
			lc.z = lc.z.With(zap.String(fieldName, val))
		}
	}

	return lc
}

// clone 深度拷贝 logger.
func (l *logger) clone() *logger {
	copied := *l
	return &copied
}

// logf 通用格式化日志方法封装
func (l *logger) logf(level zapcore.Level, format string, args ...any) {
	switch level {
	case zapcore.DebugLevel:
		l.z.Sugar().Debugf(format, args...)
	case zapcore.InfoLevel:
		l.z.Sugar().Infof(format, args...)
	case zapcore.WarnLevel:
		l.z.Sugar().Warnf(format, args...)
	case zapcore.ErrorLevel:
		l.z.Sugar().Errorf(format, args...)
	case zapcore.PanicLevel:
		l.z.Sugar().Panicf(format, args...)
	case zapcore.FatalLevel:
		l.z.Sugar().Fatalf(format, args...)
	}
}

// logw 通用结构化日志方法封装
func (l *logger) logw(level zapcore.Level, msg string, args ...any) {
	switch level {
	case zapcore.DebugLevel:
		l.z.Sugar().Debugw(msg, args...)
	case zapcore.InfoLevel:
		l.z.Sugar().Infow(msg, args...)
	case zapcore.WarnLevel:
		l.z.Sugar().Warnw(msg, args...)
	case zapcore.ErrorLevel:
		l.z.Sugar().Errorw(msg, args...)
	case zapcore.PanicLevel:
		l.z.Sugar().Panicw(msg, args...)
	case zapcore.FatalLevel:
		l.z.Sugar().Fatalw(msg, args...)
	}
}
