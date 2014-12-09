/*
* @Author: itiwll
* @Date:   2014-11-28 16:11:25
* @Last Modified by:   itiwll
* @Last Modified time: 2014-12-02 14:05:52
 */

package main

import (
	"fmt"
	"testing"
)

var (
	s    = newSnake()
	food = fooder{[][2]int{{2, 4}, {1, 6}}, 1}
)

func TestSnakeMove(t *testing.T) {
	fmt.Println("测试移动")
	fmt.Println(s)
	s.Move(&food)
	fmt.Println(s)
}
