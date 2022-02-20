package main

import (
	"fmt"
	"strings"

	"github.com/golang-demos/chalk"
)

func getIndent(depth int) string {
	return strings.Repeat("    ", depth)
}

func info(depth int, args ...interface{}) {
	fmt.Print(getIndent(depth))
	fmt.Print(chalk.BlueLight().Bold())
	fmt.Print("[INFO]")
	fmt.Print(chalk.Reset())
	fmt.Print(" ")
	fmt.Println(args...)
}

func warn(depth int, args ...interface{}) {
	fmt.Print(getIndent(depth))
	fmt.Print(chalk.YellowLight().Bold())
	fmt.Print("[WARN]")
	fmt.Print(chalk.Reset())
	fmt.Print(" ")
	fmt.Println(args...)
}

func success(depth int, args ...interface{}) {
	fmt.Print(getIndent(depth))
	fmt.Print(chalk.GreenLight().Bold())
	fmt.Print("[SUCCESS]")
	fmt.Print(chalk.Reset())
	fmt.Print(" ")
	fmt.Println(args...)
}

func error(depth int, args ...interface{}) {
	fmt.Print(getIndent(depth))
	fmt.Print(chalk.RedLight().Bold())
	fmt.Print("[ERROR]")
	fmt.Print(chalk.Reset())
	fmt.Print(" ")
	fmt.Println(args...)
}