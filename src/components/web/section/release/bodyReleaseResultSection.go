package release

type BodyReleaseResultSection struct {
	ResultHeader string
	Result string
}

func (page BodyReleaseResultSection)GetTemplateName() string {
	return "web/release/result/bodyReleaseResultTemplate.tmpl"
}
