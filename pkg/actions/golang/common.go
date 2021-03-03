package golang

import (
	"github.com/EnvCLI/normalize-ci/pkg/common"
	"github.com/PhilippHeuer/cid/pkg/common/api"
	"github.com/PhilippHeuer/cid/pkg/common/command"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

// DetectGolangProject checks if the target directory is a go project
func DetectGolangProject(projectDir string) bool {
	// go.mod
	if _, err := os.Stat(projectDir+"/go.mod"); !os.IsNotExist(err) {
		log.Debug().Str("file", projectDir+"/go.mod").Msg("found go.mod")
		return true
	}

	return false
}

func crossCompile(env []string, goos string, goarch string) {
	buildAt := time.Now().UTC().Format(time.RFC3339)
	log.Info().Str("goos", goos).Str("goarch", goarch).Msg("running go build")

	fileExt := ""
	if goos == "windows" {
		fileExt = ".exe"
	}

	env = api.GetEffectiveEnv(env)
	env = common.SetEnvironment(env, "CGO_ENABLED", "false")
	env = common.SetEnvironment(env, "GOPROXY", "https://goproxy.io,direct")
	env = common.SetEnvironment(env, "GOOS", goos)
	env = common.SetEnvironment(env, "GOARCH", goarch)
	command.RunCommand(`go build -o `+Config.Paths.Artifact+`/`+goos+`_`+goarch+fileExt+` -ldflags "-s -w -X main.Version=`+common.GetEnvironmentOrDefault(env, "NCI_COMMIT_REF_RELEASE", "")+` -X main.CommitHash=`+common.GetEnvironmentOrDefault(env, "NCI_COMMIT_SHA_SHORT", "")+` -X main.BuildAt=`+buildAt+`" .`, env)
}
