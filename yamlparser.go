package translate

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// LoadYaml loads a language from a yaml data source
func LoadYaml(lang io.Reader, id string) (*Language, error) {

	items := map[string]Value{}

	data, err := ioutil.ReadAll(lang)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &items)
	if err != nil {
		return nil, err
	}

	return &Language{ID: id, Keys: items}, nil
}
