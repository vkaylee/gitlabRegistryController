package gitlabRegistry

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type GitlabRegistry struct {
	Domain *string
	AuthToken *string
	NameSpace *string
	ProjectName *string
	ProjectId *string
	RepoId *int
	BaseUrl *string
	RepoTagUrl *string
	SpecificTag *string
	Regex *string
	NumToHold *int
}
type Message struct {
	ID int `json:"id"`
}
type DeleteWithRegex struct {
	Regex    string `json:"name_regex"`
	NumToHold int `json:"keep_n"`
}
// setDomain sets the domain field's value.
func (g *GitlabRegistry) setDomain(v *string) {
	g.Domain = v
}
// setAuthToken sets the authToken field's value.
func (g *GitlabRegistry) setAuthToken(v *string) {
	g.AuthToken = v
}
// setNameSpace sets the nameSpace field's value.
func (g *GitlabRegistry) setNameSpace(v *string) {
	g.NameSpace = v
}
// setProjectName sets the projectName field's value.
func (g *GitlabRegistry) setProjectName(v *string) {
	g.ProjectName = v
}
// setSpecificTag sets the specificTag field's value.
func (g *GitlabRegistry) setSpecificTag(v *string) {
	g.SpecificTag = v
}
// setNumToHold sets the hold field's value.
func (g *GitlabRegistry) setNumToHold(v *int) {
	g.NumToHold = v
}
// setRegex sets the regex field's value.
func (g *GitlabRegistry) setRegex(v *string) {
	g.Regex = v
}
func (g *GitlabRegistry) setProjectId() {
	id := fmt.Sprintf("%s%%2F%s",*g.NameSpace,*g.ProjectName)
	g.ProjectId = &id
}
func (g *GitlabRegistry) prepare() {
	g.parseInput()
	g.setProjectId()
	g.generateBaseUrl()
	g.getRepoId()
	g.generateRepositoryTagUrl()
}
func (g *GitlabRegistry) Run() {
	g.prepare()
	if g.SpecificTag != nil {
		g.deleteSpecificTag()
	}
	if g.Regex != nil {
		g.deleteWithRegex()
	}
}
func (g *GitlabRegistry) generateBaseUrl() {
	baseUrl := fmt.Sprintf("%s/projects/%s/registry/repositories", *g.Domain, *g.ProjectId)
	g.BaseUrl = &baseUrl
}
func (g *GitlabRegistry) generateRepositoryTagUrl() {
	repoTagUrl := fmt.Sprintf("%s/%d/tags", *g.BaseUrl, *g.RepoId)
	g.RepoTagUrl = &repoTagUrl
}
func (g *GitlabRegistry) deleteSpecificTag() {
	client := &http.Client{}
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/%s", *g.RepoTagUrl, *g.SpecificTag),
		nil,
		)
	g.failOnError(err, "Error setting http request")
	req.Header.Add("PRIVATE-TOKEN", *g.AuthToken)
	resp, err := client.Do(req)
	g.failOnError(err, "Error deleting url")
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		log.Printf("%s has just been deleted successfully", *g.SpecificTag)
	} else {
		log.Printf("%s was not deleted.", *g.SpecificTag)
	}
}
func (g *GitlabRegistry) deleteWithRegex() {
	client := &http.Client{}
	values := DeleteWithRegex{
		Regex:     *g.Regex,
		NumToHold: *g.NumToHold,
	}
	jsonValue, _ := json.Marshal(values)
	req, err := http.NewRequest(http.MethodDelete,*g.RepoTagUrl,bytes.NewBuffer(jsonValue))
	g.failOnError(err, "Error setting http request")
	req.Header.Add("PRIVATE-TOKEN", *g.AuthToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	g.failOnError(err, "Error deleting url")
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		log.Printf("The request to delete with %s regex has just been done successfully", *g.Regex)
	} else {
		log.Printf("The request to delete with %s regex has not been done yet", *g.Regex)
		body, err := ioutil.ReadAll(resp.Body)
		g.failOnError(err, "Error reading json body")
		fmt.Println(string(body))
	}
}
func (g *GitlabRegistry) getRepoId() {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, *g.BaseUrl, nil)
	g.failOnError(err, "Error setting http request")
	req.Header.Add("PRIVATE-TOKEN", *g.AuthToken)
	resp, err := client.Do(req)
	g.failOnError(err, "Error getting BaseUrl")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	g.failOnError(err, "Error reading all")
	fmt.Println(string(body))
	msgs := make([]Message,0)
	err = json.Unmarshal(body, &msgs)
	g.failOnError(err, "Error encode json")
	g.RepoId = &msgs[0].ID
}
func (g *GitlabRegistry) failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func (g *GitlabRegistry) msgAndExit(msg string) {
	log.Fatalf("Error: %s", msg)
}
func (g *GitlabRegistry) parseInput() {
	domain := flag.String("domain", os.Getenv("CI_API_V4_URL"), "a base url of your gitlab with api version, ex: https://gitlab.example.com/api/v4")
	authToken := flag.String("authToken", os.Getenv("AUTH_TOKEN"), "a token that is used to auth with gitlab")
	nameSpace := flag.String("nameSpace", os.Getenv("CI_PROJECT_NAMESPACE"), "a namespace of your project")
	projectName := flag.String("projectName", os.Getenv("CI_PROJECT_NAME"), "a project name of your project")
	specificTag := flag.String("specificTag", "", "a image tag that you want to delete")
	regex := flag.String("regex", "", "a regex pattern to match all images")
	numToHold := flag.Int("hold", 3, "a volume of images that you want to keep from latest to older")

	flag.Parse()

	if *domain == "" {
		g.msgAndExit(fmt.Sprintf("Please input with -%s or set %s", "domain", "CI_API_V4_URL"))
	}
	if *authToken == "" {
		g.msgAndExit(fmt.Sprintf("Please input with -%s or set %s", "authToken", "CI_REGISTRY_PASSWORD"))
	}
	if *nameSpace == "" {
		g.msgAndExit(fmt.Sprintf("Please input with -%s or set %s", "nameSpace", "CI_PROJECT_NAMESPACE"))
	}
	if *projectName == "" {
		g.msgAndExit(fmt.Sprintf("Please input with -%s or set %s", "projectName", "CI_PROJECT_NAME"))
	}

	if *regex != "" && *specificTag != "" {
		log.Fatalf("Option -regex and -specificTag can not go together!")
	}

	if *specificTag != "" {
		g.setSpecificTag(specificTag)
	}
	if *regex != "" {
		g.setRegex(regex)
	}
	g.setDomain(domain)
	g.setAuthToken(authToken)
	g.setNameSpace(nameSpace)
	g.setProjectName(projectName)
	g.setNumToHold(numToHold)
}