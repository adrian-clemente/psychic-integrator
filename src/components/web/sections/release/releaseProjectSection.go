package release

type ReleaseProjectSection struct {
	ProjectValue string
	ProjectName string
}

func (page ReleaseProjectSection)GetTemplateName() string {
	return "release/projectsReleaseTemplate.tmpl"
}
