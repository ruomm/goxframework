/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/4/11 23:16
 * @version 1.0
 */
package loggerx

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

type MutilLoggerX struct {
	// 默认的Loggerx
	loggerX *LoggerX
	// 其他的Loggerx
	mutilLoggerXMap map[string]*LoggerX
	mutilHandler    XMutilLoggerHandler
}

var MutilLogger *MutilLoggerX

func InitMutilLogger(logConfig interface{}, workDirPath string, instanceName string, skipHandler XCallerSkipHandler, contextHandler XContextHandler) error {
	loggerx, err := generateLoggerX(logConfig, workDirPath, instanceName, 3, skipHandler, contextHandler)
	if err != nil {
		return err
	}
	if nil == MutilLogger {
		mutilLoggerx := MutilLoggerX{
			loggerX:         nil,
			mutilLoggerXMap: make(map[string]*LoggerX),
		}
		MutilLogger = &mutilLoggerx
	}
	if len(instanceName) > 0 {
		MutilLogger.mutilLoggerXMap[instanceName] = loggerx
	} else {
		MutilLogger.loggerX = loggerx
	}
	return nil
}

type XMutilLoggerHandler func(message string) string

func ConfigMutilLoggerHandler(handler XMutilLoggerHandler) {
	if nil == MutilLogger {
		mutilLoggerx := MutilLoggerX{
			loggerX:         nil,
			mutilLoggerXMap: make(map[string]*LoggerX),
		}
		MutilLogger = &mutilLoggerx
	}
	MutilLogger.mutilHandler = handler
}

func (mutilLogger MutilLoggerX) getLoogerX(msg string) (*LoggerX, string) {
	if nil == mutilLogger.mutilHandler {
		return mutilLogger.loggerX, msg
	} else {
		instanceName := mutilLogger.mutilHandler(msg)
		if len(instanceName) <= 0 {
			return mutilLogger.loggerX, msg
		}
		_, exists := mutilLogger.mutilLoggerXMap[instanceName]
		if exists {
			if strings.HasSuffix(msg, instanceName) {
				return mutilLogger.mutilLoggerXMap[instanceName], msg[len(instanceName):]
			} else {
				return mutilLogger.mutilLoggerXMap[instanceName], msg
			}
		} else {
			return mutilLogger.loggerX, msg
		}
	}
}

func (mutilLogger MutilLoggerX) Log(lvl zapcore.Level, message string, fields ...zap.Field) {
	looger, msg := mutilLogger.getLoogerX(message)
	looger.Log(lvl, msg, fields...)
}

func (mutilLogger MutilLoggerX) Debug(message string, fields ...zap.Field) {
	looger, msg := mutilLogger.getLoogerX(message)
	looger.Debug(msg, fields...)
}

func (mutilLogger MutilLoggerX) Info(message string, fields ...zap.Field) {
	looger, msg := mutilLogger.getLoogerX(message)
	looger.Info(msg, fields...)
}

func (mutilLogger MutilLoggerX) Warn(message string, fields ...zap.Field) {
	looger, msg := mutilLogger.getLoogerX(message)
	looger.Warn(msg, fields...)
}

func (mutilLogger MutilLoggerX) Error(message string, fields ...zap.Field) {
	looger, msg := mutilLogger.getLoogerX(message)
	looger.Error(msg, fields...)
}

func (mutilLogger MutilLoggerX) DPanic(message string, fields ...zap.Field) {
	looger, msg := mutilLogger.getLoogerX(message)
	looger.DPanic(msg, fields...)
}

func (mutilLogger MutilLoggerX) Panic(message string, fields ...zap.Field) {
	looger, msg := mutilLogger.getLoogerX(message)
	looger.Panic(msg, fields...)
}

func (mutilLogger MutilLoggerX) Fatal(message string, fields ...zap.Field) {
	looger, msg := mutilLogger.getLoogerX(message)
	looger.Fatal(msg, fields...)
}

func (mutilLogger MutilLoggerX) LogWithCtx(ctx context.Context, lvl zapcore.Level, message string, fields ...zap.Field) {
	looger, msg := mutilLogger.getLoogerX(message)
	looger.LogWithCtx(ctx, lvl, msg, fields...)
}

func (mutilLogger MutilLoggerX) DebugWithCtx(ctx context.Context, message string, fields ...zap.Field) {
	looger, msg := mutilLogger.getLoogerX(message)
	looger.DebugWithCtx(ctx, msg, fields...)
}

func (mutilLogger MutilLoggerX) InfoWithCtx(ctx context.Context, message string, fields ...zap.Field) {
	looger, msg := mutilLogger.getLoogerX(message)
	looger.InfoWithCtx(ctx, msg, fields...)
}

func (mutilLogger MutilLoggerX) WarnWithCtx(ctx context.Context, message string, fields ...zap.Field) {
	looger, msg := mutilLogger.getLoogerX(message)
	looger.WarnWithCtx(ctx, msg, fields...)
}

func (mutilLogger MutilLoggerX) ErrorWithCtx(ctx context.Context, message string, fields ...zap.Field) {
	looger, msg := mutilLogger.getLoogerX(message)
	looger.ErrorWithCtx(ctx, msg, fields...)
}

func (mutilLogger MutilLoggerX) DPanicWithCtx(ctx context.Context, message string, fields ...zap.Field) {
	looger, msg := mutilLogger.getLoogerX(message)
	looger.DPanicWithCtx(ctx, msg, fields...)
}

func (mutilLogger MutilLoggerX) PanicWithCtx(ctx context.Context, message string, fields ...zap.Field) {
	looger, msg := mutilLogger.getLoogerX(message)
	looger.PanicWithCtx(ctx, msg, fields...)
}

func (mutilLogger MutilLoggerX) FatalWithCtx(ctx context.Context, message string, fields ...zap.Field) {
	looger, msg := mutilLogger.getLoogerX(message)
	looger.FatalWithCtx(ctx, msg, fields...)
}
