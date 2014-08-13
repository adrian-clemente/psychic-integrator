package email

type ReleaseEmail struct {
	JiraIssues []JiraIssueEmail
}

func (page ReleaseEmail)GetTemplateName() string {
	return "email/releaseEmail.html"
}
