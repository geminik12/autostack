/**FileHeader
 * @Author: Liangkang Zhang
 * @Date: 2026/1/24 00:47:23
 * @LastEditors: Liangkang Zhang
 * @LastEditTime: 2026/1/25 00:24:45
 * @Description:
 * @Copyright: Copyright (©)}) 2026 Liangkang Zhang<lkzhang98@gmail.com>. All rights reserved. Use of this source code is governed by a MIT style license that can be found in the LICENSE file.. All rights reserved.
 * @Email: lkzhang98@gmail.com
 * @Repository: https://github.com/geminik12/autostack
 */
// Package log provides a simple logger.

package log

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
)

var (
	infoStr       = "%s[info] "
	warnStr       = "%s[warn] "
	errStr        = "%s[error] "
	traceStr      = "[%s][%.3fms] [rows:%v] %s"
	traceWarnStr  = "%s %s[%.3fms] [rows:%v] %s"
	traceErrStr   = "%s %s[%.3fms] [rows:%v] %s"
	slowThreshold = 200 * time.Millisecond
)

func (l *logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	lc := l.clone()
	switch {
	case level <= gormlogger.Silent:
		lc.z = lc.z.WithOptions(zap.IncreaseLevel(zapcore.FatalLevel + 1))
	case level <= gormlogger.Error:
		lc.z = lc.z.WithOptions(zap.IncreaseLevel(zapcore.ErrorLevel))
	case level <= gormlogger.Warn:
		lc.z = lc.z.WithOptions(zap.IncreaseLevel(zapcore.WarnLevel))
	case level <= gormlogger.Info:
		lc.z = lc.z.WithOptions(zap.IncreaseLevel(zapcore.InfoLevel))
	default:
	}

	return lc
}

func (l *logger) Info(ctx context.Context, msg string, keyvals ...any) {
	l.AddCallerSkip(1).Infof(infoStr+msg, keyvals...)
}

func (l *logger) Warn(ctx context.Context, msg string, keyvals ...any) {
	l.AddCallerSkip(1).Warnf(warnStr+msg, keyvals...)
}

func (l *logger) Error(ctx context.Context, msg string, keyvals ...any) {
	l.AddCallerSkip(1).Errorf(errStr+msg, keyvals...)
}

func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	l = l.AddCallerSkip(1).(*logger)
	switch {
	case err != nil:
		sql, rows := fc()
		if rows == -1 {
			l.logf(zapcore.ErrorLevel, traceErrStr, "", err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.logf(zapcore.ErrorLevel, traceErrStr, "", err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > slowThreshold && slowThreshold != 0:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", slowThreshold)
		if rows == -1 {
			l.logf(zapcore.WarnLevel, traceWarnStr, "", slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.logf(zapcore.WarnLevel, traceWarnStr, "", slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	default:
		sql, rows := fc()
		if rows == -1 {
			l.logf(zapcore.InfoLevel, traceStr, "", float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.logf(zapcore.InfoLevel, traceStr, "", float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
