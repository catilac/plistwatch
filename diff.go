package main

import (
	"fmt"
	"reflect"
	"strings"
	"os/exec"

	"github.com/catilac/plistwatch/go-plist"
)

func Diff(d1 map[string]interface{}, d2 map[string]interface{}) error {
	// check for additions and changes of domains
	for domain, v2 := range d2 {
		if v1, ok := d1[domain]; ok {
			// compare v1 and v2
			prev := v1.(map[string]interface{})
			curr := v2.(map[string]interface{})

			// check for deleted keys
			for key, _ := range prev {
				if _, ok := curr[key]; !ok {
					fmt.Printf("defaults delete \"%s\" \"%s\"\n", domain, key)
				}
			}

			for key, currVal := range curr {
				prevVal, ok := prev[key]
				if !ok || !cmp(prevVal, currVal) {
					// add this key
					s, err := marshal(currVal)
					if err != nil {
						return err
					}
					
					out, _ := exec.Command("defaults", "read-type", domain, key).Output()
					typ := strings.TrimSpace(strings.Replace(string(out), "Type is ", "", -1))
					
					value := ""
					switch typ {
					case "boolean":
						if *s == "1" {
							value = "-bool true"
						} else {
							value = "-bool false"
						}
						break
					case "integer":
					case "float":
					case "date":
						value = "-" + typ + " " + *s
						break
					// strings, arrays and dicts
					default:
						value = "'" + *s + "'"
						break
					}
					fmt.Printf("defaults write \"%s\" \"%s\" %s\n", key, domain, value)
				}
			}
		} else {
			s, err := marshal(v2)
			if err != nil {
				return err
			}
			fmt.Printf("defaults write \"%s\" '%v'\n", domain, *s)
		}
	}

	// check for deletions
	for domain, _ := range d1 {
		if _, ok := d2[domain]; !ok {
			fmt.Printf("defaults delete \"%s\"\n", domain)
		}
	}

	return nil
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
	}

	return true
}

func marshal(v interface{}) (*string, error) {
	bytes, err := plist.Marshal(v, plist.OpenStepFormat)
	if err != nil {
		return nil, err
	}

	s := string(bytes)

	return &s, nil
}
