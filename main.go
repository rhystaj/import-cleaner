package main

import (
	"fmt"
	"importcleaner/run"
	"os"
)

func main() {
	targetRootDir := os.Args[1]
	workingDir, _ := os.Getwd()

	runArgs := run.RunArgs{
		WorkingDir:    workingDir,
		TargetRootDir: targetRootDir,
	}

	runError := run.Run(runArgs)
	if runError != nil {
		fmt.Println(runError.Error())
		os.Exit(1)
	}
}
