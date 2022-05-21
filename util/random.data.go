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

		"https://ik.imagekit.io/xbkhabiqcy9/img/gym_EUxo-U9Wn.webp?ik-sdk-version=javascript-1.4.3&updatedAt=1653050207758",
		"https://ik.imagekit.io/xbkhabiqcy9/img/monerochan-wildwest_14HYrFk5T.webp?ik-sdk-version=javascript-1.4.3&updatedAt=1653050207768",
		"https://ik.imagekit.io/xbkhabiqcy9/img/weddingdress_hHmMAUFSF.webp?ik-sdk-version=javascript-1.4.3&updatedAt=1653050207585",
		"https://ik.imagekit.io/xbkhabiqcy9/img/boat_IGucyrxms.webp?ik-sdk-version=javascript-1.4.3&updatedAt=1653050207903",
		"https://ik.imagekit.io/xbkhabiqcy9/img/naked_apron_1-1_nODYXGMGD.webp?ik-sdk-version=javascript-1.4.3&updatedAt=1653050242172",
		"https://ik.imagekit.io/xbkhabiqcy9/img/knighthood_1-1_4GJG8bIx3.webp?ik-sdk-version=javascript-1.4.3&updatedAt=1653050242402",
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
