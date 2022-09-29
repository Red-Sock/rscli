package randomizer

import "math/rand"

var goodbuys = []string{
	`was a pleasure to work with you!`,
	`see ya!`,
}

func GoodGoodBuy() string {
	return goodbuys[rand.Intn(len(goodbuys))]
}
