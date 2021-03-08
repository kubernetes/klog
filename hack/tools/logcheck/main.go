/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"go/ast"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
	"golang.org/x/tools/go/packages"
)

// Doc explaining the tool.
const Doc = "Tool to check use of unstructured logging patterns."

// Analyzer runs static analysis.
var Analyzer = &analysis.Analyzer{
	Name: "logcheck",
	Doc:  Doc,
	Run:  run,
}

func main() {
	handleErrors()
	singlechecker.Main(Analyzer)
}

func run(pass *analysis.Pass) (interface{}, error) {

	for _, file := range pass.Files {

		ast.Inspect(file, func(n ast.Node) bool {

			// We are intrested in function calls, as we want to detect klog.* calls
			// passing all function calls to checkForFunctionExpr
			if fexpr, ok := n.(*ast.CallExpr); ok {

				checkForFunctionExpr(fexpr.Fun, pass)
			}

			return true
		})
	}
	return nil, nil
}

// checkForFunctionExpr checks for unstructured logging function, prints error if found any.
func checkForFunctionExpr(fun ast.Expr, pass *analysis.Pass) {

	/* we are extracting external package function calls e.g. klog.Infof fmt.Printf
	   and eliminating calls like setLocalHost()
	   basically function calls that has selector expression like .
	*/
	if selExpr, ok := fun.(*ast.SelectorExpr); ok {
		// extracting function Name like Infof
		fName := selExpr.Sel.Name

		// for nested function cases klog.V(1).Infof scenerios
		// if selExpr.X contains one more caller expression which is selector expression
		// we are extracting klog and discarding V(1)
		if n, ok := selExpr.X.(*ast.CallExpr); ok {
			if _, ok = n.Fun.(*ast.SelectorExpr); ok {
				selExpr = n.Fun.(*ast.SelectorExpr)
			}
		}

		// extracting package name
		pName, ok := selExpr.X.(*ast.Ident)

		// Matching if package name is klog and any unstructured logging function is used.
		if ok && pName.Name == "klog" && isUnstructured((fName)) {

			msg := fmt.Sprintf("unstructured logging function %q should not be used", fName)
			pass.Report(analysis.Diagnostic{
				Pos:     fun.Pos(),
				Message: msg,
			})
		}
	}
}

func isUnstructured(fName string) bool {

	// List of klog functions we do not want to use after migration to structured logging.
	unstrucured := []string{
		"Infof", "Info", "Infoln", "InfoDepth",
		"Warning", "Warningf", "Warningln", "WarningDepth",
		"Error", "Errorf", "Errorln", "ErrorDepth",
		"Fatal", "Fatalf", "Fatalln", "FatalDepth",
	}

	for _, name := range unstrucured {
		if fName == name {
			return true
		}
	}

	return false
}

// handleErrors exits with code 0 for "no Go files" errors
// and code 1 for all other errors while loading packages
// that match provided pattern.
func handleErrors() {
	pkgs, err := packages.Load(nil, os.Args[1:]...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// count of "no Go files" errors
	goFileErrorsCount := 0
	// count of errors excluding "no Go files" errors
	otherErrorsCount := 0
	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		for _, err := range pkg.Errors {
			fmt.Fprintln(os.Stderr, err)
			if strings.Contains(err.Msg, "no Go files") {
				goFileErrorsCount++
				continue
			}
			otherErrorsCount++
		}
	})
	if otherErrorsCount > 0 { // if there are errors other than "no Go files" errors, exit with code 1
		os.Exit(1)
	} else if otherErrorsCount == 0 && goFileErrorsCount > 0 { // if there are only "no Go files" errors, exit with code 0
		os.Exit(0)
	}
}
