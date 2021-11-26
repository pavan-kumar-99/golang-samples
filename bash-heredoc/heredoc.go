package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/MakeNowJust/heredoc"
)

func main() {
	commands := `
cat << EOF | bash
echo $PWD
pwd
EOF`
  
	cmd := exec.Command("/bin/bash", "-c", commands)
	fmt.Println(heredoc.Doc(commands))
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", stdoutStderr)
}
