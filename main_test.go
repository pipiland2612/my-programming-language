package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func runCLI(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	cmd.Dir = "."
	cmd.Env = append(cmd.Environ(), "GOCACHE="+t.TempDir())

	out, err := cmd.Output()
	stdout := string(out)

	stderr := ""
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			stderr = string(ee.Stderr)
		}
	}

	return stdout, stderr, err
}

func TestCLI_PrintingExample_BlackBox(t *testing.T) {
	stdout, stderr, err := runCLI(t, filepath.Join("examples", "printing.mepl"))
	require.NoError(t, err, "expected success; stderr:\n%s", stderr)

	expected := strings.Join([]string{
		"42",
		"true",
		"false",
		"Hello, World!",
		"(1, true)",
		"[1, 2, 3]",
		"inl 5",
		"inr true",
		"([1, 2], hello)",
		"(1, (2, (3, 4)))",
		"<function>",
		"Hello World!",
		"",
	}, "\n")

	require.Equal(t, expected, stdout)
	require.Empty(t, stderr)
}

func TestCLI_TypeErrorExample_BlackBox(t *testing.T) {
	_, stderr, err := runCLI(t, filepath.Join("examples", "integers-errors.mepl"))
	require.Error(t, err, "expected non-zero exit for type error example")
	require.Contains(t, stderr, "type error")
	require.Contains(t, stderr, "operator '+' expects Int")
}
