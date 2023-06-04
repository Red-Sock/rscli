package randomizer

import "math/rand"

var goodbuys = []string{
	`Was a pleasure to work with you!`,
	`See ya!`,
}

func GoodGoodBuy() string {
	return goodbuys[rand.Intn(len(goodbuys))]
}
