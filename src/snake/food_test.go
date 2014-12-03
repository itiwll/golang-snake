/*
* @Author: itiwll
* @Date:   2014-12-02 09:42:42
* @Last Modified by:   itiwll
* @Last Modified time: 2014-12-02 13:52:31
 */

package main

import (
	"fmt"
	"testing"
)

var f fooder

var testSnakes []*snake

func TestFood(t *testing.T) {
	f.produceFood(0, 0, 100, 100, testSnakes)
	f.produceFood(0, 0, 100, 100, testSnakes)
	f.produceFood(0, 0, 100, 100, testSnakes)
	f.clearFood(2)
	fmt.Println(f.Foods)
	if false {
		t.Fail()
	}
}
