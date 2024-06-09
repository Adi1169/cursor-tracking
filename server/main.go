package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

func main() {
	redisAddr := flag.String("redis-addr", "localhost:6379", "redis server address")
	redisPassword := flag.String("redis-password", "", "redis server password")

	// Initialize Redis client
	rdb := NewRedisHandler(redis.NewClient(&redis.Options{
		Addr:     *redisAddr,
		Password: *redisPassword,
	}))

	webSocketHandler := &webSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		rdb:    rdb,
		cStore: NewConnectionStore(),
	}
	http.Handle("/ws", webSocketHandler)

	log.Println("http server started on :8101")
	err := http.ListenAndServe(":8101", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
