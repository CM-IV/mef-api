package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {

	rand.Seed(time.Now().UnixNano())

}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

//Generate random string
func RandomString(n int) string {

	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {

		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()

}

func RandomImage() string {

	images := []string{

		"https://i.ibb.co/5Rmp215/monero-min.jpg",
		"https://i.ibb.co/DrhHLJ6/anotha-one-min.jpg",
		"https://i.ibb.co/nz7mRHv/MONE-min.jpg",
		"https://i.ibb.co/0Jy36bX/Monerooo-min.jpg",
		"https://i.ibb.co/6RpBs7Q/Yon-monery-min.jpg",
	}

	n := len(images)
	return images[rand.Intn(n)]

}

func RandomOwner() string {

	return RandomString(10)

}

func RandomTitle() string {

	return RandomString(10)

}

func RandomSubtitle() string {

	return RandomString(8)

}

func RandomContent() string {

	return RandomString(30)

}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
