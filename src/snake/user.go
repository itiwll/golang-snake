package main

var userId int = 0

type user struct {
	id    int
	snake snake
}

func (u *user) newSnake() {
	u.snake = newSnake()
}

func newUser() (u user) {
	u = user{userId, newSnake()}
	userId++
	return
}
