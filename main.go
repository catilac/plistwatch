package main

import (
	"bytes"
	"fmt"
	//"io/ioutil"
	"os"
	"os/exec"
	"time"

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

	for {
		data, err := getDefaults()
		if _, err = plist.Unmarshal(data.Bytes(), &curr); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		if prev != nil {
			if err = Diff(prev, curr); err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		}

		prev = curr
		curr = nil

		time.Sleep(1 * time.Second)
	}
}
