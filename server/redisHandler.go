package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

// RedisPayload is a struct that contains the user's information and cursor position
type RedisPayload struct {
	UserInfo
	CursorPosition
}

type RedisHandler struct {
	*redis.Client
}

func NewRedisHandler(r *redis.Client) RedisHandler {
	return RedisHandler{r}
}

// Broadcast reads from the  channel and sends messages to clients
func (r RedisHandler) Broadcast(ctx context.Context, channel string, cStore *ConnectionStore) {
	pubsub := r.Subscribe(ctx, channel)
	defer pubsub.Close()

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Println("error occurred while setting up subscription:", err)
		return
	}

	ch := pubsub.Channel()
	for {
		redisMsg, ok := <-ch
		if !ok {
			break
		}
		var msg RedisPayload
		err := json.Unmarshal([]byte(redisMsg.Payload), &msg)
		if err != nil {
			log.Println("error occurred while unmarshalling message:", err)
			continue
		}

		// Send it to every client that are subscribed to the channel of this message
		cStore.Iterate(func(key string, value map[*websocket.Conn]bool) {
			if key == channel {
				for conn := range value {
					// Broadcast the message to the client concurrently
					go func(conn *websocket.Conn) {
						err := conn.WriteJSON(msg)
						if err != nil {
							log.Printf("error: %v", err)
						}
					}(conn)
				}
			}
		})
	}
}

func (r RedisHandler) SendPayload(ctx context.Context, msg RedisPayload) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = r.Publish(ctx, msg.Channel, string(jsonMsg)).Err()
	if err != nil {
		return err
	}
	log.Println(string(jsonMsg) + " message sent successfully on:" + msg.Channel)
	return nil
}
