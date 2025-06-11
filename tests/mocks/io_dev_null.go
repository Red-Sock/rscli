package mocks

import (
	"github.com/Red-Sock/rscli/internal/io/colors"
)

type IoDevNul struct {
}

func (i IoDevNul) Println(_ ...string) {

}

func (i IoDevNul) Print(_ string) {
}

func (i IoDevNul) PrintlnColored(_ colors.Color, _ ...string) {
}

func (i IoDevNul) PrintColored(_ colors.Color, _ string) {
}

func (i IoDevNul) Error(_ string) {
}

func (i IoDevNul) GetInput() (string, error) {
	return "", nil
}

func (i IoDevNul) GetInputOneOf(options []string) string {
	return ""
}
