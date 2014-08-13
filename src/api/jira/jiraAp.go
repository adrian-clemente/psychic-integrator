package jira

import "net/http"
import "bytes"
import "fmt"
import "encoding/json"
import "io/ioutil"
import "api/config"
import "log"

type JiraIssue struct {
	Fields JiraIssueFields
}

type JiraIssueFields struct {
	Summary string
	Key string
}

type JiraSession struct {
	Session SessionMap
}

type SessionMap struct {
	Value Session
}

type Session string

const JIRA_LOGIN_URL = "rest/auth/1/session"
const JIRA_ISSUE_RETRIEVE_URL = "rest/api/2/issue/%v?fields=summary"
const JIRA_ISSUE_CREATE_URL = "rest/api/2/issue"
const JIRA_ISSUE_TRANSITION_URL = "rest/api/2/issue/%v/transitions"
const JIRA_ISSUE_BROWSE_URL = "browse/%v"

const GET = "GET"
const POST = "POST"

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

func RetrieveIssue(issueKey string) JiraIssueFields {

	log.Printf("Retrieving JIRA issue %v", issueKey)

	session := Login()
	retrieveIssueUrl := fmt.Sprintf(getJiraUrl(JIRA_ISSUE_RETRIEVE_URL), issueKey)

	body := doRequest(retrieveIssueUrl, []byte(""), GET, string(session))
	var jsontype JiraIssue

	fmt.Println(string(body))
	json.Unmarshal(body, &jsontype)

	return jsontype.Fields
}

func CreateReleaseIssue(session Session) string {

	log.Println("Creating JIRA release issue")
	createIssueUrl := getJiraUrl(JIRA_ISSUE_CREATE_URL)

	body := doRequest(createIssueUrl, []byte(`
		{"fields":
			{"project": {
				"id": "11497"
			},
			"summary": "Release",
			"issuetype": {
				"id": "3"
			},
			"assignee": {
            	"name": "customersupport"
        	},
			"priority": {
				"id": "2"
			},
			"labels": [
				"release"
			],
			"description": "Release"
		}
	}`), POST, string(session))

	var jsontype JiraIssueFields
	json.Unmarshal(body, &jsontype)

	return jsontype.Key
}

func CloseIssue(session Session, issueKey string) {

	log.Println("Closing JIRA issue %v", issueKey)
	transitionIssueUrl := getJiraUrl(JIRA_ISSUE_TRANSITION_URL)

	doRequest(fmt.Sprintf(transitionIssueUrl, issueKey), []byte(`{
		"transition": {
			"id": "101"
		}
	}`), POST, string(session))
}

func doRequest(requestUrl string, requestData []byte, requestType string, sessionToken string) []byte {
	req, _ := http.NewRequest(requestType, requestUrl, bytes.NewBuffer(requestData))
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
