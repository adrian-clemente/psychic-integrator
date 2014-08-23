package release

import "components/web/element"

type BodyReleaseMainSection struct {
	Commits []element.CommitElement
	Projects []element.OptionElement
	Project string
}

func (page BodyReleaseMainSection)GetTemplateName() string {
	return "web/release/main/bodyReleaseMainTemplate.tmpl"
}
