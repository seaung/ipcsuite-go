package camera

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/seaung/ipcsuite-go/pkg/utils"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v3"
)

type NsePoc struct {
	Name   string            `yaml:"name"`
	Set    map[string]string `yaml:"set"`
	Rule   []RequestRule     `yaml:"rule"`
	Author Author            `yaml:"author"`
}

type RequestRule struct {
	Method     string            `yaml:"method"`
	Cache      string            `yaml:"cache"`
	Path       string            `yaml:"path"`
	Headers    map[string]string `yaml:"header"`
	Expression string            `yaml:"expression"`
}

type Author struct {
	Github       string   `yaml:"github"`
	VulnLink     []string `yaml:"vuln_link"`
	SolutionLink []string `yaml:"solution_link"`
	Description  string   `yaml:"description"`
	CVENumbr     []string `yaml:"cve"`
	Risk         string   `yaml:"risk"`
}

func loadYamlFile(filename string) (*NsePoc, error) {
	poc := &NsePoc{}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, poc)
	if err != nil {
		return nil, err
	}

	return poc, nil
}

func LoadNsePoc(filename string) (*NsePoc, error) {
	return loadYamlFile(filename)
}

func f(filename string) []string {
	return funk.UniqString(singleNsePoc(filename))
}

func singleNsePoc(filename string) []string {
	var files []string

	if strings.HasPrefix(filename, ".yml") || strings.HasPrefix(filename, ".yaml") {
		if utils.IsFileExists(filename) {
			files = append(files, filename)
		}
	}

	if strings.Contains(filename, "*") && strings.Contains(filename, "/") {
		absDirectory, _ := filepath.Abs(filename)
		baseName := filepath.Base(filename)
		ymlFiles := utils.GetFilenames(filepath.Dir(absDirectory), "yml")

		for _, file := range ymlFiles {
			base := filepath.Base(file)
			if len(baseName) == 1 && baseName == "*" {
				files = append(files, file)
				continue
			}

			if re, err := regexp.Compile(baseName); err != nil {
				if strings.Contains(file, baseName) {
					files = append(files, file)
				}
			} else {
				if re.MatchString(base) {
					files = append(files, file)
				}
			}
		}
	}

	return files
}

func LoadMutilNsePocs(filepath string) []*NsePoc {
	var pocs []*NsePoc
	filepath = strings.ReplaceAll(filepath, "\\", "/")
	for _, file := range f(filepath) {
		if poc, err := loadYamlFile(file); err == nil {
			pocs = append(pocs, poc)
		}
	}
	return pocs
}
