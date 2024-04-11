package loggerx

import (
	"errors"
	"fmt"
	"github.com/ruomm/goxframework/gox/refx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"runtime"
	"strings"
)

type LoggerX struct {
	ZapLogger          *zap.Logger
	Development        bool
	xCallerSkipHandler XCallerSkipHandler
	workPath           string
	lenWorkPath        int
}

// var Logger *zap.Logger
var Logger *LoggerX

type XCallerSkipHandler func(file string, line int) int

func ConfigCallerSkipHandler(handler XCallerSkipHandler) {
	if nil == Logger {
		return
	} else {
		Logger.xCallerSkipHandler = handler
	}
}
func InitLogger(logConfig interface{}, workDirPath string, handler XCallerSkipHandler) error {
	loggerx, err := generateLoggerX(logConfig, workDirPath, "", 1, handler)
	if err != nil {
		return err
	} else {
		Logger = loggerx
		return nil
	}
}

func generateLoggerX(logConfig interface{}, workDirPath string, instanceName string, callerSkip int, handler XCallerSkipHandler) (*LoggerX, error) {
	//workPath = workDirPath
	//lenWorkPath = len(workPath)
	//加载配置文件
	logConfigInit := LogConfigs{}
	errG, transFailsKeys := refx.XRefStructCopy(logConfig, &logConfigInit)
	if len(instanceName) > 0 {
		logConfigInit.InstanceName = instanceName
	}
	if errG != nil {
		return nil, errG
	}
	if len(transFailsKeys) > 0 {
		return nil, errors.New("logger config init err, some field config err:" + strings.Join(transFailsKeys, ","))
	}
	//xCallerSkipHandler = callerSkipHandler
	//serviceField = zap.String("service", logConfigInit.ServiceName)
	//if len(logConfigInit.InstanceName) > 0 {
	//	zapField := zap.String("instance", logConfigInit.InstanceName)
	//	instanceField = &zapField
	//} else {
	//	instanceField = nil
	//}
	// 开始配置文件
	initFields := getInitFields(&logConfigInit)
	encoder := getLogEncoder()
	writer := getLogWriter(&logConfigInit)
	level := getLogLevel(&logConfigInit)
	core := zapcore.NewCore(encoder, writer, level)
	caller := zap.AddCaller()
	zap.AddCallerSkip(1)
	// 开启文件及行号
	development := zap.Development()
	// 构造日志
	zapLogger := zap.New(core, caller, development, zap.Fields(initFields...))
	// 构造日志
	//Logger = zap.New(core, caller, zap.Fields(initFields...))
	//Logger = zap.New(core, caller)
	//Logger = zapLogger
	loogerx := LoggerX{
		ZapLogger:          zapLogger,
		Development:        true,
		xCallerSkipHandler: handler,
		workPath:           workDirPath,
		lenWorkPath:        len(workDirPath),
	}
	return &loogerx, nil
}

func generateZapLogger(logConfig interface{}, workDirPath string, instanceName string, handler XCallerSkipHandler) (*zap.Logger, error) {
	//workPath = workDirPath
	//lenWorkPath = len(workPath)
	//加载配置文件
	logConfigInit := LogConfigs{}
	errG, transFailsKeys := refx.XRefStructCopy(logConfig, &logConfigInit)
	if len(instanceName) > 0 {
		logConfigInit.InstanceName = instanceName
	}
	if errG != nil {
		return nil, errG
	}
	if len(transFailsKeys) > 0 {
		return nil, errors.New("logger config init err, some field config err:" + strings.Join(transFailsKeys, ","))
	}
	//xCallerSkipHandler = callerSkipHandler
	//serviceField = zap.String("service", logConfigInit.ServiceName)
	//if len(logConfigInit.InstanceName) > 0 {
	//	zapField := zap.String("instance", logConfigInit.InstanceName)
	//	instanceField = &zapField
	//} else {
	//	instanceField = nil
	//}
	// 开始配置文件
	initFields := getInitFields(&logConfigInit)
	encoder := getLogEncoder()
	writer := getLogWriter(&logConfigInit)
	level := getLogLevel(&logConfigInit)
	core := zapcore.NewCore(encoder, writer, level)
	caller := zap.AddCaller()
	zap.AddCallerSkip(1)
	// 开启文件及行号
	development := zap.Development()
	// 构造日志
	zapLogger := zap.New(core, caller, development, zap.Fields(initFields...))
	// 构造日志
	//Logger = zap.New(core, caller, zap.Fields(initFields...))
	//Logger = zap.New(core, caller)
	//Logger = zapLogger
	return zapLogger, nil
}

