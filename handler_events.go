package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

// Handler for the events stream

// TODO: now the eventsstream shows which category changed, or whether the collection of categories changed
func eventsStream(notifier *Notifier) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Logger().Infof("eventsStream called")
		receiverID, err := notifier.NewReceiver()
		if err != nil {
			ctx.Logger().Errorf("eventsStream Error creating receiver: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)

		}

		// close notifier to clean up notifier channel if connection closes:

		closeNotify := ctx.Response().CloseNotify()
		go func() {
			<-closeNotify
			// connection close, do cleanup
			ctx.Logger().Infof("eventsStream: close notify triggered")
			notifier.RemoveReceiver(receiverID)
		}()

		ctx.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		ctx.Response().WriteHeader(http.StatusOK)
		ctx.Response().Flush()
		for {
			// listen on channel:
			cmd := notifier.Listen(receiverID)

			// write message to client
			_, err := ctx.Response().Write([]byte(fmt.Sprintf("data: {\"cmd\": \"%s\"}\n\n", cmd)))
			if err != nil {
				notifier.RemoveReceiver(receiverID)
				msg := fmt.Sprintf("Error writing to stream: %v", err)
				ctx.Logger().Infof("eventsStream: Error %v", msg)
				return echo.NewHTTPError(http.StatusInternalServerError, msg)
			}
			ctx.Response().Flush()
		}
	}
}
