package main

import (
	"bufio"
	"fmt"
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

const inspectCmd = "docker inspect --format='%s,%s' $(docker ps -aq)"
const labelFormat = `{{ index .Config.Labels "com.docker.compose.service"}}`
const statusFormat = "{{ .State.Status }}"

func DockerInspect() ([]container, error) {
	containers := []container{}
	cmd := exec.Command("bash", "-c", fmt.Sprintf(inspectCmd, labelFormat, statusFormat))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		info := strings.Split(scanner.Text(), ",")
		if len(info[0]) > 0 {
			containers = append(containers, container{info[0], unit[info[1]]})
		}
	}

	return containers, cmd.Wait()
}
