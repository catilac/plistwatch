package main

import (
	"bytes"
	"flag"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func filterDomains(m map[string]any, include, exclude []string) {
	maps.DeleteFunc(m, func(k string, v any) bool {
		// Allow every domain by default
		if len(include) == 0 {
			return false
		}
		for _, pattern := range include {
			if matched, _ := filepath.Match(pattern, strings.ToLower(k)); matched {
				return false
			}
		}
		return true
	})
	maps.DeleteFunc(m, func(k string, v any) bool {
		for _, pattern := range exclude {
			if matched, _ := filepath.Match(pattern, strings.ToLower(k)); matched {
				return true
			}
		}
		return false
	})
}

func main() {
	var include []string
	var exclude []string

	flag.Func("filter", "a comma-separated list of `domains`. Prefix names with \"!\" to exclude them. Supports globbing.", func(s string) error {
		for _, v := range strings.Split(s, ",") {
			v = strings.ToLower(strings.TrimSpace(v))
			domain, found := strings.CutPrefix(v, "!")
			// Users might write "! com.apple.dock" so we trim again
			domain = strings.TrimSpace(domain)
			if domain == "" {
				continue
			}
			if found {
				exclude = append(exclude, domain)
			} else {
				include = append(include, domain)
			}
		}
		return nil
	})
	flag.Parse()

	var prev map[string]interface{}
	var curr map[string]interface{}

	for {
		data, err := getDefaults()
		if _, err = plist.Unmarshal(data.Bytes(), &curr); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		filterDomains(curr, include, exclude)

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
