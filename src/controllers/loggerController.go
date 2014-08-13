package controllers

import "net/http"
import "components/web/pages"
import "api/logging"
import "fmt"

func ViewLoggerHandler(w http.ResponseWriter, r *http.Request) {
	loggerLines := logging.RetrieveLoggerLines();
	mainPage := pages.LoggerPage{loggerLines}
	mainPageContent := mainPage.GetContent();
	fmt.Fprintf(w, mainPageContent)
}
