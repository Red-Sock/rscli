package internal

import "os"

func RunCMD(args map[string][]string) {
	if len(args) == 0 {
		println("no args given")
		os.Exit(1)
	}
	// todo
	//err := plugins[args[0]].Run(args[1:])
	//if err != nil {
	//	println(err.Error())
	//	os.Exit(1)
	//}
}
