package main

import (
	"math/rand"
	"time"
)

func nextRandom(min, max int) int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	next := r.Intn(max-min+1) + min
	return int(next)
}
