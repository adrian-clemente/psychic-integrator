package repository

import "fmt"
import "strings"
import "regexp"
import "api/config"

type Branch string
type Repository string

const (
	DEVELOP_BRANCH Branch = "develop"
	MASTER_BRANCH Branch = "master"

	AUTHOR_TOKEN = "Author:"
	DATE_TOKEN = "Date:"
	JIRA_ISSUE_ID_REGEX = "CS-\\d+"
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
	rawCommitText := executeCommand(logCommand)
	return parseCommitResponse(rawCommitText)
}

func Merge(repository Repository, currentBranch Branch, mergeBranch Branch) {
	repoPath := getLocalRepositoryPath(repository)
	Checkout(repository, currentBranch)
	executeCommandWithParams("git", "-C", repoPath, "merge", "--no-ff", "-m",  "Merge with " + string(mergeBranch),
		string(mergeBranch))
}

func AddAll(repository Repository) {
	repoPath := getLocalRepositoryPath(repository)
	addCommand := fmt.Sprintf("git -C %v add --all", repoPath);
	executeCommand(addCommand)
}

func Checkout(repository Repository, branch Branch) {
	repoPath := getLocalRepositoryPath(repository)
	checkoutCommand := fmt.Sprintf("git -C %v checkout -b %v origin/%v", repoPath, branch, branch);
	executeCommand(checkoutCommand)
}

func ChangeBranch(repository Repository, branch Branch) {
	repoPath := getLocalRepositoryPath(repository)
	checkoutCommand := fmt.Sprintf("git -C %v checkout %v", repoPath, branch);
	executeCommand(checkoutCommand)
}

func Commit(repository Repository, message string, jiraTicket string) {
	repoPath := getLocalRepositoryPath(repository)
	commitMessage := jiraTicket + message
	executeCommandWithParams("git", "-C", repoPath, "commit", "-m", commitMessage)
}

func Push(repository Repository, branch string) {
	repoPath := getLocalRepositoryPath(repository)
	pushCommand := fmt.Sprintf("git -C %v push origin %v", repoPath, branch);
	executeCommand(pushCommand)
}

func Clone(repository Repository) {
	localRepositoryPathFmt := getLocalRepositoryPath(repository)
	extRepositoryPathFmt := getExternalRepositoryPath(repository)

	cloneCommand := fmt.Sprintf("git clone %v %v", extRepositoryPathFmt, localRepositoryPathFmt);
	executeCommand(cloneCommand)
}

func CommitDiff(repository Repository, firstBranch Branch, secondBranch Branch) []CommitData {
	repoPath := getLocalRepositoryPath(repository)
	diffCommand := fmt.Sprintf("git -C %v log %v..%v", repoPath, firstBranch, secondBranch);
	rawCommitText := executeCommand(diffCommand)
	return parseCommitResponse(rawCommitText)
}

func getLocalRepositoryPath(repository Repository) string {
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
