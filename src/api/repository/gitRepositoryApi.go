package repository

import "fmt"
import "os"
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
	MERGE_TOKEN string = "Merge:"
	COMMIT_TOKEN string = "commit"
	JIRA_ISSUE_ID_REGEX string = "\\[CS-\\d+\\]"
	HASH_REGEX string = "[a-z0-9]{40}"
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

func Merge(repository Repository, currentBranch Branch, mergeBranch Branch, jiraTicket string) {
	repoPath := GetLocalRepositoryPath(repository)
	ChangeBranch(repository, currentBranch)

	mergeBranchString := "origin/" + string(mergeBranch)
	currentBranchString := string(currentBranch)

	commitText := fmt.Sprintf("[%v] Merge %v into %v", jiraTicket, mergeBranchString, currentBranchString)

	command.ExecuteCommandWithParams("git", "-C", repoPath, "merge", "--no-ff", "-m",  commitText, mergeBranchString)
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
	commitMessage := fmt.Sprintf("[%v] %v", jiraTicket, message)
	command.ExecuteCommandWithParams("git", "-C", repoPath, "commit", "-m", commitMessage)
}

func Push(repository Repository, branch Branch) {
	repoPath := GetLocalRepositoryPath(repository)
	pushCommand := fmt.Sprintf("git -C %v push origin %v", repoPath, branch);
	command.ExecuteCommand(pushCommand)
}

func Clone(repository Repository) {
	localRepositoryPathFmt := GetLocalRepositoryPath(repository)
	if _, err := os.Stat(localRepositoryPathFmt); os.IsNotExist(err) {
		extRepositoryPathFmt := getExternalRepositoryPath(repository)
		cloneCommand := fmt.Sprintf("git clone %v %v", extRepositoryPathFmt, localRepositoryPathFmt);
		command.ExecuteCommand(cloneCommand)

		Checkout(repository, DEVELOP_BRANCH)
		Checkout(repository, MASTER_BRANCH)
	}
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
	hashRegex := regexp.MustCompile(HASH_REGEX)

	var commits []CommitData

	for _, value := range strings.Split(rawCommitText, "commit") {
		if (len(value) > 0) {

			var commitHash, author, date, jiraTicket, text string

			for _, value := range strings.Split(value, "\n") {
				value = strings.TrimSpace(value)
				if value != "" {
					if strings.Contains(value, AUTHOR_TOKEN) {
						author = strings.TrimSpace(strings.Split(value, AUTHOR_TOKEN)[1]);
					} else if strings.Contains(value, DATE_TOKEN)  {
						date = strings.TrimSpace(strings.Split(value, DATE_TOKEN)[1]);
					} else if jiraTicketRegex.FindString(value) != "" {
						jiraTicket = jiraTicketRegex.FindAllString(value, 1)[0]
						jiraTicket = jiraTicket[1:len(jiraTicket)-1] //Remove []
						text = strings.TrimSpace(jiraTicketRegex.Split(value, 2)[1]);
					} else if hashRegex.FindString(value) != "" {
						commitHash = strings.TrimSpace(value);
					} else if strings.Contains(value, MERGE_TOKEN)  {
						// Do nothing
					} else {
						text = value
					}
				}
			}
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
