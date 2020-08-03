package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/catilac/plistwatch/go-plist"
)

func getDefaults() (bytes.Buffer, error) {
	var out bytes.Buffer
	cmd := exec.Command("defaults", "read")
	cmd.Env = os.Environ()
	cmd.Stdout = &out
	err := cmd.Run()
	return out, err
}

func main() {
	var prev map[string]interface{}
	var curr map[string]interface{}

	data, err := ioutil.ReadFile("data1.plist")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = plist.Unmarshal(data, &prev)

	data, err = ioutil.ReadFile("data2.plist")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = plist.Unmarshal(data, &curr)

	err = Diff(prev, curr)
	if err != nil {
		fmt.Errorf("Error: ", err)
		os.Exit(-1)
	}
}
