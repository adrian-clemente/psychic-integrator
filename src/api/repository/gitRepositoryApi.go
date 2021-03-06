package repository

import "fmt"
import "os"
import "strings"
import "regexp"
import "api/config"
import "api/command"
import "log"

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
	rawCommitText, _ := command.ExecuteCommand(logCommand)
	return parseCommitResponse(rawCommitText)
}

func Merge(repository Repository, currentBranch Branch, mergeBranch Branch, jiraTicket string) error {
	if jiraTicket == "" {
		jiraTicket = "CS-0"
	}

	repoPath := GetLocalRepositoryPath(repository)
	ChangeBranch(repository, currentBranch)

	mergeBranchString := "origin/" + string(mergeBranch)
	currentBranchString := string(currentBranch)

	commitText := fmt.Sprintf("[%v] Merge %v into %v", jiraTicket, mergeBranchString, currentBranchString)

	_, err := command.ExecuteCommandWithParams("git", "-C", repoPath, "merge", "--no-ff", "-m",  commitText,
		mergeBranchString)

	return err
}

func AddAll(repository Repository) {
	repoPath := GetLocalRepositoryPath(repository)
	addCommand := fmt.Sprintf("git -C %v add --all", repoPath)
	command.ExecuteCommand(addCommand)
}

func Checkout(repository Repository, branch Branch) {
	repoPath := GetLocalRepositoryPath(repository)
	checkoutCommand := fmt.Sprintf("git -C %v checkout -b %v origin/%v", repoPath, branch, branch)
	command.ExecuteCommand(checkoutCommand)
}

func ChangeBranch(repository Repository, branch Branch) {
	repoPath := GetLocalRepositoryPath(repository)
	checkoutCommand := fmt.Sprintf("git -C %v checkout %v", repoPath, branch)
	command.ExecuteCommand(checkoutCommand)
}

func Commit(repository Repository, message string, jiraTicket string) {
	if jiraTicket == "" {
		jiraTicket = "CS-0"
	}

	repoPath := GetLocalRepositoryPath(repository)
	commitMessage := fmt.Sprintf("[%v] %v", jiraTicket, message)
	command.ExecuteCommandWithParams("git", "-C", repoPath, "commit", "-m", commitMessage)
}

func Pull(repository Repository, branch Branch) error {
	repoPath := GetLocalRepositoryPath(repository)
	ChangeBranch(repository, branch)
	pushCommand := fmt.Sprintf("git -C %v pull origin %v", repoPath, branch)
	_, err := command.ExecuteCommand(pushCommand)
	return err
}

func Push(repository Repository, branch Branch) error {
	repoPath := GetLocalRepositoryPath(repository)
	pushCommand := fmt.Sprintf("git -C %v push origin %v", repoPath, branch)
	_, err := command.ExecuteCommand(pushCommand)
	return err
}

func Clone(repository Repository) {
	localRepositoryPathFmt := GetLocalRepositoryPath(repository)
	if _, err := os.Stat(localRepositoryPathFmt); os.IsNotExist(err) {

		//Try clone repository from external path
		extRepositoryPathFmt := getExternalRepositoryPath(repository, "repository.external.path")
		cloneCommand := fmt.Sprintf("git clone %v %v", extRepositoryPathFmt, localRepositoryPathFmt)
		_, err := command.ExecuteCommand(cloneCommand)
		if err != nil {
			//Try clone repository from external github path
			extRepositoryPathFmt = getExternalRepositoryPath(repository, "repository.external.github.path")
			cloneCommand = fmt.Sprintf("git clone %v %v", extRepositoryPathFmt, localRepositoryPathFmt)
			_, err := command.ExecuteCommand(cloneCommand)
			if err == nil {
				checkoutLocalBranches(repository)
			}
		} else {
			checkoutLocalBranches(repository)
		}
	}
}

func checkoutLocalBranches(repository Repository) {
	Checkout(repository, DEVELOP_BRANCH)
	Checkout(repository, MASTER_BRANCH)
}

func CommitDiff(repository Repository, firstBranch Branch, secondBranch Branch) []CommitData {

	Pull(repository, firstBranch)
	Pull(repository, secondBranch)

	repoPath := GetLocalRepositoryPath(repository)
	diffCommand := fmt.Sprintf("git -C %v log %v..%v", repoPath, firstBranch, secondBranch)
	rawCommitText, _ := command.ExecuteCommand(diffCommand)
	return parseCommitResponse(rawCommitText)
}

func GetLocalRepositoryPath(repository Repository) string {
	localRepositoryPath := config.GetProperty("repository.local.path")
	localRepositoryPathFmt := fmt.Sprintf(localRepositoryPath, repository)

	return localRepositoryPathFmt
}

func Delete(repository Repository) {
	localRepositoryPath := GetLocalRepositoryPath(repository)
	err := os.RemoveAll(localRepositoryPath)

	if err != nil {
		log.Println(err)
		return
	} else {
		log.Printf("Repository deleted %v", localRepositoryPath)
	}
}

func getExternalRepositoryPath(repository Repository, path string) string {
	extRepositoryPath := config.GetProperty(path)
	extRepositoryPathFmt := fmt.Sprintf(extRepositoryPath, repository)

	return extRepositoryPathFmt
}

func parseCommitResponse(rawCommitText string) []CommitData {
	jiraTicketRegex := regexp.MustCompile(JIRA_ISSUE_ID_REGEX)
	hashRegex := regexp.MustCompile(HASH_REGEX)

	var commits []CommitData
	var commitHash, author, date, jiraTicket, text string

	for _, value := range strings.Split(rawCommitText, "\n") {
		value = strings.TrimSpace(value)
		if value != "" {
			if strings.Contains(value, AUTHOR_TOKEN) {
				author = strings.TrimSpace(strings.Split(value, AUTHOR_TOKEN)[1])
			} else if strings.Contains(value, DATE_TOKEN)  {
				date = strings.TrimSpace(strings.Split(value, DATE_TOKEN)[1])
			} else if jiraTicketRegex.FindString(value) != "" {
				jiraTicket = jiraTicketRegex.FindAllString(value, 1)[0]
				jiraTicket = jiraTicket[1:len(jiraTicket)-1] //Remove []
				text = strings.TrimSpace(jiraTicketRegex.Split(value, 2)[1])
			} else if hashRegex.FindString(value) != "" {
				if commitHash != "" {
					commits = append(commits, CommitData{commitHash, author, date, text, jiraTicket})
					text = ""
					jiraTicket = ""
					author = ""
					date = ""
				}

				commitHash = hashRegex.FindString(value)
			} else if strings.Contains(value, MERGE_TOKEN)  {
				// Do nothing
			} else {
				text = text + "</br>" + value
			}
		}
	}

	return commits
}

func removeToken(originalString string, token string) string {
	splitString := strings.Split(originalString, token)
	if (len(splitString) > 1) {
		return strings.TrimSpace(splitString[1])
	} else {
		return originalString
	}
}
