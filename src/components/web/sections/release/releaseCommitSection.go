package release

type ReleaseCommitSection struct {
	Hash string
	Author string
	Date string
	Text string
	JiraTicket string
}

func (page ReleaseCommitSection)GetTemplateName() string {
	return "release/releaseCommitTemplate.tmpl"
}
