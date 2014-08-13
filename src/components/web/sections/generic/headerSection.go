package generic

type HeaderSection struct {
	Title string
}

func (page HeaderSection)GetTemplateName() string {
	return "generic/headerTemplate.tmpl"
}
