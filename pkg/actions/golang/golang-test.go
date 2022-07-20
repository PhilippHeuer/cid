package golang

import (
	"errors"
	"github.com/cidverse/cidverseutils/pkg/filesystem"
	"path/filepath"
	"strings"

	"github.com/cidverse/cid/pkg/common/api"
	"github.com/cidverse/cid/pkg/common/command"
	"github.com/cidverse/cid/pkg/repoanalyzer/analyzerapi"
	"github.com/rs/zerolog/log"
)

type TestActionStruct struct{}

// GetDetails retrieves information about the action
func (action TestActionStruct) GetDetails(ctx *api.ActionExecutionContext) api.ActionDetails {
	return api.ActionDetails{
		Name:             "golang-test",
		Version:          "0.1.0",
		UsedTools:        []string{"go"},
		ToolDependencies: GetToolDependencies(ctx),
	}
}

// Check evaluates if the action should be executed or not
func (action TestActionStruct) Check(ctx *api.ActionExecutionContext) bool {
	return ctx.CurrentModule != nil && ctx.CurrentModule.BuildSystem == analyzerapi.BuildSystemGoMod
}

// Execute runs the action
func (action TestActionStruct) Execute(ctx *api.ActionExecutionContext, state *api.ActionStateContext) error {
	// config
	coverageFile := filepath.Join(ctx.Paths.ArtifactModule(ctx.CurrentModule.Slug), "coverage.out")
	coverageJSON := filepath.Join(ctx.Paths.ArtifactModule(ctx.CurrentModule.Slug), "coverage.json")

	// test
	var testArgs []string
	testArgs = append(testArgs, `go test`)
	testArgs = append(testArgs, `-vet off`)
	testArgs = append(testArgs, `-cover`)
	testArgs = append(testArgs, `-coverprofile `+coverageFile)
	testArgs = append(testArgs, `-covermode=count`)
	testArgs = append(testArgs, `./...`)
	testResult := command.RunOptionalCommand(strings.Join(testArgs, " "), ctx.Env, ctx.CurrentModule.Directory)
	if testResult != nil {
		return errors.New("go unit tests failed. Cause: " + testResult.Error())
	}

	// get report
	covOut, covOutErr := command.RunCommandAndGetOutput("go tool cover -func "+coverageFile, ctx.Env, ctx.CurrentModule.Directory)
	if covOutErr != nil {
		return errors.New("failed to retrieve go coverage report. Cause: " + covOutErr.Error())
	}

	// parse report
	report := ParseCoverageProfile(covOut)
	log.Info().Float64("coverage", report.Percent).Msg("calculated final code coverage")

	// json report
	jsonOut, jsonOutErr := command.RunCommandAndGetOutput("go test -coverprofile "+coverageFile+" -covermode=count -json ./...", ctx.Env, ctx.CurrentModule.Directory)
	if jsonOutErr != nil {
		return errors.New("failed to retrieve go coverage report. Cause: " + covOutErr.Error())
	}
	_ = filesystem.SaveFileText(coverageJSON, jsonOut) //nolint:errcheck

	return nil
}

func init() {
	api.RegisterBuiltinAction(TestActionStruct{})
}
