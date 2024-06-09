package main

import (
	"context"
	"fmt"
	"math/rand"
)

// UserInfo is a struct that contains the user's name, cursor color, and the channel they belong to
type UserInfo struct {
	UserName    string `json:"userName"`
	Channel     string `json:"channel"`
	CursorColor string `json:"cursorColor"`
}

// CursorPosition is a struct that contains the x and y coordinates of the cursor
type CursorPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func NewUserInfo(ctx context.Context, rdb RedisHandler, userName, channel string) (UserInfo, error) {
	uinfo := UserInfo{UserName: userName, Channel: channel}
	err := uinfo.SetCursorColor(ctx, rdb)
	if err != nil {
		return UserInfo{}, err
	}
	return uinfo, nil
}

func (msg *UserInfo) SetCursorColor(ctx context.Context, rdb RedisHandler) (err error) {
	msg.CursorColor, err = RandomColor(ctx, rdb, msg.Channel)
	if err != nil {
		return err
	}
	return nil
}

// RandomColor generates a random RGB color that's not already in use
func RandomColor(ctx context.Context, rdb RedisHandler, key string) (string, error) {
	for {
		r := rand.Intn(256)
		g := rand.Intn(256)
		b := rand.Intn(256)
		color := fmt.Sprintf("rgb(%d, %d, %d)", r, g, b)

		// Check if the color is already in use
		exists, err := rdb.SIsMember(ctx, key, color).Result()
		if err != nil {
			return "", err
		}
		if !exists {
			// The color is not in use, so add it to the set and return it
			err = rdb.SAdd(ctx, key, color).Err()
			if err != nil {
				return "", err
			}
			return color, nil
		}
		// The color is in use, so generate a new color
	}
}
