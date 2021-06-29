package app

import (
	"math/rand"
	"strings"
	"time"
)

var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var numbers = "1234567890"
var pool = strings.Split(alphabet+numbers, "")

func generate(length uint) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	result := []string{}

	for i := 0; i < int(length); i++ {
		r := r1.Intn(len(pool))
		result = append(result, pool[r])
	}

	return strings.Join(result, "")
}
