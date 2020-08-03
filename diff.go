package main

import (
	"fmt"
	"reflect"

	"github.com/catilac/plistwatch/go-plist"
)

func Diff(d1 map[string]interface{}, d2 map[string]interface{}) error {
	// check for additions and changes of domains
	for domain, v2 := range d2 {
		if v1, ok := d1[domain]; ok {
			// compare v1 and v2
			if keyptr, eq := cmp(v1, v2); !eq {
				if keyptr != nil {

					// if there is a key, we know this is a map[string]interface{}
					v := v2.(map[string]interface{})
					s, err := marshal(v[*keyptr])
					if err != nil {
						return err
					}

					fmt.Printf("defaults write \"%s\" \"%s\" '%v'\n", domain, *keyptr, *s)
				} else {
					s, err := marshal(v2)
					if err != nil {
						return err
					}
					fmt.Printf("defaults write \"%s\" '%v'\n", domain, *s)
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

func cmp(a interface{}, b interface{}) (*string, bool) {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return nil, false
	}

	switch valA := a.(type) {
	case string:
		return nil, a.(string) == b.(string)
	case int:
		return nil, a.(int) == b.(int)
	case []interface{}:
		valB := b.([]interface{})

		if len(valA) != len(valB) {
			return nil, false
		}
		for i := range valA {
			if key, eq := cmp(valA[i], valB[i]); !eq {
				return key, false
			}
		}
	case map[string]interface{}:
		valB := b.(map[string]interface{})
		if len(valA) != len(valB) {
			return nil, false
		}

		for k := range valA {
			if _, eq := cmp(valA[k], valB[k]); !eq {
				return &k, false
			}
		}
	default:
		fmt.Println("DEFAULT ERR", valA)
	}

	return nil, true
}

func marshal(v interface{}) (*string, error) {
	bytes, err := plist.Marshal(v, plist.OpenStepFormat)
	if err != nil {
		return nil, err
	}

	s := string(bytes)

	return &s, nil

}
