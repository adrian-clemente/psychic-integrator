package controllers

import "net/http"
import "components/web/page"
import "api/logging"
import "api/jenkins"
import "fmt"

func ViewDeployHandler(w http.ResponseWriter, r *http.Request) {
	loggerLines := logging.RetrieveLoggerLines();
	mainPage := page.LoggerPage{loggerLines}
	mainPageContent := mainPage.GetContent();

	jenkins.BuildProject()

	fmt.Fprintf(w, mainPageContent)
}
