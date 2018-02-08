package authfile

import "fmt"

var (
	// Name of the program
	Name string
	// GitCommit hash
	GitCommit string
	// Version build
	Version string
	//HumanVersion easy readable for Humans
	HumanVersion = fmt.Sprintf("%s %s (%s)", Name, Version, GitCommit)
)
