package main

import "net/http"
import "controllers"
import "api/logging"
import "io/ioutil"
import "api/config"

func main() {

	config.LoadConfig()

	startLoggingProcess()
	http.HandleFunc("/logger/", controllers.ViewLoggerHandler)
	http.HandleFunc("/release/", controllers.ViewReleaseHandler)
	http.HandleFunc("/release/execute/", controllers.PerformReleaseHandler)
	http.HandleFunc("/deploy/", controllers.ViewDeployHandler)
	http.HandleFunc("/release/commits/", controllers.ViewReleaseCommitsHandler)

	http.Handle("/static/", http.FileServer(http.Dir(".")))

	http.ListenAndServe(":8080", nil)
}

func startLoggingProcess() {
	var c chan string = make(chan string)
	files, _ := ioutil.ReadDir("files")
	for _, file := range files {
		go logging.ProcessFile("files", file.Name(), c)
	}
	go logging.PrintLoggerLine(c)
}
