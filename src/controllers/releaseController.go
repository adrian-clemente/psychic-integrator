package controllers

import "net/http"
import "components/web/page"
import "components/web/element"
import "fmt"
import "api/repository"
import "api/email"
import "api/project"

func ViewReleaseHandler(w http.ResponseWriter, r *http.Request) {
	projectName := r.URL.Query()["project"]
	repositoryName := repository.Repository(project.TEST_PROJECT)
	fmt.Println(projectName)
	if len(projectName) > 0 {
		repositoryName = repository.Repository(projectName[0])
	}

	//Retrieve the commits data that are going to be push to MASTER_BRANCH
	var releaseCommitsSections []element.CommitElement
	commits := repository.CommitDiff(repositoryName, repository.MASTER_BRANCH, repository.DEVELOP_BRANCH)
	for _, value := range commits {
		releaseCommitsSections = append(releaseCommitsSections, element.CommitElement{value.Hash,
				value.Author, value.Date, value.Text, value.JiraTicket})
	}

	var releaseProjects []element.OptionElement
	for _, project := range project.GetProjects() {
		releaseProjects = append(releaseProjects, element.OptionElement{ string(project.ProjectKey),
				string(project.ProjectName) })
	}

	var releaseTypes []element.OptionElement
	for _, releaseType := range project.GetReleaseTypes() {
		releaseTypes = append(releaseTypes, element.OptionElement{ string(releaseType.ReleaseTypeKey),
				string(releaseType.ReleaseTypeName) })
	}

	releasePage := page.ReleaseHandlerPage{releaseCommitsSections, releaseProjects, releaseTypes}
	releasePageContent := releasePage.GetContent();
	fmt.Fprintf(w, releasePageContent)
}

func PerformReleaseHandler(w http.ResponseWriter, r *http.Request) {
	projectName := r.FormValue("project")
	repositoryName := repository.Repository(projectName)
	releaseType := project.ReleaseTypeKey(r.FormValue("type"))

	repository.Clone(repositoryName)
	repository.Checkout(repositoryName, repository.DEVELOP_BRANCH)
	repository.ChangeBranch(repositoryName, repository.MASTER_BRANCH)

	project.IncrementVersionByReleaseType(repositoryName, releaseType)
	projectVersion := project.PrintVersion(repositoryName)

	//Retrieve the commits data that are going to be push to MASTER_BRANCH
	commits := repository.CommitDiff(repositoryName, repository.MASTER_BRANCH, repository.DEVELOP_BRANCH)
	var releaseCommitsSections []element.CommitElement
	for _, value := range commits {
		releaseCommitsSections = append(releaseCommitsSections, element.CommitElement{value.Hash,
			value.Author, value.Date, value.Text, value.JiraTicket})
	}
	//repository.Merge(repository.MASTER_BRANCH, repository.DEVELOP_BRANCH)

	//session := jira.Login()
	//jiraIssueKey := jira.CreateReleaseIssue(session)

	//repository.AddAll()
	//repository.Commit("\"Release of FAVOR version 1.3.2\"", jiraIssueKey)
	//repository.Push(repository.MASTER_BRANCH)

	//jira.CloseIssue(session, jiraIssueKey)
	//commits := repository.Log(20, repository.MASTER_BRANCH)

	email.GenerateReleaseEmail(project.FAVOR_PROJECT, projectVersion, commits)

	releasePage := page.ReleasePerformedPage{projectName, true}
	releasePageContent := releasePage.GetContent();
	fmt.Fprintf(w, releasePageContent)
}
