package main

import (
	"fmt"
	"github.com/codesoap/gosort"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Give one file as an argument.")
		os.Exit(1)
	}
	funcs, err := gosort.ListFunctionCallsWithinFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
	tree := &gosort.FuncTree{}
	tree.SetValue(funcs)
	tree.BranchOutByTopology()
	tree.BranchOutByCallOrder([]string{})
	printOrder(tree)
}

func printOrder(tree *gosort.FuncTree) {
	if tree.Left != nil {
		printOrder(tree.Left)
	}
	if tree.Right != nil {
		printOrder(tree.Right)
	}
	if tree.Value() != nil {
		for _, f := range tree.Order() {
			fmt.Println(f)
		}
	}
}
