package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"labsystem/configs"
	"labsystem/model"
	"os"
)

var Log *zap.Logger

func init() {
	var options []zap.Option
	logConfig := configs.NewLogConfig()
	// env
	encodeConfig, opt := parseLogEnv(logConfig.Env)
	// writer
	writer := parseWriter(logConfig.Output)
	// level
	logLevel := parseLogLevel(logConfig.Level)
	// options
	if opt != nil {
		options = append(options, opt)
	}
	options = append(options, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	// log common configuration
	encodeConfig.EncodeTime = zapcore.TimeEncoderOfLayout(model.TimeFormat)
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder

	logCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encodeConfig), writer, logLevel)
	Log = zap.New(logCore, options...)
}

func parseLogEnv(env configs.Environment) (zapcore.EncoderConfig, zap.Option) {
	switch env {
	case configs.Production:
		return zap.NewProductionEncoderConfig(), nil
	case configs.Development:
		return zap.NewDevelopmentEncoderConfig(), zap.Development()
	default:
		panic("invalid env param")
	}
}

func parseWriter(writer string) zapcore.WriteSyncer {
	var logWriter io.Writer
	switch writer {
	case "stdout":
		logWriter = os.Stdout
	default:
		var err error
		logWriter, err = os.Open(configs.CurProjectPath() + writer)
		if err != nil {
			panic("initialize logger failed:" + err.Error())
		}
	}

	return zapcore.AddSync(logWriter)
}

func parseLogLevel(level string) zapcore.LevelEnabler {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}
