package jira

import "net/http"
import "bytes"
import "fmt"
import "encoding/json"
import "io/ioutil"
import "api/config"
import "log"
import "time"

type JiraIssue struct {
	Fields JiraIssueFields `json:"fields"`
}

type JiraIssueFieldId struct {
	Id string `json:"id"`
}

type JiraIssueFieldName struct {
	Name string `json:"name"`
}

type JiraIssueFields struct {
	Project JiraIssueFieldId `json:"project"`
	Summary string `json:"summary"`
	IssueType JiraIssueFieldId `json:"issuetype"`
	Assignee JiraIssueFieldName `json:"assignee"`
	Priority JiraIssueFieldId `json:"priority"`
	Labels []string `json:"labels"`
	Description string `json:"description"`
}

type JiraIssueKey struct {
	Key string `json:"key"`
}

type JiraSession struct {
	Session SessionMap
}

type SessionMap struct {
	Value Session
}

type JiraVersionFields struct {
	Name        string `json:"name"`
	Released    bool `json:"released"`
	ReleaseDate string `json:"releaseDate"`
	Project     string `json:"project"`
}

type Session string

const JIRA_LOGIN_URL = "rest/auth/1/session"
const JIRA_ISSUE_RETRIEVE_URL = "rest/api/2/issue/%v?fields=summary"
const JIRA_ISSUE_CREATE_URL = "rest/api/2/issue"
const JIRA_ISSUE_UPDATE_URL = "rest/api/2/issue/%v"
const JIRA_ISSUE_TRANSITION_URL = "rest/api/2/issue/%v/transitions"
const JIRA_PROJECT_CREATION_URL = "rest/api/2/version"
const JIRA_ISSUE_BROWSE_URL = "browse/%v"

const GET = "GET"
const POST = "POST"
const PUT = "PUT"

func Login() Session {
	username := config.GetProperty("jira.auth.username")
	password := config.GetProperty("jira.auth.password")

	log.Printf("Login to JIRA with the user %v", username)

	var requestData = []byte(fmt.Sprintf(`{"username":"%v", "password": "%v"}`, username, password))
	loginUrl := getJiraUrl(JIRA_LOGIN_URL)

	body := doRequest(loginUrl, requestData, POST, "")
	var jsontype JiraSession
	json.Unmarshal(body, &jsontype)

	return jsontype.Session.Value
}

func GetJiraIssueBrowseUrl(issueKey string) string {
	return fmt.Sprintf(getJiraUrl(JIRA_ISSUE_BROWSE_URL), issueKey)
}

func RetrieveIssue(session Session, issueKey string) JiraIssueFields {

	var fields JiraIssueFields
	if issueKey != "" {
		log.Printf("Retrieving JIRA issue %v", issueKey)
		retrieveIssueUrl := fmt.Sprintf(getJiraUrl(JIRA_ISSUE_RETRIEVE_URL), issueKey)

		body := doRequest(retrieveIssueUrl, []byte(""), GET, string(session))

		var jsontype JiraIssue
		json.Unmarshal(body, &jsontype)
		fields = jsontype.Fields
	}

	return fields
}

func CreateReleaseIssue(session Session, project string) string {

	log.Println("Creating JIRA release issue")
	createIssueUrl := getJiraUrl(JIRA_ISSUE_CREATE_URL)

	projectId := JiraIssueFieldId{ config.GetProperty("jira.project.id") }
	summary := fmt.Sprintf("[RELEASE] %v", project)
	issueType := JiraIssueFieldId{ "3" }
	assignee := JiraIssueFieldName{ config.GetProperty("jira.assignee.name") }
	priority := JiraIssueFieldId{ "2" }
	labels := []string{"release"}
	description := fmt.Sprintf("Release of %v", project)

	jiraIssueFields := JiraIssueFields{projectId, summary, issueType, assignee, priority, labels, description }
	jiraIssue := JiraIssue {jiraIssueFields}

	jiraIssueJson, err := json.Marshal(jiraIssue)
	if err != nil {
		log.Print(err)
	}

	body := doRequest(createIssueUrl, jiraIssueJson, POST, string(session))
	var jsontype JiraIssueKey
	json.Unmarshal(body, &jsontype)

	return jsontype.Key
}

func CloseIssue(session Session, issueKey string) {
	log.Printf("Closing JIRA issue %v", issueKey)
	transitionIssueUrl := getJiraUrl(JIRA_ISSUE_TRANSITION_URL)
	doRequest(fmt.Sprintf(transitionIssueUrl, issueKey), []byte(`{
		"transition": {
			"id": "101"
		}
	}`), POST, string(session))
}

func UpdateIssueVersion(session Session, issueKey string, version string) {
	log.Printf("Update JIRA issue %v", issueKey)
	updateIssueUrl := getJiraUrl(JIRA_ISSUE_UPDATE_URL)
	doRequest(fmt.Sprintf(updateIssueUrl, issueKey), []byte(fmt.Sprintf(`{
		"fields": {
			"fixVersions": [{
				"name" : "%v"
			}]
		}
	}`, version)), PUT, string(session))
}

func CreateVersion(session Session, version string, project string) string {

	log.Printf("Creating JIRA release version for %v with version %v", project, version)
	versionUrl := getJiraUrl(JIRA_PROJECT_CREATION_URL)
	versionName := fmt.Sprintf("%v - %v", project, version)
	t := time.Now()
	dateFmt := fmt.Sprintf(t.Format("2006-01-02"))

	jiraVersionStruct := JiraVersionFields {versionName, true, dateFmt, "CS"}
	jiraVersionStructJson, _ := json.Marshal(jiraVersionStruct)

	doRequest(versionUrl, jiraVersionStructJson, POST, string(session))
	return versionName
}

func doRequest(requestUrl string, requestData []byte, requestType string, sessionToken string) []byte {
	req, err := http.NewRequest(requestType, requestUrl, bytes.NewBuffer(requestData))
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Content-Type", "application/json")

	if (len(sessionToken) > 0 ) {
		cookie := http.Cookie{Name:"JSESSIONID", Value: sessionToken}
		req.AddCookie(&cookie)
	}

	var body []byte
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Print(err)
	} else {
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
		}
	}
	return body;
}

func getJiraUrl(restUrl string) string {
	return config.GetProperty("jira.url.domain") + restUrl
}
