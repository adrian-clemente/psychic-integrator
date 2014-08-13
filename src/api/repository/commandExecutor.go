package repository

import (
	"log"
	"os/exec"
	"sync"
	"strings"
)

func executeCommand(cmd string) string {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	log.Println("command is ", cmd)
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head,parts...).Output()
	if err != nil {
		log.Println(err)
	}
	wg.Done() // Need to signal to waitgroup that this goroutine is done

	return string(out)
}

func executeCommandWithParams(cmd string, params ...string) string {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	out, err := exec.Command(cmd, params...).Output()
	if err != nil {
		log.Println(err)
	}
	wg.Done() // Need to signal to waitgroup that this goroutine is done

	return string(out)
}

