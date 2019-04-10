package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type container struct {
	name   string
	status int
}

var unit = map[string]int{
	"running":    1,
	"created":    0,
	"restarting": 0,
	"removing":   0,
	"paused":     0,
	"exited":     0,
	"dead":       0,
}

func Containers() ([]container, error) {
	cmd := exec.Command("docker", "ps", "-aq")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(stdout)
	containers := []container{}
	for scanner.Scan() {
		instance := scanner.Text()
		name, nameErr := name(instance)
		status, statusErr := status(instance)
		if nameErr != nil || statusErr != nil {
			log.Panic(fmt.Sprintf("Failed to get name or status of container '%s'", instance))
			continue
		}
		containers = append(containers, container{name, unit[status]})
	}
	return containers, nil
}

func name(instance string) (string, error) {
	stdout, err := exec.Command("docker", "inspect", "--format='{{ index .Config.Labels \"com.docker.compose.service\"}}'", instance).Output()
	return strings.Trim(string(stdout), "\n'"), err
}

func status(instance string) (string, error) {
	stdout, err := exec.Command("docker", "inspect", "--format='{{.State.Status}}'", instance).Output()
	return strings.Trim(string(stdout), "\n'"), err
}
