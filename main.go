package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/catilac/plistwatch/go-plist"
)

type dict map[string]interface{}

func getDefaults() (bytes.Buffer, error) {
	var out bytes.Buffer
	cmd := exec.Command("defaults", "read")
	cmd.Env = os.Environ()
	cmd.Stdout = &out
	err := cmd.Run()
	return out, err
}

func main() {
	data, err := ioutil.ReadFile("data1.plist")
	if err != nil {
		fmt.Println(err)
		return
	}

	var dict1 dict
	var dict2 dict

	format, err := plist.Unmarshal(data, &dict1)
	fmt.Println("DEBUG: ", format, err, dict1)

	data, err = ioutil.ReadFile("data2.plist")
	if err != nil {
		fmt.Println(err)
		return
	}

	format, err = plist.Unmarshal(data, &dict2)
	fmt.Println("DEBUG: ", format, err, dict2)

	indent, err := plist.MarshalIndent(dict1, 3, "  ")
	fmt.Println("OK:\n", string(indent))

}
