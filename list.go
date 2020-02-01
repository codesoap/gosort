package gosort

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// ListFunctionCallsWithinFile returns a map, where the keys are
// functions that are defined in the given file and the values are lists
// of functions in the same file, that are called by the function in the
// key. The values are in the order they are first called.
func ListFunctionCallsWithinFile(filename string) (map[string][]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return nil, err
	}

	var funcs = map[string][]string{}
	for _, decl := range f.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			caller := funcDecl.Name.Name
			funcs[caller] = extractFunctionCalls(funcDecl.Body.List...)
		}
	}
	return removeNonLocal(funcs), nil
}

func removeNonLocal(funcs map[string][]string) map[string][]string {
	var localFuncs = map[string][]string{}
	for caller, callees := range funcs {
		var localCallees = []string{}
		for _, callee := range callees {
			if _, ok := funcs[callee]; ok {
				localCallees = append(localCallees, callee)
			}
		}
		localFuncs[caller] = localCallees
	}
	return localFuncs
}

func extractFunctionCalls(stmts ...ast.Stmt) []string {
	var callees = []string{}
	for _, stmt := range stmts {
		switch v := stmt.(type) {
		case *ast.LabeledStmt:
			callees = appendNewFunctionCalls(callees, v.Stmt)
		case *ast.ExprStmt:
			callees = appendNewFunctionCall(callees, v.X)
		case *ast.SendStmt:
			callees = appendNewFunctionCall(callees, v.Value)
		case *ast.AssignStmt:
			for _, expr := range v.Rhs {
				callees = appendNewFunctionCall(callees, expr)
			}
		case *ast.GoStmt:
			callees = appendNewFunctionCall(callees, v.Call)
		case *ast.DeferStmt:
			callees = appendNewFunctionCall(callees, v.Call)
		case *ast.ReturnStmt:
			for _, expr := range v.Results {
				callees = appendNewFunctionCall(callees, expr)
			}
		case *ast.BlockStmt:
			callees = appendNewFunctionCalls(callees, v.List...)
		case *ast.IfStmt:
			callees = appendNewFunctionCalls(callees, v.Init)
			callees = appendNewFunctionCalls(callees, v.Body)
		case *ast.CaseClause:
			for _, expr := range v.List {
				callees = appendNewFunctionCall(callees, expr)
			}
			callees = appendNewFunctionCalls(callees, v.Body...)
		case *ast.SwitchStmt:
			callees = appendNewFunctionCalls(callees, v.Init)
			callees = appendNewFunctionCalls(callees, v.Body.List...)
		case *ast.TypeSwitchStmt:
			callees = appendNewFunctionCalls(callees, v.Init)
			callees = appendNewFunctionCalls(callees, v.Body.List...)
		case *ast.ForStmt:
			callees = appendNewFunctionCalls(callees, v.Init)
			callees = appendNewFunctionCall(callees, v.Cond)
			callees = appendNewFunctionCalls(callees, v.Post)
			callees = appendNewFunctionCalls(callees, v.Body)
		case *ast.RangeStmt:
			callees = appendNewFunctionCall(callees, v.X)
			callees = appendNewFunctionCalls(callees, v.Body)
		}
	}
	return callees
}

func appendNewFunctionCalls(callees []string, stmts ...ast.Stmt) []string {
	if stmts == nil {
		return callees
	}
	for _, callee := range extractFunctionCalls(stmts...) {
		if !contains(callees, callee) {
			callees = append(callees, callee)
		}
	}
	return callees
}

func appendNewFunctionCall(callees []string, expr ast.Expr) []string {
	switch v := expr.(type) {
	case *ast.CallExpr:
		switch fun := v.Fun.(type) {
		case *ast.SelectorExpr:
			callee := fun.Sel.Name
			if !contains(callees, callee) {
				callees = append(callees, callee)
			}
		case *ast.Ident:
			callee := fun.Name
			if !contains(callees, callee) {
				callees = append(callees, callee)
			}
		default:
			panic("Unhandeled CallExpr!")
		}
	}
	return callees
}
