package pages

import "components/web/sections/release"
import "components/web/sections/generic"
import "components/printer"

type ReleaseHandlerPage struct {
	ReleaseCommitsSections []release.ReleaseCommitSection
	ReleaseProjects []release.ReleaseProjectSection
}

func (page *ReleaseHandlerPage)GetContent() string {

	headerSection := generic.HeaderSection{"Release manager"}
	bodySection := release.BodyReleaseHandlerSection{page.ReleaseCommitsSections, page.ReleaseProjects}
	footerSection := generic.FooterSection{}
	mainSection := generic.MainSection{headerSection, bodySection, footerSection}
	printerPage := printer.PrinterPage{}

	content, _ := printerPage.PrintContent(mainSection);

	return content;
}
