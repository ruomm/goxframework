package loggerx

import (
	"errors"
	"fmt"
	"github.com/ruomm/goxframework/gox/corex"
	"github.com/ruomm/goxframework/gox/refx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

var Logger *zap.Logger
var rootDirPath string
var rootDirLength int
var serviceField zap.Field
var instanceField zap.Field

type XCallerSkipHandler func(file string, line int) int

var xCallerSkipHandler XCallerSkipHandler = nil

func ConfigCallerSkipHandler(handler XCallerSkipHandler) {
	xCallerSkipHandler = handler
}
func InitLogger(logConfig interface{}, callerSkipHandler XCallerSkipHandler) error {
	// 获取当前文件绝对路径，可以减少路径长度
	//var thisPath = "logger/zap_logger.go"
	//_, file, _, _ := runtime.Caller(0)
	//if len(file) > len(thisPath) {
	//	rootDirLength = len(file) - len(thisPath)
	//	rootDirPath = file[0:rootDirLength]
	//} else {
	//	rootDirLength = 0
	//	rootDirPath = ""
	//}
	logConfigInit := LogConfigs{}
	errG, transFailsKeys := refx.XRefStructCopy(logConfig, &logConfigInit)
	if errG != nil {
		return errG
	}
	if len(transFailsKeys) > 0 {
		return errors.New("logger config init err, some field config err:" + strings.Join(transFailsKeys, ","))
	}
	xCallerSkipHandler = callerSkipHandler
	rootDirPath = corex.GetCurrentPath()
	rootDirLength = len(rootDirPath)
	serviceField = zap.String("service", logConfigInit.ServiceName)
	instanceField = zap.String("instance", logConfigInit.InstanceName)
	//initFields := getInitFields(&logConfigInit)
	encoder := getLogEncoder()
	writer := getLogWriter(&logConfigInit)
	level := getLogLevel(&logConfigInit)
	core := zapcore.NewCore(encoder, writer, level)
	caller := zap.AddCaller()
	zap.AddCallerSkip(1)
	// 开启文件及行号
	//development := zap.Development()
	// 构造日志
	//Logger = zap.New(core, caller, development, zap.Fields(initFields...))
	//Logger = zap.New(core, caller, zap.Fields(initFields...))
	Logger = zap.New(core, caller)
	return nil
}

func getLogWriter(logConfig *LogConfigs) zapcore.WriteSyncer {
	fileName := fmt.Sprintf("./logs/%s-%s.log", logConfig.ServiceName, logConfig.InstanceName)
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
	fields = append(fields, zap.String("instance", logConfig.InstanceName))
	return fields
}

func Info(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Logger.Info(message, fields...)
}
func Debug(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Logger.Debug(message, fields...)
}
func Error(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Logger.Error(message, fields...)
}
func Warn(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	Logger.Warn(message, fields...)
}
func getCallerInfoForLog() (callerFields []zap.Field) {

	pc, file, line, ok := runtime.Caller(2) // 回溯两层，拿到写日志的调用方的函数信息
	if nil != xCallerSkipHandler {
		callSkip := xCallerSkipHandler(file, line)
		if callSkip > 0 {
			pc, file, line, ok = runtime.Caller(2 + callSkip)
		}
	}
	if !ok {
		return
	}
	var realFile string
	if rootDirLength > 0 && len(file) > rootDirLength {
		realFile = file[rootDirLength:]
	} else {
		realFile = file
	}
	funcName := runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName) //Base函数返回路径的最后一个元素，只保留函数名
	//callerFields = append(callerFields, zap.String("func", funcName), zap.String("file", realFile), zap.Int("line", line))
	callerFields = append(callerFields, serviceField, instanceField, zap.String("func", funcName), zap.String("lineNo", realFile+":"+strconv.Itoa(line)))
	return
}
