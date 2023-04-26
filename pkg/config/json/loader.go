package json

import (
	"encoding/json"
	"os"
)

func Load(file string, conf interface{}) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.UseNumber()
	if err = dec.Decode(conf); err != nil {
		return err
	}
	return nil
}
