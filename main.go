package main

import (
	"bufio"
	"fmt"
	"my-programming-language/ast"
	"my-programming-language/checker"
	"my-programming-language/evaluator"
	"my-programming-language/lexer"
	"my-programming-language/parser"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		runREPL(nil, nil, nil)
		return
	}

	filename := args[0]
	repl := false
	if len(args) > 1 && args[1] == "-i" {
		repl = true
	}

	typeEnv, valEnv, err := runFile(filename, checker.NewTypeEnv(nil), evaluator.NewEnv(nil), make(map[string]bool))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if repl {
		runREPL(typeEnv, valEnv, nil)
	}
}

func runFile(filename string, typeEnv *checker.TypeEnv, valEnv *evaluator.Env, imported map[string]bool) (*checker.TypeEnv, *evaluator.Env, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot resolve path '%s': %v", filename, err)
	}

	if imported[absPath] {
		return typeEnv, valEnv, nil // already imported
	}
	imported[absPath] = true

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot read file '%s': %v", filename, err)
	}

	tokens, err := lexer.New(string(data)).Tokenize()
	if err != nil {
		return nil, nil, err
	}

	prog, err := parser.New(tokens).ParseProgram()
	if err != nil {
		return nil, nil, err
	}

	// Process imports first
	dir := filepath.Dir(absPath)
	for _, decl := range prog.Declarations {
		if imp, ok := decl.(*ast.ImportExpr); ok {
			importPath := filepath.Join(dir, imp.Path)
			typeEnv, valEnv, err = runFile(importPath, typeEnv, valEnv, imported)
			if err != nil {
				return nil, nil, fmt.Errorf("import '%s': %v", imp.Path, err)
			}
		}
	}

	// Type check
	chk := checker.NewWithEnv(typeEnv)
	if err := chk.CheckProgram(prog); err != nil {
		return nil, nil, err
	}

	// Evaluate
	ev := evaluator.NewWithEnv(valEnv)
	if _, err := ev.EvalProgram(prog); err != nil {
		return nil, nil, err
	}

	return chk.Env(), ev.Env(), nil
}

func runREPL(typeEnv *checker.TypeEnv, valEnv *evaluator.Env, imported map[string]bool) {
	if typeEnv == nil {
		typeEnv = checker.NewTypeEnv(nil)
	}
	if valEnv == nil {
		valEnv = evaluator.NewEnv(nil)
	}

	chk := checker.NewWithEnv(typeEnv)
	ev := evaluator.NewWithEnv(valEnv)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("MEPL v0.1.0 - Type :quit to exit")

	for {
		fmt.Print("mepl> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == ":quit" || line == ":q" {
			break
		}

		// Try parsing as a program (declarations) first, then as expression
		result, err := evalInput(line, chk, ev)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			continue
		}
		if result != nil {
			fmt.Printf("= %s\n", result.String())
		}
	}
}

func evalInput(input string, chk *checker.Checker, ev *evaluator.Evaluator) (evaluator.Value, error) {
	tokens, err := lexer.New(input).Tokenize()
	if err != nil {
		return nil, err
	}

	// Try as declaration first (let/type at top level)
	if len(tokens) > 0 && (tokens[0].Literal == "let" || tokens[0].Literal == "type") {
		p := parser.New(tokens)
		prog, err := p.ParseProgram()
		if err == nil && len(prog.Declarations) > 0 {
			if err := chk.CheckProgram(prog); err != nil {
				return nil, err
			}
			val, err := ev.EvalProgram(prog)
			if err != nil {
				return nil, err
			}
			return val, nil
		}
	}

	// Parse as expression
	p := parser.New(tokens)
	expr, err := p.ParseSingleExpr()
	if err != nil {
		return nil, err
	}

	t, err := chk.Check(expr)
	if err != nil {
		return nil, err
	}
	_ = t

	val, err := ev.Eval(expr)
	if err != nil {
		return nil, err
	}

	return val, nil
}
