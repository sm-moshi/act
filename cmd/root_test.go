package cmd

import (
	"context"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func configureTestInput(t *testing.T, input *Input) {
	t.Helper()
	cacheDir := t.TempDir()
	input.actionCachePath = path.Join(cacheDir, "act")
	input.cacheServerPath = path.Join(cacheDir, "actcache")
}

func configureRunnerTestInput(t *testing.T, input *Input) {
	t.Helper()
	configureTestInput(t, input)
	t.Setenv("DOCKER_CONFIG", t.TempDir())
	input.platforms = []string{"ubuntu-latest=node:16-buster-slim"}
	input.workdir = "../pkg/runner/testdata/"
	input.workflowsPath = "./basic/push.yml"
	input.noCacheServer = true
}

func skipIfRestrictedDockerRuntime(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}

	message := err.Error()
	if strings.Contains(message, "Docker daemon socket") ||
		strings.Contains(message, "error getting credentials") ||
		strings.Contains(message, "operation not permitted") {
		t.Skipf("skipping Docker-dependent test in restricted environment: %v", err)
	}
}

func TestReadSecrets(t *testing.T) {
	secrets := map[string]string{}
	ret := readEnvsEx(path.Join("testdata", "secrets.yml"), secrets, true)
	assert.True(t, ret)
	assert.Equal(t, `line1
line2
line3
`, secrets["MYSECRET"])
}

func TestReadEnv(t *testing.T) {
	secrets := map[string]string{}
	ret := readEnvs(path.Join("testdata", "secrets.yml"), secrets)
	assert.True(t, ret)
	assert.Equal(t, `line1
line2
line3
`, secrets["mysecret"])
}

func TestListOptions(t *testing.T) {
	input := &Input{}
	configureTestInput(t, input)
	rootCmd := createRootCommand(context.Background(), input, "")
	input.listOptions = true
	err := newRunCommand(context.Background(), input)(rootCmd, []string{})
	assert.NoError(t, err)
}

func TestRun(t *testing.T) {
	input := &Input{}
	rootCmd := createRootCommand(context.Background(), input, "")
	configureRunnerTestInput(t, input)
	err := newRunCommand(context.Background(), input)(rootCmd, []string{})
	skipIfRestrictedDockerRuntime(t, err)
	assert.NoError(t, err)
}

func TestRunPush(t *testing.T) {
	input := &Input{}
	rootCmd := createRootCommand(context.Background(), input, "")
	configureRunnerTestInput(t, input)
	err := newRunCommand(context.Background(), input)(rootCmd, []string{"push"})
	skipIfRestrictedDockerRuntime(t, err)
	assert.NoError(t, err)
}

func TestRunPushJsonLogger(t *testing.T) {
	input := &Input{}
	rootCmd := createRootCommand(context.Background(), input, "")
	configureRunnerTestInput(t, input)
	input.jsonLogger = true
	err := newRunCommand(context.Background(), input)(rootCmd, []string{"push"})
	skipIfRestrictedDockerRuntime(t, err)
	assert.NoError(t, err)
}

func TestFlags(t *testing.T) {
	for _, f := range []string{"graph", "list", "bug-report", "man-page"} {
		t.Run("TestFlag-"+f, func(t *testing.T) {
			input := &Input{}
			rootCmd := createRootCommand(context.Background(), input, "")
			configureRunnerTestInput(t, input)
			err := rootCmd.Flags().Set(f, "true")
			assert.NoError(t, err)
			err = newRunCommand(context.Background(), input)(rootCmd, []string{})
			if f == "bug-report" {
				skipIfRestrictedDockerRuntime(t, err)
			}
			assert.NoError(t, err)
		})
	}
}

func TestReadArgsFile(t *testing.T) {
	tables := []struct {
		path  string
		split bool
		args  []string
		env   map[string]string
	}{
		{
			path:  path.Join("testdata", "simple.actrc"),
			split: true,
			args:  []string{"--container-architecture=linux/amd64", "--action-offline-mode"},
		},
		{
			path:  path.Join("testdata", "env.actrc"),
			split: true,
			env: map[string]string{
				"FAKEPWD": "/fake/test/pwd",
				"FOO":     "foo",
			},
			args: []string{
				"--artifact-server-path", "/fake/test/pwd/.artifacts",
				"--env", "FOO=prefix/foo/suffix",
			},
		},
		{
			path:  path.Join("testdata", "split.actrc"),
			split: true,
			args:  []string{"--container-options", "--volume /foo:/bar --volume /baz:/qux --volume /tmp:/tmp"},
		},
	}
	for _, table := range tables {
		t.Run(table.path, func(t *testing.T) {
			for k, v := range table.env {
				t.Setenv(k, v)
			}
			args := readArgsFile(table.path, table.split)
			assert.Equal(t, table.args, args)
		})
	}
}
