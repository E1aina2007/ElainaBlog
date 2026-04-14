package zaplogger

import (
	"ElainaWeb/config"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 全局 Zap 日志实例
var Logger *zap.Logger

// InitLogger 初始化 Zap 日志实例，根据配置设置日志输出、编码格式和日志级别
func InitLogger() *zap.Logger {
	zapCfg := config.GlobalConfig.Zap

	// 创建日志文件写入器（基于 lumberjack 实现日志轮转）
	writeSyncer := getLogWriter(zapCfg.FileName, zapCfg.MaxSize, zapCfg.MaxBackups, zapCfg.MaxAge)

	// 如果开启了控制台打印，则同时输出到文件和控制台
	if zapCfg.IsConsolePrint {
		writeSyncer = zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout))
	}

	// 获取日志编码器（JSON 格式）
	encoder := getEncoder()

	// 解析配置中的日志级别字符串
	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(zapCfg.Level))
	if err != nil {
		log.Fatalf("解析日志级别失败: %v", err)
	}

	// 组装核心：编码器 + 写入器 + 日志级别
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)
	// 创建 Logger，AddCaller 会在日志中记录调用位置
	logger := zap.New(core, zap.AddCaller())
	return logger
}

// getLogWriter 创建日志写入器，使用 lumberjack 实现日志文件自动轮转
// filename: 日志文件路径, maxSize: 单文件最大MB, maxBackups: 保留旧文件数, maxAge: 保留天数
func getLogWriter(filename string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filename,   // 日志文件路径
		MaxSize:    maxSize,    // 单个日志文件最大大小（MB）
		MaxBackups: maxBackups, // 保留的旧日志文件最大数量
		MaxAge:     maxAge,     // 旧日志文件保留的最大天数
	}
	return zapcore.AddSync(lumberjackLogger)
}

// getEncoder 创建日志编码器，使用 JSON 格式输出，自定义时间和级别的编码方式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionConfig()
	encoderConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder         // 时间格式：ISO8601（如 2026-04-11T16:47:00+08:00）
	encoderConfig.EncoderConfig.TimeKey = "time"                                // 时间字段名
	encoderConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder       // 日志级别大写（INFO, WARN, ERROR）
	encoderConfig.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder // 耗时以秒为单位
	encoderConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder       // 调用位置使用短路径

	return zapcore.NewJSONEncoder(encoderConfig.EncoderConfig)
}
