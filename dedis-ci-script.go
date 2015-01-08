package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// get all changed files
	cmd := exec.Command("git", "diff", "--name-only", "origin/master")
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("error running git diff: ", err)
		os.Exit(1)
	}
	fnames := strings.Split(string(output), "\n")
	fmt.Println("Files Changed: ", fnames)
	good := true
	for _, fn := range fnames {
		// go source code must be properly formatted
		if strings.HasSuffix(fn, ".go") {
			fmtCmd := exec.Command("gofmt -l", fn)
			fmtCmd.Stderr = os.Stderr
			out, err := fmtCmd.Output()
			if err != nil {
				fmt.Println("Error Running go fmt: ", err)
				os.Exit(1)
			}
			// if go fmt returns anything that means the file has been
			// formatted and did not conform.
			if len(out) != 0 {
				fmt.Println("File Not Properly Formatted: ", fn)
				good = false
			}
		}
	}
	if good == false {
		fmt.Println("Failed: Files Not Properly Formatted")
		os.Exit(1)
	}
	tests := exec.Command("go", "test", "-v", "./...")
	tests.Stderr = os.Stderr
	tests.Stdout = os.Stdout
	err = tests.Run()
	if err != nil {
		fmt.Println("Tests Failed")
		os.Exit(1)
	}
}
