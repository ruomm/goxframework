/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/4/11 23:16
 * @version 1.0
 */
package loggerx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MutilLoggerX struct {
	// 默认的Loggerx
	loggerX *LoggerX
	// 其他的Loggerx
	mutilLoggerXMap map[string]*LoggerX
	mutilHandler    XMutilLoggerHandler
}

var MutilLogger *MutilLoggerX

func (*MutilLoggerX) InitLogger(logConfig interface{}, workDirPath string, instanceName string, handler XCallerSkipHandler) error {
	loggerx, err := generateLoggerX(logConfig, workDirPath, instanceName, 1, handler)
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

func (mutilLogger MutilLoggerX) getLoogerX(msg string) *LoggerX {
	if nil == mutilLogger.mutilHandler {
		return mutilLogger.loggerX
	} else {
		instanceName := mutilLogger.mutilHandler(msg)
		if len(instanceName) <= 0 {
			return mutilLogger.loggerX
		}
		_, exists := mutilLogger.mutilLoggerXMap[instanceName]
		if exists {
			return mutilLogger.mutilLoggerXMap[instanceName]
		} else {
			return mutilLogger.loggerX
		}
	}
}

func (mutilLogger MutilLoggerX) Log(lvl zapcore.Level, msg string, fields ...zap.Field) {
	looger := mutilLogger.getLoogerX(msg)
	looger.Log(lvl, msg, fields...)
}

func (mutilLogger MutilLoggerX) Debug(message string, fields ...zap.Field) {
	looger := mutilLogger.getLoogerX(message)
	looger.Debug(message, fields...)
}

func (mutilLogger MutilLoggerX) Info(message string, fields ...zap.Field) {
	looger := mutilLogger.getLoogerX(message)
	looger.Info(message, fields...)
}

func (mutilLogger MutilLoggerX) Warn(message string, fields ...zap.Field) {
	looger := mutilLogger.getLoogerX(message)
	looger.Warn(message, fields...)
}

func (mutilLogger MutilLoggerX) Error(message string, fields ...zap.Field) {
	looger := mutilLogger.getLoogerX(message)
	looger.Error(message, fields...)
}

func (mutilLogger MutilLoggerX) DPanic(message string, fields ...zap.Field) {
	looger := mutilLogger.getLoogerX(message)
	looger.DPanic(message, fields...)
}

func (mutilLogger MutilLoggerX) Panic(message string, fields ...zap.Field) {
	looger := mutilLogger.getLoogerX(message)
	looger.Panic(message, fields...)
}

func (mutilLogger MutilLoggerX) Fatal(message string, fields ...zap.Field) {
	looger := mutilLogger.getLoogerX(message)
	looger.Fatal(message, fields...)
}
