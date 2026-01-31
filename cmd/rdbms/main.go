package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/migeru111/rdbms_go/pkg/executor"
	"github.com/migeru111/rdbms_go/pkg/parser"
	"github.com/migeru111/rdbms_go/pkg/storage"
	"github.com/migeru111/rdbms_go/pkg/types"
)

func main() {
	db := storage.NewMemoryStorage()
	exec := executor.NewExecutor(db)

	fmt.Println("Welcome to RDBMS-Go!")
	fmt.Println("Type SQL commands or 'exit' to quit.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("rdbms> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if strings.ToLower(input) == "exit" || strings.ToLower(input) == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		// Remove trailing semicolon if present
		input = strings.TrimSuffix(input, ";")

		p := parser.NewParser(input)
		stmt, err := p.Parse()
		if err != nil {
			fmt.Printf("Parse error: %s\n", err)
			continue
		}

		result, err := exec.Execute(stmt)
		if err != nil {
			fmt.Printf("Execution error: %s\n", err)
			continue
		}

		printResult(result)
	}
}

func printResult(result *types.Result) {
	if result.Message != "" {
		fmt.Println(result.Message)
		return
	}

	if len(result.Columns) == 0 {
		return
	}

	// Calculate column widths
	widths := make([]int, len(result.Columns))
	for i, col := range result.Columns {
		widths[i] = len(col)
	}
	for _, row := range result.Rows {
		for i, val := range row {
			if len(val) > widths[i] {
				widths[i] = len(val)
			}
		}
	}

	// Print header separator
	printSeparator(widths)

	// Print header
	fmt.Print("|")
	for i, col := range result.Columns {
		fmt.Printf(" %-*s |", widths[i], col)
	}
	fmt.Println()

	// Print header separator
	printSeparator(widths)

	// Print rows
	for _, row := range result.Rows {
		fmt.Print("|")
		for i, val := range row {
			fmt.Printf(" %-*s |", widths[i], val)
		}
		fmt.Println()
	}

	// Print footer separator
	printSeparator(widths)

	fmt.Printf("%d row(s)\n", len(result.Rows))
}

func printSeparator(widths []int) {
	fmt.Print("+")
	for _, w := range widths {
		fmt.Print(strings.Repeat("-", w+2) + "+")
	}
	fmt.Println()
}
