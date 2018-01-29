package authfile

import "fmt"

// Version build
const Version = "0.0.1.dev"

var (
	// Name of the program
	Name string
	// GitCommit hash
	GitCommit string
	//HumanVersion easy readable for Humans
	HumanVersion = fmt.Sprintf("%s v%s (%s)", Name, Version, GitCommit)
)
