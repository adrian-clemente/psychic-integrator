package release

type BodyReleasePerformedSection struct {
	Commits []ReleaseCommitSection
}

func (page BodyReleasePerformedSection)GetTemplateName() string {
	return "release/bodyReleaseTemplate.tmpl"
}
