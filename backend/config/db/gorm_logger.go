// gorm_logger.go 将 GORM 日志桥接到项目的 zap logger
package db

import (
	"ElainaBlog/pkg/zaplogger"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger 适配 GORM 的 logger.Interface，输出到 zap
type GormLogger struct {
	SlowThreshold        time.Duration
	IgnoreRecordNotFound bool
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		SlowThreshold:        200 * time.Millisecond,
		IgnoreRecordNotFound: true,
	}
}

func (l *GormLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...any) {
	zaplogger.Logger.Sugar().Infof(msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	zaplogger.Logger.Sugar().Warnf(msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...any) {
	zaplogger.Logger.Sugar().Errorf(msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil && !(l.IgnoreRecordNotFound && errors.Is(err, gorm.ErrRecordNotFound)) {
		zaplogger.Logger.Error("SQL Error",
			zap.String("sql", sql),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.Error(err),
		)
		return
	}

	if l.SlowThreshold > 0 && elapsed > l.SlowThreshold {
		zaplogger.Logger.Warn("Slow SQL",
			zap.String("sql", sql),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
		)
		return
	}

	zaplogger.Logger.Debug("SQL",
		zap.String("sql", sql),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
	)
}
