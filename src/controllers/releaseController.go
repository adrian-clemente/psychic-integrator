package controllers

import "net/http"
import "components/web/pages"
import "components/web/sections/release"
import "fmt"
import "api/repository"
import "api/email"
import "api/project"
import "api/jenkins"

func ViewReleaseHandler(w http.ResponseWriter, r *http.Request) {

	jenkins.BuildProject()

	repositoryName := repository.Repository(project.TEST_PROJECT)

	//Retrieve the commits data that are going to be push to MASTER_BRANCH
	var releaseCommitsSections []release.ReleaseCommitSection
	commits := repository.CommitDiff(repositoryName, repository.MASTER_BRANCH, repository.DEVELOP_BRANCH)
	for _, value := range commits {
		releaseCommitsSections = append(releaseCommitsSections, release.ReleaseCommitSection{value.Hash,
				value.Author, value.Date, value.Text, value.JiraTicket})
	}

	var releaseProjects []release.ReleaseProjectSection
	for _, project := range project.GetProjects() {
		releaseProjects = append(releaseProjects, release.ReleaseProjectSection{ string(project.ProjectKey),
				string(project.ProjectName) })
	}

	releasePage := pages.ReleaseHandlerPage{releaseCommitsSections, releaseProjects}
	releasePageContent := releasePage.GetContent();
	fmt.Fprintf(w, releasePageContent)
}

func PerformReleaseHandler(w http.ResponseWriter, r *http.Request) {

	projectName := r.FormValue("project")
	//releaseType := r.FormValue("type")
	repositoryName := repository.Repository(projectName)


	repository.Clone(repositoryName)
	repository.Checkout(repositoryName, repository.DEVELOP_BRANCH)
	repository.ChangeBranch(repositoryName, repository.MASTER_BRANCH)

	//Retrieve the commits data that are going to be push to MASTER_BRANCH
	commits := repository.CommitDiff(repositoryName, repository.MASTER_BRANCH, repository.DEVELOP_BRANCH)
	var releaseCommitsSections []release.ReleaseCommitSection
	for _, value := range commits {
			releaseCommitsSections = append(releaseCommitsSections, release.ReleaseCommitSection{value.Hash,
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

	email.GenerateReleaseEmail(project.FAVOR_PROJECT, "1.3.2", commits)

	releasePage := pages.ReleasePerformedPage{}
	releasePageContent := releasePage.GetContent();
	fmt.Fprintf(w, releasePageContent)
}
