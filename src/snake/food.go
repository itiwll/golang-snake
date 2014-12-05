/*
* @Author: itiwll
* @Date:   2014-12-02 09:20:22
* @Last Modified by:   itiwll
* @Last Modified time: 2014-12-05 10:28:19
 */

package main

import (
	"math/rand"
)

var randSrc int64 = 0

type fooder struct {
	Foods    [][2]int
	quantity int
}

func (f *fooder) setQuantity(n int) {
	f.quantity = n
}

func (f *fooder) produceFood(gameMap maper, Snakes []*snake) {
S:
	x := randValue(0, gameMap.Width)
	y := randValue(0, gameMap.Height)

	for _, v := range Snakes {
		for _, vv := range v.Body {
			if vv[0] == x && vv[1] == y {
				goto S
			}
		}
	}

	for _, v := range f.Foods {
		if v[0] == x && v[1] == y {
			goto S
		}
	}

	f.Foods = append(f.Foods, [2]int{x, y})
}

func (f *fooder) clearFood(i int) {
	f.Foods = append(f.Foods[:i], f.Foods[i+1:]...)
}

func randValue(begin int, size int) (value int) {
	randSrc++
	r := rand.New(rand.NewSource(randSrc))
	value = r.Intn(size) + begin
	return
}
