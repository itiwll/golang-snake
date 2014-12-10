/*
* @Author: itiwll
* @Date:   2014-11-28 13:14:15
* @Last Modified by:   itiwll
* @Last Modified time: 2014-12-02 14:34:40
 */

package main

import (
	"fmt"
)

var snakeId int = 0

type snake struct {
	Name     string //名字
	Id       int
	Body     [][2]int //身体
	length   int      //长度
	direcion int      //方向
	Staust   int      //状态
}

// 移动 吃food
func (s *snake) Move(f *fooder) (isEat bool) {
	isEat = false

D:
	// 计算新头的位置
	head := s.Body[len(s.Body)-1]
	second := s.Body[len(s.Body)-2]

	switch s.direcion {
	case 1:
		head[1] = head[1] - 1
	case 2:
		head[0] = head[0] + 1
	case 3:
		head[1] = head[1] + 1
	case 4:
		head[0] = head[0] - 1
	}

	if head[0] == second[0] && head[1] == second[1] {
		fmt.Println("方向错误")
		switch s.direcion {
		case 1:
			s.direcion = 3
		case 2:
			s.direcion = 4
		case 3:
			s.direcion = 1
		case 4:
			s.direcion = 2
		}
		goto D
	}

	s.Body = append(s.Body, head)

	// 是否吃到food
	for i, food := range f.Foods {
		if food[0] == head[0] && food[1] == head[1] {
			// 吃
			f.clearFood(i)
			isEat = true
			fmt.Println(s.Name + ":吃")
			break
		}
	}
	// 移动
	if !isEat {
		s.Body = s.Body[1:]
	}

	return
}

// 死掉
func (s *snake) die() {
	s.Staust = 0
}

func newSnake() (s *snake) {
	s = &snake{
		"",
		snakeId,
		[][2]int{{1, 1}, {1, 2}, {1, 3}},
		3,
		3,
		1,
	}
	snakeId++
	return
}
