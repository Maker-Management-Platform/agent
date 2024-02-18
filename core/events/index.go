package events

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func index(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().WriteHeader(http.StatusOK)
	uuid := uuid.New().String()

	//uuid := "batata"
	sender := NewSSESender(c.Response())

	go func() {

		time.Sleep(500 * time.Millisecond)
		err := sender.send(&Message{
			Event: "connect",
			Data: map[string]string{
				"uuid": uuid,
			},
		})
		if err != nil {
			log.Println(err)
		}
	}()

	eventChan, unregister := RegisterSession(uuid)

	for {
		select {
		case <-c.Request().Context().Done():
			unregister()
			return nil
		case s, ok := <-eventChan:
			if !ok {
				log.Println("Event chan closed, closing client")
				close(eventChan)
				return nil
			}
			err := sender.send(s)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
}
