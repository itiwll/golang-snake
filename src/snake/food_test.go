/*
* @Author: itiwll
* @Date:   2014-12-02 09:42:42
* @Last Modified by:   itiwll
* @Last Modified time: 2014-12-05 10:30:39
 */

package main

import (
	"fmt"
	"testing"
)

var (
	f fooder
	m maper = maper{100, 100}
)

var testSnakes []*snake

func TestFood(t *testing.T) {
	f.produceFood(m, testSnakes)
	f.produceFood(m, testSnakes)
	f.produceFood(m, testSnakes)
	f.clearFood(2)
	fmt.Println(f.Foods)
	if false {
		t.Fail()
	}
}
