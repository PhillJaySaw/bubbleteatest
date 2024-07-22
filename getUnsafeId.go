package main

import (
	"math/rand/v2"
	"time"
)

func getUnsafeId() string {
	epochTime := time.Now().Unix()
	randInt := rand.IntN(10000) + 1

	return string(epochTime) + string(randInt)
}
