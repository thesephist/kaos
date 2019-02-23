package wordid

import (
	"math/rand"
	"time"
)

const LENGTH int = 10

func Generate() string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	b := []byte{}
	for i := 0; i < LENGTH; i++ {
		b = append(b, byte(r.Intn(26)+97))
	}
	return string(b)
}
