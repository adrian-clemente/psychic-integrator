package pages

import "components/web/sections/release"
import "components/web/sections/generic"
import "components/printer"

type ReleasePerformedPage struct {
	ReleaseCommitsSections []release.ReleaseCommitSection
}

func (page *ReleasePerformedPage)GetContent() string {

	headerSection := generic.HeaderSection{"Release manager"}
	bodySection := release.BodyReleasePerformedSection{page.ReleaseCommitsSections}
	footerSection := generic.FooterSection{}
	mainSection := generic.MainSection{headerSection, bodySection, footerSection}
	printerPage := printer.PrinterPage{}

	content, _ := printerPage.PrintContent(mainSection);

	return content;
}
