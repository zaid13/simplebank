package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklomnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().Unix())

}

func RandomInt(min, max int64) int64 {

	return min + rand.Int63n(max-min+1) //0->max-min
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

//RandomOwner generates random owner name
func RandomOwner() string {
	return RandomString(6)
}

//RandomMoeny generates random amount of money

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

//RandomCurrency generates random amount currency

func RandomCurrency() string {

	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]

}

//Random Email gernates random email

func RandomEmail() string {

	return fmt.Sprintf(RandomString(6), "@", RandomString(4), ".", RandomString(3))
}
