package utils

import (
	"bufio"
	"os"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func GetAnswer(s string) string {
	println(s)
	answ, _ := reader.ReadString('\n')
	return strings.Replace(answ, "\n", "", -1)
}
