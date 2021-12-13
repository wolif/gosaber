package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Load(file string, conf interface{}) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	bs , err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bs, conf)
	if err != nil {
		return err
	}
	return nil
}
