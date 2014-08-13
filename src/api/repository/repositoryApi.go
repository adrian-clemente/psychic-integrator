package repository

type RepositoryApi interface {

	Log(repository Repository, numCommits int, branch Branch) []CommitData

	Merge(repository Repository, currentBranch Branch, mergeBranch Branch)

	AddAll(repository Repository)

	Checkout(repository Repository, branch Branch)

	ChangeBranch(repository Repository, branch Branch)

	Commit(repository Repository, message string, jiraTicket string)

	Push(repository Repository, branch string)

	Clone(repository Repository)

	CommitDiff(repository Repository, firstBranch Branch, secondBranch Branch) []CommitData
}
