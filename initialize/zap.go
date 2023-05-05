package initialize

import (
	"data-transform/global"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var level zapcore.Level

func Zap() (logger *zap.Logger) {
	//if info, err := os.Stat(global.Config.Zap.Director); err != nil && os.IsNotExist(err) || !info.IsDir() { // 目录不存在则创建目录
	//	fmt.Printf("create %v directory\n", global.Config.Zap.Director)
	//	_ = os.Mkdir(global.Config.Zap.Director, os.ModePerm)
	//}

	switch global.Config.Zap.Level { // 初始化配置文件的Level
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	config := zapcore.EncoderConfig{
		MessageKey:   "msg",                       //结构化（json）输出：msg的key
		LevelKey:     "level",                     //结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:      "ts",                        //结构化（json）输出：时间的key（INFO，WARN，ERROR等）
		CallerKey:    "file",                      //结构化（json）输出：打印日志的文件对应的Key
		EncodeLevel:  zapcore.CapitalLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: zapcore.ShortCallerEncoder,  //采用短文件路径编码输出（test/main.go:14 ）
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}, //输出的时间格式
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		}, //
	}
	//////自定义日志级别：自定义Info级别
	//infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl < zapcore.WarnLevel && lvl >= level
	//})
	////自定义日志级别：自定义Warn级别
	//warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl >= zapcore.WarnLevel && lvl >= level
	//})
	// 获取io.Writer的实现
	// 实现多个输出
	core := zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), level) //同时将日志输出到控制台，NewJSONEncoder 是结构化输出
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
	return logger
}
