package input

import (
	"bufio"
	"os"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

// GetAnswer получение ответа из консоли
func GetAnswer(s string) string {
	println(s)
	answ, _ := reader.ReadString('\n')
	return strings.Replace(answ, "\n", "", -1)
}
