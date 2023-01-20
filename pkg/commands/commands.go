package commands

func init() {
	//var err error
	//rsCLI, err = os.Executable()
	//if err != nil {
	//	panic(err)
	//}
	//_, rsCLI = path.Split(rsCLI)
}

func RsCLI() string {
	return rsCLI
}

var rsCLI string

const (
	GetUtil = "get"
	Delete  = "del"
)
