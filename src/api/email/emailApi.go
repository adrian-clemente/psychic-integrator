package email

import (
	"log"
	"net/smtp"
	"strconv"
    "components/printer"
	"components/email"
	"api/config"
	"api/repository"
	"api/jira"
	"fmt"
)

type EmailUser struct {
	Username    string
	Password    string
	EmailServer string
	Port        int
}

type SmtpTemplateData struct {
	Subject string
	Body    string
}

func GenerateReleaseEmail(project string, version string, commits []repository.CommitData) {

	emailReceiver := config.GetProperty("email.auth.user")
	log.Printf("Sending release email to: %v", emailReceiver)

	jiraIssuesEmailMap := make(map[string]bool)
	var jiraIssuesEmail []email.JiraIssueEmail
	for _, commit := range commits {
		if _, exists := jiraIssuesEmailMap[commit.JiraTicket]; !exists {
			jiraIssueFields := jira.RetrieveIssue(commit.JiraTicket)
			jiraIssuesEmailMap[commit.JiraTicket] = true
			jiraIssuesEmail = append(jiraIssuesEmail, email.JiraIssueEmail{ commit.JiraTicket,
					jira.GetJiraIssueBrowseUrl(commit.JiraTicket), jiraIssueFields.Summary })
		}
	}

	printerPage := printer.PrinterPage{}
	content, _ := printerPage.PrintContent(email.ReleaseEmail{jiraIssuesEmail, project, version});

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n";
	subject := fmt.Sprintf("Subject: Release of %v %v\n", project, version)
	msg := []byte(subject + mime + content)

	sendEmail(msg, []string{emailReceiver})

	log.Printf("Email has been sent to: %v", emailReceiver)
}

func sendEmail(message []byte, emailReceivers []string) {
	emailAuthUser := config.GetProperty("email.auth.user")
	emailAuthPassword := config.GetProperty("email.auth.password")
	emailUser := &EmailUser{emailAuthUser, emailAuthPassword, "smtp.gmail.com", 587}

	auth := smtp.PlainAuth("", emailUser.Username, emailUser.Password, emailUser.EmailServer)
	addr := emailUser.EmailServer + ":"+ strconv.Itoa(emailUser.Port)
	err := smtp.SendMail(addr, auth, emailUser.Username, emailReceivers, message)
	if err != nil {
		log.Println(err)
	}
}
