package cmds

import (
	"os"
	"strings"

	"github.com/Red-Sock/rscli/pkg/service/project/config-processor/config"

	"github.com/Red-Sock/rscli/internal/randomizer"
	"github.com/Red-Sock/rscli/internal/utils/input"
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

		answ := input.GetAnswer("file " + c.GetPath() + " already exists. Do you want to override it? Y(es)/N(o)")
		answ = strings.ToLower(strings.Replace(answ, "\n", "", -1))

		if !strings.HasPrefix(answ, "y") {
			println("aborting config creation. " + randomizer.GoodGoodBuy())
		}

		err = c.ForceWrite()
		if err != nil {
			println(err.Error())
			return
		}
	}
	println("successfully create config at " + c.GetPath() + ". " + randomizer.GoodGoodBuy())
}
