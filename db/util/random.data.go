package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {

	rand.Seed(time.Now().UnixNano())

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

func RandomTitle() string {

	return RandomString(10)

}

func RandomSubtitle() string {

	return RandomString(8)

}

func RandomContent() string {

	return RandomString(30)

}
