package util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"go.uber.org/zap/zapcore"
)

type DiscordHook struct {
	AcceptedLevels []zapcore.Level
	HookUrl        string
	Username       string
}

type Payload struct {
	Username string  `json:"username"`
	Embeds   []Embed `json:"embeds"`
}

type Embed struct {
	Author      Author  `json:"author"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Color       int64   `json:"color"`
	Footer      Footer  `json:"footer"`
	Timestamp   string  `json:"timestamp,omitempty"`
	Fields      []Field `json:"fields,omitempty"`
	Thumbnail   Url     `json:"thumbnail,omitempty"`
	Image       Url     `json:"image,omitempty"`
}

type Author struct {
	Name string `json:"name"`
}

type Footer struct {
	Text string `json:"text"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Url struct {
	Url string `json:"url"`
}

func NewDiscordHook(hookUrl string) *DiscordHook {
	return &DiscordHook{
		HookUrl:  hookUrl,
		Username: "Logger",
	}
}

func (dh *DiscordHook) Send(e zapcore.Entry, fields []zapcore.Field) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	payload := Payload{
		Username: dh.Username,
		Embeds: []Embed{
			{
				Author: Author{
					Name: e.Caller.Function,
				},
				Title:       e.Level.CapitalString(),
				Description: e.Message,
				Color:       LevelColorMap[e.Level],
				Footer: Footer{
					Text: hostname,
				},
				Timestamp: e.Time.UTC().Format("2006-01-02T15:04:05.999Z"),
			}},
	}
	for _, f := range fields {
		switch f.Type {
		case zapcore.StringType:
			switch f.Key {
			case "thumb":
				payload.Embeds[0].Thumbnail = Url{
					Url: f.String,
				}
			case "image":
				payload.Embeds[0].Image = Url{
					Url: f.String,
				}
			default:
				payload.Embeds[0].Fields = append(payload.Embeds[0].Fields, Field{
					Name:   f.Key,
					Value:  f.String,
					Inline: false,
				})
			}
		case zapcore.Int16Type, zapcore.Int32Type, zapcore.Int64Type:
			payload.Embeds[0].Fields = append(payload.Embeds[0].Fields, Field{
				Name:   f.Key,
				Value:  strconv.FormatInt(f.Integer, 10),
				Inline: false,
			})
		case zapcore.Float32Type, zapcore.Float64Type:
			payload.Embeds[0].Fields = append(payload.Embeds[0].Fields, Field{
				Name:   f.Key,
				Value:  strconv.FormatFloat(float64(f.Integer), 'f', -1, 64),
				Inline: false,
			})
		case zapcore.BoolType:
			payload.Embeds[0].Fields = append(payload.Embeds[0].Fields, Field{
				Name:   f.Key,
				Value:  strconv.FormatBool(f.Integer != 0),
				Inline: false,
			})
		default:
			payload.Embeds[0].Fields = append(payload.Embeds[0].Fields, Field{
				Name:   f.Key,
				Value:  f.String,
				Inline: false,
			})
		}
	}
	req, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	res, err := http.Post(dh.HookUrl, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

var LevelColorMap = map[zapcore.Level]int64{
	zapcore.DebugLevel: 0x95a5a6,
	zapcore.InfoLevel:  0x3498db,
	zapcore.WarnLevel:  0xf1c40f,
	zapcore.ErrorLevel: 0xe67e22,
	zapcore.FatalLevel: 0xe74c3c,
	zapcore.PanicLevel: 0x9b59b6,
}
