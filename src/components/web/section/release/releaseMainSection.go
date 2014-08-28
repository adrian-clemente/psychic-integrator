package release

import "components/web/element"

type BodyReleaseMainSection struct {
	Projects []element.OptionElement
	CommitsSection CommitsSection
}

func (page BodyReleaseMainSection)GetTemplateName() string {
	return "web/release/bodyReleaseMainTemplate.tmpl"
}
