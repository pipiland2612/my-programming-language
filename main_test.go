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

func TestCLI_PrintingExample(t *testing.T) {
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

func TestCLI_TypeErrorExample(t *testing.T) {
	_, stderr, err := runCLI(t, filepath.Join("examples", "integers-errors.mepl"))
	require.Error(t, err, "expected non-zero exit for type error example")
	require.Contains(t, stderr, "type error")
	require.Contains(t, stderr, "operator '+' expects Int")
}

func TestCLI_TuplesExample(t *testing.T) {
	stdout, stderr, err := runCLI(t, filepath.Join("examples", "tuples.mepl"))
	require.NoError(t, err, "expected success; stderr:\n%s", stderr)

	require.Contains(t, stdout, "(1, true, hello)")
	require.Contains(t, stdout, "1\ntrue\nhello\n")
	require.Contains(t, stdout, "(10, 20, 30, 40)")
	require.Contains(t, stdout, "40\n")
}

func TestCLI_TuplesErrorExample(t *testing.T) {
	_, stderr, err := runCLI(t, filepath.Join("examples", "tuples-errors.mepl"))
	require.Error(t, err, "expected non-zero exit for type error")
	require.Contains(t, stderr, "type error")
	require.Contains(t, stderr, "out of bounds")
}

func TestCLI_RecordsExample(t *testing.T) {
	stdout, stderr, err := runCLI(t, filepath.Join("examples", "records.mepl"))
	require.NoError(t, err, "expected success; stderr:\n%s", stderr)

	require.Contains(t, stdout, "{name = Alice, age = 30}")
	require.Contains(t, stdout, "Alice\n")
	require.Contains(t, stdout, "30\n")
	require.Contains(t, stdout, "on\n")
}

func TestCLI_RecordsErrorExample(t *testing.T) {
	_, stderr, err := runCLI(t, filepath.Join("examples", "records-errors.mepl"))
	require.Error(t, err, "expected non-zero exit for type error")
	require.Contains(t, stderr, "type error")
	require.Contains(t, stderr, "no field")
}

func TestCLI_LoopsExample(t *testing.T) {
	stdout, stderr, err := runCLI(t, filepath.Join("examples", "loops.mepl"))
	require.NoError(t, err, "expected success; stderr:\n%s", stderr)

	// Basic loop: 0,1,2,3,4
	require.Contains(t, stdout, "0\n1\n2\n3\n4\n")
	// Squares: 1,4,9,16,25
	require.Contains(t, stdout, "1\n4\n9\n16\n25\n")
}

func TestCLI_StringsExample(t *testing.T) {
	stdout, stderr, err := runCLI(t, filepath.Join("examples", "strings.mepl"))
	require.NoError(t, err, "expected success; stderr:\n%s", stderr)

	require.Contains(t, stdout, "Hello, World!")
	require.Contains(t, stdout, "5\n") // length
	require.Contains(t, stdout, "H\n") // charAt 0
	require.Contains(t, stdout, "o\n") // charAt 4
	require.Contains(t, stdout, "MEPL\n")
}

func TestCLI_ADTExample(t *testing.T) {
	stdout, stderr, err := runCLI(t, filepath.Join("examples", "adt.mepl"))
	require.NoError(t, err, "expected success; stderr:\n%s", stderr)

	// Color enum
	require.Contains(t, stdout, "red\ngreen\nblue\n")
	// Option type
	require.Contains(t, stdout, "0\n42\n")
	// Shape area: Circle 5 -> 75, Rectangle (3,4) -> 12
	require.Contains(t, stdout, "75\n12\n")
	// Recursive Expr: Add(Lit 10, Neg(Lit 3)) -> 7
	require.Contains(t, stdout, "7\n")
	// Equality
	require.Contains(t, stdout, "true\nfalse\ntrue\nfalse\ntrue\n")
}

func TestCLI_AllExamplesRun(t *testing.T) {
	// Verify every non-error example runs without error
	examples := []string{
		"variables.mepl", "integers.mepl", "booleans.mepl",
		"functions.mepl", "pairs.mepl", "lists.mepl", "sums.mepl",
		"recursion.mepl", "declarations.mepl", "imports.mepl",
		"comments.mepl", "printing.mepl", "strings.mepl",
		"tuples.mepl", "records.mepl", "loops.mepl", "adt.mepl",
	}
	for _, ex := range examples {
		t.Run(ex, func(t *testing.T) {
			_, stderr, err := runCLI(t, filepath.Join("examples", ex))
			require.NoError(t, err, "example %s failed; stderr:\n%s", ex, stderr)
		})
	}
}

func TestCLI_AllErrorExamplesReject(t *testing.T) {
	// Verify every error example fails with a type error
	examples := []string{
		"variables-errors.mepl", "integers-errors.mepl", "booleans-errors.mepl",
		"functions-errors.mepl", "pairs-errors.mepl",
		"lists-errors.mepl", "sums-errors.mepl",
		"recursion-errors.mepl", "declarations-errors.mepl",
		"imports-errors.mepl", "comments-errors.mepl", "printing-errors.mepl",
		"strings-errors.mepl", "tuples-errors.mepl", "records-errors.mepl", "loops-errors.mepl",
		"adt-errors.mepl",
	}
	for _, ex := range examples {
		t.Run(ex, func(t *testing.T) {
			_, stderr, err := runCLI(t, filepath.Join("examples", ex))
			require.Error(t, err, "expected error for %s", ex)
			require.Contains(t, stderr, "type error", "expected 'type error' in stderr for %s: %s", ex, stderr)
		})
	}
}
