package element

type CommitElement struct {
	Hash string
	Author string
	Date string
	Text string
	JiraTicket string
	JiraUrl string
}

func (page CommitElement)GetTemplateName() string {
	return "web/element/commitTemplate.tmpl"
}
