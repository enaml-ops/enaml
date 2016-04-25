package diff

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// decodeYaml takes a reader to a YAML file and unmarshals it to the given struct.
func decodeYaml(r io.Reader, v interface{}) error {
	bytes, err := ioutil.ReadAll(r)
	if err == nil {
		yaml.Unmarshal(bytes, v)
	}
	return err
}
