package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"

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

func cmp(a interface{}, b interface{}) bool {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}

	switch valA := a.(type) {
	case string:
		return a.(string) == b.(string)
	case int:
		return a.(int) == b.(int)
	case []interface{}:
		valB := b.([]interface{})

		if len(valA) != len(valB) {
			return false
		}
		for i := range valA {
			if !cmp(valA[i], valB[i]) {
				return false
			}
		}
	case map[string]interface{}:
		valB := b.(map[string]interface{})
		if len(valA) != len(valB) {
			return false
		}

		for k := range valA {
			if !cmp(valA[k], valB[k]) {
				return false
			}
		}
	default:
		fmt.Println("DEFAULT ERR", valA)
	}

	return true
}

func diff(d1 dict, d2 dict) {
	// check for additions and changes of domains
	for domain, v2 := range d2 {
		if v1, ok := d1[domain]; ok {
			// compare v1 and v2
			if !cmp(v1, v2) {
				fmt.Printf("defaults write \"%s\" \"%v\"\n", domain, v2)
			}
		} else {
			fmt.Printf("defaults write \"%s\" \"%v\"\n", domain, v2)
		}
	}

	// check for deletions
	for domain, _ := range d1 {
		if _, ok := d2[domain]; !ok {
			fmt.Printf("defaults delete \"%s\"\n", domain)
		}
	}
}

func main() {
	var dict1 dict
	var dict2 dict

	data, err := ioutil.ReadFile("data1.plist")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = plist.Unmarshal(data, &dict1)

	data, err = ioutil.ReadFile("data2.plist")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = plist.Unmarshal(data, &dict2)

	diff(dict1, dict2)

}
