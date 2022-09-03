package internal

import "github.com/Red-Sock/rscli/internal/config"

var commands = map[string]func(){
	config.Command: config.Run,
}
