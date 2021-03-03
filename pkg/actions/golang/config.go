package golang

import (
	"github.com/PhilippHeuer/cid/pkg/common/api"
	"github.com/jinzhu/configor"
)

var Config = struct {
	Paths api.PathConfig
	GoLang struct {
		Platform []struct {
			Goos string `required:"true"`
			Goarch string `required:"true"`
		}
	}
}{}

func loadConfig(projectDirectory string) {
	configor.Load(&Config, projectDirectory + "/mpi.yml")
}