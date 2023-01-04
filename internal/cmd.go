package internal

func RunCMD(args []string) {
	if len(args) == 0 {
		println("no args given")
		return
	}
	plugins[args[0]].Run(args[1:])
}
