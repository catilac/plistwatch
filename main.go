package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"howett.net/plist"
)

func main() {
	var out bytes.Buffer
	cmd := exec.Command("defaults", "read")
	cmd.Env = os.Environ()
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		fmt.Println("ERROR", err)
		return
	}

	// prints out the output of `defaults read`
	fmt.Println(out.String())

	var dict map[string]interface{}
	format, err := plist.Unmarshal(out.Bytes(), dict)
	fmt.Println("DEBUG: ", format, err, dict)
}
