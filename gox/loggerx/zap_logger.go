package loggerx

import (
	"context"
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
	xContextHandler    XContextHandler
	workPath           string
	lenWorkPath        int
}

// var Logger *zap.Logger
var Logger *LoggerX

type XCallerSkipHandler func(file string, line int) int

type XContextHandler func(context.Context) ([]zap.Field, string)

func InitLogger(logConfig interface{}, workDirPath string, skipHandler XCallerSkipHandler, contextHandler XContextHandler) error {
	loggerx, err := generateLoggerX(logConfig, workDirPath, "", 1, skipHandler, contextHandler)
	if err != nil {
		return err
	} else {
		Logger = loggerx
		return nil
	}
}
func GenerateLogger(logConfig interface{}, workDirPath string, instanceName string, callerSkip int, skipHandler XCallerSkipHandler, contextHandler XContextHandler) (*LoggerX, error) {
	return generateLoggerX(logConfig, workDirPath, instanceName, callerSkip, skipHandler, contextHandler)
}
func generateLoggerX(logConfig interface{}, workDirPath string, instanceName string, callerSkip int, skipHandler XCallerSkipHandler, contextHandler XContextHandler) (*LoggerX, error) {
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
	encoder := getLogEncoder(logConfigInit.TextMode)
	writer := getLogWriter(&logConfigInit)
	level := getLogLevel(&logConfigInit)
	core := zapcore.NewCore(encoder, writer, level)
	caller := zap.AddCaller()
	zap.AddCallerSkip(callerSkip)
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
		xCallerSkipHandler: skipHandler,
		xContextHandler:    contextHandler,
		workPath:           workDirPath,
		lenWorkPath:        len(workDirPath),
	}
	return &loogerx, nil
}

func generateZapLogger(logConfig interface{}, workDirPath string, instanceName string, skipHandler XCallerSkipHandler) (*zap.Logger, error) {
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
	encoder := getLogEncoder(logConfigInit.TextMode)
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

func getLogEncoder(textMode bool) zapcore.Encoder {
	encodeLevel := zapcore.LowercaseLevelEncoder
	if textMode {
		encodeLevel = zapcore.CapitalLevelEncoder
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:  "time",
		LevelKey: "level",
		NameKey:  "logger",
		//CallerKey:      "lineNum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,                    // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	if textMode {
		return zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
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
func (looger LoggerX) ConfigCallerSkipHandler(handler XCallerSkipHandler) {
	if nil == Logger {
		return
	} else {
		looger.xCallerSkipHandler = handler
	}
}

func (looger LoggerX) ConfigContextHandler(handler XContextHandler) {
	looger.xContextHandler = handler
}

func (looger LoggerX) Log(lvl zapcore.Level, msg string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(callerFields, fields...)
	}
	looger.ZapLogger.Log(lvl, msg, fields...)
}

func (looger LoggerX) Debug(message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(callerFields, fields...)
	}
	looger.ZapLogger.Debug(message, fields...)
}

func (looger LoggerX) Info(message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(callerFields, fields...)
	}
	looger.ZapLogger.Info(message, fields...)
}

func (looger LoggerX) Warn(message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(callerFields, fields...)
	}
	looger.ZapLogger.Warn(message, fields...)
}

func (looger LoggerX) Error(message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(callerFields, fields...)
	}
	looger.ZapLogger.Error(message, fields...)
}

func (looger LoggerX) DPanic(message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(callerFields, fields...)
	}
	looger.ZapLogger.DPanic(message, fields...)
}

func (looger LoggerX) Panic(message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(callerFields, fields...)
	}
	looger.ZapLogger.Panic(message, fields...)
}

func (looger LoggerX) Fatal(message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(callerFields, fields...)
	}
	looger.ZapLogger.Fatal(message, fields...)
}

