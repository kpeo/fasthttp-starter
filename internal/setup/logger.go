package setup

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

type LoggerConfig struct {
	Level        string
	TraceLevel   string
	Format       string
	Debug        bool
	Color        bool
	FullCaller   bool
	NoDisclaimer bool
	Sampling     *zap.SamplingConfig
}

const (
	defaultSamplingInitial    = 100
	defaultSamplingThereafter = 100
)

func NewLogger(lcfg *LoggerConfig, app *Settings) (*zap.Logger, error) {
	var cfg zap.Config
	if lcfg.Debug {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	if lcfg.Sampling != nil {
		cfg.Sampling = lcfg.Sampling
	}

	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stdout"}

	cfg.Encoding = lcfg.SafeFormat()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if lcfg.Color {
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	if lcfg.FullCaller {
		cfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	}
	cfg.Level = SafeLevel(lcfg.Level, zapcore.InfoLevel)
	traceLevel := SafeLevel(lcfg.TraceLevel, zapcore.WarnLevel)
	l, err := cfg.Build(
		// enable trace for current log-level only
		zap.AddStacktrace(traceLevel))
	if err != nil {
		return nil, err
	}
	// disable showing app's name & version
	if lcfg.NoDisclaimer {
		return l, nil
	}

	return l.With(
		zap.String("app_name", app.Name),
		zap.String("app_version", app.Version),
	), nil
}

func NewSugaredLogger(log *zap.Logger) *zap.SugaredLogger {
	return log.Sugar()
}

func NewLoggerConfig(v *viper.Viper) *LoggerConfig {
	cfg := &LoggerConfig{
		Debug:        v.GetBool("debug"),
		Level:        v.GetString("logger.level"),
		TraceLevel:   v.GetString("logger.trace_level"),
		Format:       v.GetString("logger.format"),
		Color:        v.GetBool("logger.color"),
		FullCaller:   v.GetBool("logger.full_caller"),
		NoDisclaimer: v.GetBool("logger.no_disclaimer"),
	}

	if v.IsSet("logger.sampling") {
		cfg.Sampling = &zap.SamplingConfig{
			Initial:    defaultSamplingInitial,
			Thereafter: defaultSamplingThereafter,
		}
		if val := v.GetInt("logger.sampling.initial"); val > 0 {
			cfg.Sampling.Initial = val
		}
		if val := v.GetInt("logger.sampling.thereafter"); val > 0 {
			cfg.Sampling.Thereafter = val
		}
	}
	return cfg
}

func SafeLevel(level string, defaultLevel zapcore.Level) zap.AtomicLevel {
	switch strings.ToLower(level) {
	case "debug":
		return zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "panic":
		return zap.NewAtomicLevelAt(zapcore.PanicLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		return zap.NewAtomicLevelAt(defaultLevel)
	}
}

func (lcfg LoggerConfig) SafeFormat() string {
	switch lcfg.Format {
	case "console":
	case "json":
	default:
		return "json"
	}
	return lcfg.Format
}
