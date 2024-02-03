package events

import (
	"encoding/json"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Message struct {
	Id    int
	Event string
	Data  any
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

	sender.response.Write([]byte("id: "))
	sender.response.Write([]byte(strconv.Itoa(message.Id)))
	sender.response.Write([]byte("\nevent: "))
	sender.response.Write([]byte(message.Event))
	sender.response.Write([]byte("\ndata: "))
	if err := sender.enc.Encode(message.Data); err != nil {
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
