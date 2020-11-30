package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var (
	allClass map[string]string
	c1, c2   []string
)

func init() {
	allClass = map[string]string{}
}

func getAllClass() {
	err := filepath.Walk("./books", getClass)
	if err != nil {
		log.Println("get all class error: ", err)
	}
	b, err := json.MarshalIndent(allClass, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}

func getClass(path string, info os.FileInfo, err error) error {
	if info.IsDir() && info.Name() == filepath.Base(path) {
		return nil
	}

	if err != nil {
		log.Println(err, path)
		return nil
	}
	if info.IsDir() {
		err = filepath.Walk(path, getClass)
		if err != nil {
			log.Println("get all class error: ", err)
		}
	} else {
		for k, v := range getClassInRange(path, c1, c2) {
			allClass[k] = v
		}
	}
	return nil
}

func getClassInRange(path string, class1, class2 []string) map[string]string {
	classes := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		return classes
	}

	defer f.Close()

	r := bufio.NewReader(f)
	b := []byte{}
	reg := regexp.MustCompile(`class="(.+?)"`)

	err = nil
	cls := ""
	start := false
	if len(class1) == 0 || len(class2) == 0 {
		start = true
	}

	for err != io.EOF {
		b, _, err = r.ReadLine()
		l := string(b)

		if reg.MatchString(l) {
			cls = reg.FindStringSubmatch(l)[1]
			for _, v := range class1 {
				if v == cls {
					start = true
					break
				}
			}
			for _, v := range class2 {
				if v == cls {
					return classes
				}
			}
			if start {
				classes[cls] = path
			}
		}
	}
	return classes
}