func getLogWriter(logConfig *LogConfigs) zapcore.WriteSyncer {
	fileName := ""
	if len(logConfig.InstanceName) > 0 {
		fileName = fmt.Sprintf("./logs/%s-%s.log", logConfig.ServiceName, logConfig.InstanceName)
	} else {
		fileName = fmt.Sprintf("./logs/%s.log", logConfig.ServiceName)
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    logConfig.MaxSize,
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.MaxBackups,
		Compress:   logConfig.Compress,
	}
	if logConfig.StdOut {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger),
			zapcore.AddSync(os.Stdout))
	} else {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger))
	}
}

func getLogEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:  "time",
		LevelKey: "level",
		NameKey:  "logger",
		//CallerKey:      "lineNum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogLevel(logConfig *LogConfigs) zap.AtomicLevel {
	var level zapcore.Level
	switch logConfig.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	return atomicLevel
}

func getInitFields(logConfig *LogConfigs) (fields []zap.Field) {
	fields = append(fields, zap.String("service", logConfig.ServiceName))
	if len(logConfig.InstanceName) == 0 {
		logConfig.InstanceName, _ = os.Hostname()
	}
	if len(logConfig.InstanceName) > 0 {
		fields = append(fields, zap.String("instance", logConfig.InstanceName))
	}
	return fields
}

func (looger LoggerX) Log(lvl zapcore.Level, msg string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(fields, callerFields...)
	}
	looger.ZapLogger.Log(lvl, msg, fields...)
}

func (looger LoggerX) Debug(message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(fields, callerFields...)
	}
	looger.ZapLogger.Debug(message, fields...)
}

func (looger LoggerX) Info(message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(fields, callerFields...)
	}
	looger.ZapLogger.Info(message, fields...)
}

func (looger LoggerX) Warn(message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(fields, callerFields...)
	}
	looger.ZapLogger.Warn(message, fields...)
}

func (looger LoggerX) Error(message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(fields, callerFields...)
	}
	looger.ZapLogger.Error(message, fields...)
}

func (looger LoggerX) DPanic(message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(fields, callerFields...)
	}
	looger.ZapLogger.DPanic(message, fields...)
}

func (looger LoggerX) Panic(message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(fields, callerFields...)
	}
	looger.ZapLogger.Panic(message, fields...)
}

func (looger LoggerX) Fatal(message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(fields, callerFields...)
	}
	looger.ZapLogger.Fatal(message, fields...)
}

func (looger LoggerX) getCallerInfoForLog() (callerFields []zap.Field) {
	pc, file, line, ok := runtime.Caller(2) // 回溯两层，拿到写日志的调用方的函数信息
	if nil != looger.xCallerSkipHandler {
		callSkip := looger.xCallerSkipHandler(file, line)
		if callSkip > 0 {
			pc, file, line, ok = runtime.Caller(2 + callSkip)
		}
	}
	if !ok {
		return
	}
	var realFile = looger.parseRealFile(file)
	funcName := runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName) //Base函数返回路径的最后一个元素，只保留函数名
	//callerFields = append(callerFields, zap.String("func", funcName), zap.String("file", realFile), zap.Int("line", line))
	//callerFields = append(callerFields, initFields...)
	callerFields = append(callerFields, zap.String("func", funcName), zap.String("file", realFile), zap.Int("lineNo", line))
	return
}

func (looger LoggerX) parseRealFile(filePath string) string {
	lenFile := len(filePath)
	if lenFile <= 0 {
		return filePath
	}
	//Users/qx/go/pkg/mod/gitlab.idr.ai/turing/charging-bill.git@v0.0.0-20240407094215-5b376be23fed/logger/zap_logger.go
	indexAt := strings.Index(filePath, "@v")
	if indexAt > 0 && indexAt < lenFile {
		tmpFilePath := filePath[0:indexAt]
		indexSpec01 := strings.LastIndex(tmpFilePath, "/")
		indexSpec02 := strings.LastIndex(tmpFilePath, "\\")
		indexSepc := -1
		if indexSpec01 >= indexSpec02 {
			indexSepc = indexSpec01
		} else {
			indexSepc = indexSpec02
		}
		if indexSepc >= 0 {
			return filePath[indexSepc+1:]
		} else {
			return filePath
		}
	} else if looger.lenWorkPath > 0 {
		if strings.HasPrefix(filePath, looger.workPath) {
			return filePath[looger.lenWorkPath:]
		} else {
			return filePath
		}
	} else {
		return filePath
	}
}
