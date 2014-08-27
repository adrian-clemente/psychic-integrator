package release

import "components/web/element"

type CommitsSection struct {
	ReleaseCommitsElements []element.CommitElement
	Project string
}

func (page CommitsSection)GetTemplateName() string {
	return "web/release/main/commitsReleaseTemplate.tmpl"
}
