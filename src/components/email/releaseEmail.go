package email

type ReleaseEmail struct {
	JiraIssues []JiraIssueEmail
	Project string
	Version string
}

func (page ReleaseEmail)GetTemplateName() string {
	return "email/releaseEmail.html"
}
