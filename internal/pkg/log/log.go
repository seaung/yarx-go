package log

import (
	"context"
	"sync"
	"time"

	"github.com/seaung/yarx-go/internal/pkg/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, options ...interface{})
	Info(msg string, options ...interface{})
	Wanning(msg string, options ...interface{})
	Errorw(msg string, options ...interface{})
	Fatalw(msg string, options ...interface{})
	Sync()
}

type zapLogger struct {
	z *zap.Logger
}

var _ Logger = &zapLogger{}

var (
	mutx sync.Mutex

	std = NewZapLogger(NewOptions())
)

func Init(opts *Options) {}

func NewZapLogger(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	encodingConfig := zap.NewProductionEncoderConfig()

	encodingConfig.MessageKey = "message"

	encodingConfig.TimeKey = "timestamp"

	encodingConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	encodingConfig.EncodeDuration = func(d time.Duration, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendFloat64(float64(d) / float64(time.Millisecond))
	}

	config := &zap.Config{
		DisableCaller:     opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Encoding:          opts.Format,
		EncoderConfig:     encodingConfig,
		OutputPaths:       opts.OutputPaths,
		ErrorOutputPaths:  []string{"stderr"},
	}

	z, err := config.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	logger := &zapLogger{z: z}
	zap.RedirectStdLog(z)

	return logger
}

func Sync() { std.Sync() }

func (l *zapLogger) Sync() {
	_ = l.z.Sync()
}

func Debug(msg string, options ...interface{}) {
	std.z.Sugar().Debugw(msg, options...)
}

func (l *zapLogger) Debug(msg string, options ...interface{}) {
	l.z.Sugar().Debugw(msg, options...)
}

func Info(msg string, options ...interface{}) {
	std.z.Sugar().Infow(msg, options...)
}

func (l *zapLogger) Info(msg string, options ...interface{}) {
	l.z.Sugar().Debugw(msg, options...)
}

func Wanning(msg string, options ...interface{}) {
	std.z.Sugar().Warnw(msg, options...)
}

func (l *zapLogger) Wanning(msg string, options ...interface{}) {
	l.z.Sugar().Warnw(msg, options...)
}

func Errorw(msg string, options ...interface{}) {
	std.z.Sugar().Errorw(msg, options...)
}

func (l *zapLogger) Errorw(msg string, options ...interface{}) {
	l.z.Sugar().Errorw(msg, options...)
}

func Fatalw(msg string, options ...interface{}) {
	std.z.Sugar().Fatalw(msg, options...)
}

func (l *zapLogger) Fatalw(msg string, options ...interface{}) {
	l.z.Sugar().Fatalw(msg, options...)
}

func C(ctx context.Context) *zapLogger {
	return std.C(ctx)
}

func (l *zapLogger) clone() *zapLogger {
	lc := *l
	return &lc
}

func (l *zapLogger) C(ctx context.Context) *zapLogger {
	lc := l.clone()

	if requestID := ctx.Value(constants.XRequestIDKey); requestID != nil {
		lc.z = lc.z.With(zap.Any(constants.XUsernameKey, requestID))
	}

	if userID := ctx.Value(constants.XUsernameKey); userID != nil {
		lc.z = lc.z.With(zap.Any(constants.XUsernameKey, userID))
	}

	return lc
}
