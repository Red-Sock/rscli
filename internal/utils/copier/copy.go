package copier

import (
	"encoding/json"
)

func Copy(src, dst interface{}) error {

	bts, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bts, dst)
	if err != nil {
		return err
	}

	return nil
}
