package jenkins

import "log"
import "fmt"
import "net/http"
import "api/config"
import "crypto/tls"

func BuildProject() {

	jenkinsUrl := config.GetProperty("jenkins.url")
	jenkinsUser := config.GetProperty("jenkins.auth.username")
	jenkinsToken := config.GetProperty("jenkins.auth.token")

	log.Printf("Executing build on Jenkins")

	jenkinsUrlFmt := fmt.Sprintf(jenkinsUrl, jenkinsUser, jenkinsToken, "favor.pe.deploy-to-qa")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(jenkinsUrlFmt)
}
