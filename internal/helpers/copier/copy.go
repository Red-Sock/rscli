package copier

import (
	"gopkg.in/yaml.v3"
)

func Copy(src, dst interface{}) error {
	bts, err := yaml.Marshal(src)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bts, dst)
	if err != nil {
		return err
	}

	return nil
}
