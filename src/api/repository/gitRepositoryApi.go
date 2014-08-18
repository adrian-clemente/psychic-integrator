package repository

import "fmt"
import "strings"
import "regexp"
import "api/config"
import "api/command"

type Branch string
type Repository string

const (
	DEVELOP_BRANCH Branch = "develop"
	MASTER_BRANCH Branch = "master"

	AUTHOR_TOKEN string = "Author:"
	DATE_TOKEN string = "Date:"
	JIRA_ISSUE_ID_REGEX string = "CS-\\d+"
)

type CommitData struct {
	Hash string
	Author string
	Date string
	Text string
	JiraTicket string
}

func Log(repository Repository, numCommits int, branch Branch) []CommitData {
	logCommand := fmt.Sprintf("git -C %v log -n%v %v", repository, numCommits, branch)
	rawCommitText := command.ExecuteCommand(logCommand)
	return parseCommitResponse(rawCommitText)
}

func Merge(repository Repository, currentBranch Branch, mergeBranch Branch) {
	repoPath := GetLocalRepositoryPath(repository)
	Checkout(repository, currentBranch)
	command.ExecuteCommandWithParams("git", "-C", repoPath, "merge", "--no-ff", "-m",  "Merge with " + string(mergeBranch),
		string(mergeBranch))
}

func AddAll(repository Repository) {
	repoPath := GetLocalRepositoryPath(repository)
	addCommand := fmt.Sprintf("git -C %v add --all", repoPath);
	command.ExecuteCommand(addCommand)
}

func Checkout(repository Repository, branch Branch) {
	repoPath := GetLocalRepositoryPath(repository)
	checkoutCommand := fmt.Sprintf("git -C %v checkout -b %v origin/%v", repoPath, branch, branch);
	command.ExecuteCommand(checkoutCommand)
}

func ChangeBranch(repository Repository, branch Branch) {
	repoPath := GetLocalRepositoryPath(repository)
	checkoutCommand := fmt.Sprintf("git -C %v checkout %v", repoPath, branch);
	command.ExecuteCommand(checkoutCommand)
}

func Commit(repository Repository, message string, jiraTicket string) {
	repoPath := GetLocalRepositoryPath(repository)
	commitMessage := jiraTicket + message
	command.ExecuteCommandWithParams("git", "-C", repoPath, "commit", "-m", commitMessage)
}

func Push(repository Repository, branch string) {
	repoPath := GetLocalRepositoryPath(repository)
	pushCommand := fmt.Sprintf("git -C %v push origin %v", repoPath, branch);
	command.ExecuteCommand(pushCommand)
}

func Clone(repository Repository) {
	localRepositoryPathFmt := GetLocalRepositoryPath(repository)
	extRepositoryPathFmt := getExternalRepositoryPath(repository)

	cloneCommand := fmt.Sprintf("git clone %v %v", extRepositoryPathFmt, localRepositoryPathFmt);
	command.ExecuteCommand(cloneCommand)
}

func CommitDiff(repository Repository, firstBranch Branch, secondBranch Branch) []CommitData {
	repoPath := GetLocalRepositoryPath(repository)
	diffCommand := fmt.Sprintf("git -C %v log %v..%v", repoPath, firstBranch, secondBranch);
	rawCommitText := command.ExecuteCommand(diffCommand)
	return parseCommitResponse(rawCommitText)
}

func GetLocalRepositoryPath(repository Repository) string {
	localRepositoryPath := config.GetProperty("repository.local.path")
	localRepositoryPathFmt := fmt.Sprintf(localRepositoryPath, repository)

	return localRepositoryPathFmt
}

func getExternalRepositoryPath(repository Repository) string {
	username := config.GetProperty("repository.external.username")
	password := config.GetProperty("repository.external.password")
	extRepositoryPath := config.GetProperty("repository.external.path")
	extRepositoryPathFmt := fmt.Sprintf(extRepositoryPath, username, password, repository)

	return extRepositoryPathFmt
}

func parseCommitResponse(rawCommitText string) []CommitData {
	jiraTicketRegex := regexp.MustCompile(JIRA_ISSUE_ID_REGEX)
	var commits []CommitData

	for _, value := range strings.Split(rawCommitText, "commit") {
		if (len(value) > 0) {
			commitContent := [5]string{}
			index := 0
			for _, value := range strings.Split(value, "\n") {
				if (len(value) > 0 && index < 5) {
					commitContent[index] = strings.TrimSpace(value)
					index++
				}
			}
			commitHash := commitContent[0]
			author := removeToken(commitContent[1], AUTHOR_TOKEN)
			date := removeToken(commitContent[2], DATE_TOKEN)
			jiraTicket := jiraTicketRegex.FindString(commitContent[3])
			text := removeToken(commitContent[3], "]")

			commits = append(commits, CommitData{commitHash, author, date, text, jiraTicket})
		}
	}

	return commits
}

func removeToken(originalString string, token string) string {
	splitString := strings.Split(originalString, token);
	if (len(splitString) > 1) {
		return strings.TrimSpace(splitString[1])
	} else {
		return originalString
	}
}
