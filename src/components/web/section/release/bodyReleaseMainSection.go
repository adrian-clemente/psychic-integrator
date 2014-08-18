package release

import "components/web/element"

type BodyReleaseMainSection struct {
	Commits []element.CommitElement
	Projects []element.OptionElement
	ReleaseTypes []element.OptionElement
}

func (page BodyReleaseMainSection)GetTemplateName() string {
	return "web/release/main/bodyReleaseMainTemplate.tmpl"
}
