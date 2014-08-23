package page

import "components/web/element"
import "components/web/section/generic"
import "components/web/section/release"
import "components/printer"

type ReleaseHandlerPage struct {
	ReleaseCommitsSections []element.CommitElement
	ReleaseProjects []element.OptionElement
	Project string
}

func (page *ReleaseHandlerPage)GetContent() string {

	headerSection := generic.HeaderSection{"Release manager"}
	bodySection := release.BodyReleaseMainSection{page.ReleaseCommitsSections, page.ReleaseProjects, page.Project }
	footerSection := generic.FooterSection{}
	mainSection := generic.MainSection{headerSection, bodySection, footerSection}
	printerPage := printer.PrinterPage{}

	content, _ := printerPage.PrintContent(mainSection);

	return content;
}
