package types

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type NsePocs struct {
	Description string            `yaml:"description"`
	Author      string            `yaml:"author"`
	Host        string            `yaml:"host"`
	Port        int               `yaml:"port"`
	Path        string            `yaml:"path"`
	StatusCode  int               `yaml:"status_code"`
	Headers     map[string]string `yaml:"headers"`
	Payloads    string            `yaml:"payloads"`
	Regex       string            `yaml:"regex"`
}

func LoaderSingePoc(path string) (*NsePocs, error) {
	return loaderPoc(path)
}

func loaderPoc(path string) (*NsePocs, error) {
	nse := &NsePocs{}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, nse)
	if err != nil {
		return nil, err
	}
	return nse, nil
}

func CheckSingleVuln(filepath string) {
	poc, err := LoaderSingePoc(filepath)
	if err != nil {
		return
	}

	runPoc(poc)
}

func runPoc(poc *NsePocs) {}
