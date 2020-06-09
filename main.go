package main

import (
	"github.com/vleedev/gitlabRegistryController/gitlabRegistry"
)

func main() {
	gR := gitlabRegistry.GitlabRegistry{}
	gR.Run()
}
