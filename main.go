package main

import (
	"fmt"
	"os"
)

func main() {
	targetRootDir := os.Args[1]
	workingDir, _ := os.Getwd()

	runArgs := RunArgs{
		WorkingDir:    workingDir,
		TargetRootDir: targetRootDir,
	}

	runError := Run(runArgs)
	if runError != nil {
		fmt.Println(runError.Error())
		os.Exit(1)
	}
}
