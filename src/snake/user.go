package main

import (
	"github.com/gorilla/websocket"
)

var userId int = 0

type user struct {
	id    int
	snake *snake
	conn  *websocket.Conn
}

func (u *user) newSnake()(s *snake) {
	s = u.snake = newSnake()
    return
}

func newUser() (u user) {
	u = user{userId, newSnake(), nil}
	userId++
	return
}
