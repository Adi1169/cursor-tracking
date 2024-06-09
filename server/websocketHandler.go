package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type UserInitMessage struct {
	UserName string `json:"userName"`
	Channel  string `json:"channel"`
}

type webSocketHandler struct {
	upgrader websocket.Upgrader
	rdb      RedisHandler
	cStore   *ConnectionStore
}

func (wsh *webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}
	ctx := context.Background()

	// Read the initial message from the user
	// This message contains the user's name and the channel they want to join
	userInitMessage, err := readUserInitMessage(c)
	if err != nil {
		log.Println("error occurred while reading initial message:", err)
		c.Close()
		return
	}

	// Create a new UserInfo object
	// This object contains the user's name, cursor color, and the channel they want to join
	userInfo, err := NewUserInfo(ctx, wsh.rdb, userInitMessage.UserName, userInitMessage.Channel)
	if err != nil {
		log.Println("error occurred while intialising user:", err)
		c.Close()
		return
	}

	log.Println(userInfo.UserName, userInfo.CursorColor, userInfo.Channel)

	// Add the websocket connection to the connection store based on the channel
	exists := wsh.cStore.Set(userInfo.Channel, c)
	if !exists {
		// If this is the first connection for this channel, start a new goroutine to broadcast messages
		go wsh.rdb.Broadcast(ctx, userInfo.Channel, wsh.cStore)
	}

	// Start a new goroutine to read messages from the user
	go wsh.readLoop(ctx, c, userInfo)
}

func (wsh *webSocketHandler) readLoop(ctx context.Context, w *websocket.Conn, userInfo UserInfo) {
	defer func() {
		w.Close()
		wsh.cStore.Delete(userInfo.Channel, w)
	}()
	for {
		var cp CursorPosition
		err := w.ReadJSON(&cp)
		if err != nil {
			log.Println("error occurred while reading message:", err)
			return
		}

		err = wsh.rdb.SendPayload(ctx, RedisPayload{userInfo, cp})
		if err != nil {
			log.Println("error occurred while publishing message:", err)
			return
		}
	}
}

func readUserInitMessage(w *websocket.Conn) (UserInitMessage, error) {
	var initialMsg UserInitMessage
	err := w.ReadJSON(&initialMsg)
	if err != nil {
		return UserInitMessage{}, err
	}
	return initialMsg, nil
}
