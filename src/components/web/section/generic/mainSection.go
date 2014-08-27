package generic

import "components"

type MainSection struct {
	HeaderSection HeaderSection
	BodySection components.Component
	FooterSection FooterSection
}

func (page MainSection)GetTemplateName() string {
	return "web/generic/structureTemplate.tmpl"
}
