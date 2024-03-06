package events

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

type Message struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

type sseSender struct {
	response *echo.Response
	enc      *json.Encoder
}

func NewSSESender(response *echo.Response) *sseSender {
	return &sseSender{
		response: response,
		enc:      json.NewEncoder(response),
	}
}

func (sender *sseSender) send(message *Message) error {

	sender.response.Write([]byte("data: "))
	if err := sender.enc.Encode(message); err != nil {
		return err
	}
	sender.response.Write([]byte("\n\n"))
	sender.response.Flush()
	return nil
}

func Register(e *echo.Group) {

	group := e
	group.GET("", index)
}
