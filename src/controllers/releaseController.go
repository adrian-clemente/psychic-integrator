package controllers

import "net/http"
import "components/web/page"
import "components/web/element"
import "fmt"
import "api/repository"
import "api/email"
import "api/project"
import "api/jira"
import "api/config"
import "components/web/section/release"

func ViewReleaseHandler(w http.ResponseWriter, r *http.Request) {
	projectName := string(project.FAVOR_PROJECT)
	repositoryName := repository.Repository(project.FAVOR_PROJECT)

	//First ensure the repository exists
	repository.Clone(repositoryName)

	releaseProjects := getReleaseProjects()
	releaseCommitsSections := getProjectReleaseToCommit(repositoryName)

	projectVersion := project.PrintVersion(repositoryName)

	releasePage := page.ReleaseHandlerPage{releaseCommitsSections, releaseProjects, projectName, projectVersion}
	releasePageContent := releasePage.GetContent();
	fmt.Fprintf(w, releasePageContent)
}

func ViewReleaseCommitsHandler(w http.ResponseWriter, r *http.Request) {
	projectNameParam := r.URL.Query()["project"][0]
	repositoryName := repository.Repository(projectNameParam)

	//First ensure the repository exists
	repository.Clone(repositoryName)

	releaseCommitsElements := getProjectReleaseToCommit(repositoryName)
	projectVersion := project.PrintVersion(repositoryName)
	commitSection := release.CommitsSection{releaseCommitsElements, string(projectNameParam), projectVersion}
	commitSectionContent := commitSection.GetContent()

	fmt.Fprintf(w, commitSectionContent)
}

func PerformReleaseHandler(w http.ResponseWriter, r *http.Request) {
	projectNameParam := r.FormValue("project")
	repositoryName := repository.Repository(projectNameParam)

	//First ensure the repository exists
	repository.Clone(repositoryName)

	//Create Jira Issue
	session := jira.Login()
	jiraIssueKey := jira.CreateReleaseIssue(session, projectNameParam)

	//Retrieve wich commits are going to be pushed in this release
	commitsRelease := repository.CommitDiff(repositoryName, repository.MASTER_BRANCH, repository.DEVELOP_BRANCH)

	outputMsg := "Release submit has finished correctly"
	hasError := false

	//Merge develop with master branch
	err := repository.Merge(repositoryName, repository.MASTER_BRANCH, repository.DEVELOP_BRANCH, jiraIssueKey)
	if err == nil {

		//Remove SNAPSHOT from version and add all the changes
		project.SetReleaseVersion(repositoryName)
		projectReleaseVersion := project.PrintVersion(repositoryName)
		repository.AddAll(repositoryName)

		//Create a commit with the changed version
		releaseCommitText := fmt.Sprintf("Release version %v", projectReleaseVersion)
		repository.Commit(repositoryName, releaseCommitText, jiraIssueKey)

		//Push merge and changed version
		err := repository.Push(repositoryName, repository.MASTER_BRANCH)
		if err == nil {
			//Merge master branch into develop and increment version
			err := repository.Merge(repositoryName, repository.DEVELOP_BRANCH, repository.MASTER_BRANCH, jiraIssueKey)
			if (err == nil) {
				project.IncrementMinorVersion(repositoryName)
				project.SetDevelopmentVersion(repositoryName)
				repository.AddAll(repositoryName)
				repository.Commit(repositoryName, "Merge master branch into develop branch", jiraIssueKey)
				err = repository.Push(repositoryName, repository.DEVELOP_BRANCH)
				if err != nil {
					outputMsg = "Error while pushing to the repository"
					hasError = true
				}

			} else {
				outputMsg = "Conflicts while merging MASTER branch into DEVELOP branch. You should do manually the merge"
				hasError = true
			}
			//Close ticket
			jira.CloseIssue(session, jiraIssueKey)
			//Send the email with all the commits that were merged
			email.GenerateReleaseEmail(projectNameParam, projectReleaseVersion, commitsRelease)
		} else {
			outputMsg = "Error while pushing to the repository"
			hasError = true
		}
	} else {
		outputMsg = "Conflicts while merging DEVELOP branch into MASTER branch. You should do manually the merge"
		hasError = true
	}

	releaseResultSectionContent := getReleaseContent(projectNameParam, outputMsg, hasError)
	fmt.Fprintf(w, releaseResultSectionContent)
}

func getReleaseContent(project string, outputMsg string, hasError bool) string {
	resultHeader := fmt.Sprintf("Release summary of %v:", project)

	contentColorClass := "correctResult"
	if hasError {
		contentColorClass = "errorResult"
	}

	releaseResultSection := release.BodyReleaseResultSection{resultHeader, outputMsg, contentColorClass}
	releaseResultSectionContent := releaseResultSection.GetContent();
	return releaseResultSectionContent
}

/**
 * Retrieve the release types that can be perform
 */
func getReleaseTypes() []element.OptionElement {
	var releaseTypes []element.OptionElement
	for _, releaseType := range project.GetReleaseTypes() {
		releaseTypes = append(releaseTypes, element.OptionElement{ string(releaseType.ReleaseTypeKey),
				string(releaseType.ReleaseTypeName) })
	}
	return releaseTypes
}

/**
 * Retrieve the release projects that can be executed
 */
func getReleaseProjects() []element.OptionElement {
	var releaseProjects []element.OptionElement
	for _, project := range project.GetProjects() {
		releaseProjects = append(releaseProjects, element.OptionElement{ string(project.ProjectKey),
				string(project.ProjectName) })
	}
	return releaseProjects
}

/**
 * Retrieve all the commits that are candidate to be push in the release
 */
func getProjectReleaseToCommit(repositoryName repository.Repository) []element.CommitElement {
	commits := repository.CommitDiff(repositoryName, repository.MASTER_BRANCH, repository.DEVELOP_BRANCH)
	var releaseCommitsSections []element.CommitElement
	jiraUrl := config.GetProperty("jira.url.domain")

	for _, value := range commits {
		releaseCommitsSections = append(releaseCommitsSections, element.CommitElement{value.Hash,
				value.Author, value.Date, value.Text, value.JiraTicket, jiraUrl})
	}
	return releaseCommitsSections
}
