package api

import (
	"fmt"
	"net/http"
	"time"

	"docker_service/internal/logger"

	"github.com/gdygd/goglib"
)

func sendSse(data goglib.EventData) {
	logger.Log.Print(1, "Active sse session : %v", ActivesseSessionList)
	for _, actSession := range ActivesseSessionList {
		CheckSSEMsgChannel(actSession.Key)

		SseMsgChan[actSession.Key] <- data
	}
}

// ------------------------------------------------------------------------------
// processEventMsg
// ------------------------------------------------------------------------------
func ProcessEventMsg() {
	for {
		select {
		case event := <-goglib.ChEvent:
			logger.Log.Print(1, "Get Event message [%s]", event.Msgtype)

			if len(event.Msgtype) > 0 {
				msg := &event
				sendSse(*msg)
			} else {
				logger.Log.Error("undefined sse..[%s](%d)", event.Msgtype, event.Id)
			}
		}
	}
}

func handleSSE() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get sse session key
		sessionKey := GetSSeSessionKey()
		defer func() {
			logger.Log.Print(2, "Close sse.. [%d]", sessionKey)
			ClearSSeSessionKey(sessionKey)
		}()

		if sessionKey == 0 {
			// invalid key...
			logger.Log.Error("Access handleSSE invalid key.. [%d]", sessionKey)
			<-r.Context().Done()
			return
		}

		// prepare the header
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// prepare the flusher
		flusher, _ := w.(http.Flusher)

		// trap the request under loop forever
		for {
			select {

			case <-r.Context().Done():
				return
			default:
				sseMsg, ok := PopSSEMsgChannel(sessionKey)
				if ok {
					btData := sseMsg.PrepareMessage()
					fmt.Fprintf(w, "%s\n", btData)

					flusher.Flush()
				}
			}
			time.Sleep(time.Millisecond * 5)
		}
	}
}
