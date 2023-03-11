package update

import "github.com/Red-Sock/rscli/plugins/project/processor"

func Do(p processor.Project) error {

}

var updates = []func(p processor.Project) error{
	v009
}