func (looger LoggerX) LogWithCtx(ctx context.Context, lvl zapcore.Level, msg string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(callerFields, fields...)
	}
	msgPrefix := ""
	if nil != Logger.xContextHandler {
		callerFields, msgPrefixTmp := looger.xContextHandler(ctx)
		if len(callerFields) > 0 {
			fields = append(fields, callerFields...)
		}
		msgPrefix = msgPrefixTmp

	}
	looger.ZapLogger.Log(lvl, msgPrefix+msg, fields...)
}

func (looger LoggerX) DebugWithCtx(ctx context.Context, message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(callerFields, fields...)
	}
	msgPrefix := ""
	if nil != Logger.xContextHandler {
		callerFields, msgPrefixTmp := looger.xContextHandler(ctx)
		if len(callerFields) > 0 {
			fields = append(fields, callerFields...)
		}
		msgPrefix = msgPrefixTmp

	}
	looger.ZapLogger.Debug(msgPrefix+message, fields...)
}

func (looger LoggerX) InfoWithCtx(ctx context.Context, message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(callerFields, fields...)
	}
	msgPrefix := ""
	if nil != Logger.xContextHandler {
		callerFields, msgPrefixTmp := looger.xContextHandler(ctx)
		if len(callerFields) > 0 {
			fields = append(fields, callerFields...)
		}
		msgPrefix = msgPrefixTmp

	}
	looger.ZapLogger.Info(msgPrefix+message, fields...)
}

func (looger LoggerX) WarnWithCtx(ctx context.Context, message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(callerFields, fields...)
	}
	msgPrefix := ""
	if nil != Logger.xContextHandler {
		callerFields, msgPrefixTmp := looger.xContextHandler(ctx)
		if len(callerFields) > 0 {
			fields = append(fields, callerFields...)
		}
		msgPrefix = msgPrefixTmp

	}
	looger.ZapLogger.Warn(msgPrefix+message, fields...)
}

func (looger LoggerX) ErrorWithCtx(ctx context.Context, message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(callerFields, fields...)
	}
	msgPrefix := ""
	if nil != Logger.xContextHandler {
		callerFields, msgPrefixTmp := looger.xContextHandler(ctx)
		if len(callerFields) > 0 {
			fields = append(fields, callerFields...)
		}
		msgPrefix = msgPrefixTmp

	}
	looger.ZapLogger.Error(msgPrefix+message, fields...)
}

func (looger LoggerX) DPanicWithCtx(ctx context.Context, message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(callerFields, fields...)
	}
	msgPrefix := ""
	if nil != Logger.xContextHandler {
		callerFields, msgPrefixTmp := looger.xContextHandler(ctx)
		if len(callerFields) > 0 {
			fields = append(fields, callerFields...)
		}
		msgPrefix = msgPrefixTmp

	}
	looger.ZapLogger.DPanic(msgPrefix+message, fields...)
}

func (looger LoggerX) PanicWithCtx(ctx context.Context, message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(callerFields, fields...)
	}
	msgPrefix := ""
	if nil != Logger.xContextHandler {
		callerFields, msgPrefixTmp := looger.xContextHandler(ctx)
		if len(callerFields) > 0 {
			fields = append(fields, callerFields...)
		}
		msgPrefix = msgPrefixTmp

	}
	looger.ZapLogger.Panic(msgPrefix+message, fields...)
}

func (looger LoggerX) FatalWithCtx(ctx context.Context, message string, fields ...zap.Field) {
	if looger.Development {
		callerFields := looger.getCallerInfoForLog()
		fields = append(callerFields, fields...)
	}
	msgPrefix := ""
	if nil != Logger.xContextHandler {
		callerFields, msgPrefixTmp := looger.xContextHandler(ctx)
		if len(callerFields) > 0 {
			fields = append(fields, callerFields...)
		}
		msgPrefix = msgPrefixTmp

	}
	looger.ZapLogger.Fatal(msgPrefix+message, fields...)
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
