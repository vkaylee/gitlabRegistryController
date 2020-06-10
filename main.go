package main

import (
	"github.com/vleedev/gitlabRegistryController/gitlabRegistry"
	"net/http"
)

func main() {
	gR := gitlabRegistry.GitlabRegistry{
		HttpClient:&http.Client{},
	}
	gR.Run()
}
