package os

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func CountOpenFiles() int {
	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("lsof -p %v", os.Getpid())).Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")
	return len(lines) - 1
}
