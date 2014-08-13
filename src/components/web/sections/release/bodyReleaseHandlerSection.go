package release

type BodyReleaseHandlerSection struct {
	Commits []ReleaseCommitSection
	Projects []ReleaseProjectSection
}

func (page BodyReleaseHandlerSection)GetTemplateName() string {
	return "release/bodyReleaseTemplate.tmpl"
}
