package flagbuilder

func BuildFlagArg(flag, value string) []string {
	return []string{"-" + flag, value}
}
