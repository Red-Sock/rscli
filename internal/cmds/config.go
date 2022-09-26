package cmds

import (
	"bufio"
	"github.com/Red-Sock/rscli/internal/randomizer"
	"github.com/Red-Sock/rscli/internal/service/config"
	"os"
	"strings"
)

func RunConfig(args []string) {
	c, err := config.Run(args)
	if err != nil {
		println(err.Error())
		return
	}

	err = c.TryWrite()
	if err != nil {
		if err != os.ErrExist {
			println(err.Error())
			return
		}

		reader := bufio.NewReader(os.Stdin)
		println("file " + c.GetPath() + " already exists. Do you want to override it? Y(es)/N(o)")
		anws, _ := reader.ReadString('\n')
		anws = strings.ToLower(strings.Replace(anws, "\n", "", -1))

		if strings.HasPrefix(anws, "n") {
			println("aborting config creation. " + randomizer.GoodGoodBuy())
		}
	}
}
