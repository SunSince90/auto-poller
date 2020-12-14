package main

import (
	"math/rand"
	"time"
)

func nextRandom(min, max uint) uint {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	_min := int(min)
	_max := int(max)
	next := r.Intn(_max-_min+1) + _min

	return uint(next)
}
