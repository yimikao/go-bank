package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())

}
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

const alphabets = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	var sb strings.Builder
	l := len(alphabets)
	for i := 0; i < n; i++ {
		c := alphabets[rand.Intn(l)]
		sb.WriteByte(c)
	}
	return sb.String()

}
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomOwner() string {
	return RandomString(6)
}
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
