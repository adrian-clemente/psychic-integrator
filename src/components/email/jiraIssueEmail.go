package email

type JiraIssueEmail struct {
	JiraIssueId string
	JiraIssueUrl string
	JiraIssueDesc string
}

func (page JiraIssueEmail)GetTemplateName() string {
	return "email/releaseJiraIssue.html"
}

