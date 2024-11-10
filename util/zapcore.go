package util

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CustomZapCore struct {
	AcceptedLevels []zapcore.Level
	zapcore.LevelEnabler
	enc         zapcore.Encoder
	out         zapcore.WriteSyncer
	DiscordHook *DiscordHook
}

type NewCustomZapCoreParams struct {
	Level        zapcore.Level
	Enc          zapcore.Encoder
	Out          zapcore.WriteSyncer
	LevelEnabler zapcore.LevelEnabler
	DiscordHook  *DiscordHook
}

func SetZapGlobals() {
	core, err := NewCustomZapCore(&NewCustomZapCoreParams{
		Level:        zap.InfoLevel,
		Enc:          zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		Out:          zapcore.AddSync(os.Stdout),
		LevelEnabler: zapcore.InfoLevel,
		DiscordHook:  NewDiscordHook(os.Getenv("DISCORD_WEBHOOK")),
	})
	if err != nil {
		panic(err)
	}
	logger := zap.New(zapcore.NewTee(core, zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(os.Stdout), zapcore.DebugLevel)))
	zap.ReplaceGlobals(logger)
}

func NewCustomZapCore(params *NewCustomZapCoreParams) (zapcore.Core, error) {
	core := &CustomZapCore{
		AcceptedLevels: LevelThreshold(params.Level),
		LevelEnabler:   params.LevelEnabler,
		enc:            params.Enc,
		out:            params.Out,
		DiscordHook:    params.DiscordHook,
	}
	return core, nil
}

func (c *CustomZapCore) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	addFields(clone.enc, fields)
	return clone
}

func (c *CustomZapCore) clone() *CustomZapCore {
	return &CustomZapCore{
		AcceptedLevels: c.AcceptedLevels,
		LevelEnabler:   c.LevelEnabler,
		enc:            c.enc.Clone(),
		out:            c.out,
		DiscordHook:    c.DiscordHook,
	}
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}

func (c *CustomZapCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *CustomZapCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	buf, err := c.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}
	if c.isAcceptedLevel(ent.Level) {
		err = c.DiscordHook.Send(ent, fields)
	}
	buf.Free()
	if err != nil {
		return err
	}
	if ent.Level > zapcore.ErrorLevel {
		c.Sync()
	}
	return nil
}

func (c *CustomZapCore) Sync() error {
	return c.out.Sync()
}

func (c *CustomZapCore) Levels() []zapcore.Level {
	if c.AcceptedLevels == nil {
		return AllLevels
	}
	return c.AcceptedLevels
}

func (c *CustomZapCore) isAcceptedLevel(level zapcore.Level) bool {
	for _, lv := range c.Levels() {
		if lv == level {
			return true
		}
	}
	return false
}

var AllLevels = []zapcore.Level{
	zapcore.DebugLevel,
	zapcore.InfoLevel,
	zapcore.WarnLevel,
	zapcore.ErrorLevel,
	zapcore.FatalLevel,
	zapcore.PanicLevel,
}

func LevelThreshold(l zapcore.Level) []zapcore.Level {
	for i := range AllLevels {
		if AllLevels[i] == l {
			return AllLevels[i:]
		}
	}
	return []zapcore.Level{}
}
