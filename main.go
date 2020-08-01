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

func walk(d1 dict) {
	fmt.Println("WALK")
	for k, v := range d1 {
		fmt.Println(k, v)
		if v, ok := v.(dict); ok {
			walk(v)
		}
	}
}

func main() {
	data, err := ioutil.ReadFile("mega.plist")
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

	// perform dict1 and dict2
	// we need to identify the domain, and the value changed, and then replace the whole plist, or singular value if it's top level

	walk(dict1)

}
