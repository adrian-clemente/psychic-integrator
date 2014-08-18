package generic

type FooterSection struct {}

func (page FooterSection)GetTemplateName() string {
	return "web/generic/footerTemplate.tmpl"
}
