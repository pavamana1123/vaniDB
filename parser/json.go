package main

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getVerse(path string) ([]Verse, error) {
	verse := []Verse{Verse{}}

	raw, err := getTextFromHTML(path, "TEXT", "SYNONYMS")
	if err != nil {
		return verse, err
	}

	for _, l := range strings.Split(raw, "\n") {
		switch {
		case strings.TrimSpace(l) == "":
			continue
		case l == "NEXTVERSE":
			verse[len(verse)-1].Roman = strings.Trim(verse[len(verse)-1].Roman, "\n")
			verse = append(verse, Verse{})
		default:
			verse[len(verse)-1].Roman += l + "\n"
		}
	}

	return verse, nil
}

func getSynonyms(path string) (string, error) {
	return getTextFromHTML(path, "SYNONYMS", "TRANSLATION")
}

func getTranslation(path string) (string, error) {
	return getTextFromHTML(path, "TRANSLATION", "PURPORT")
}

func getPurport(path string) (string, error) {
	return getTextFromHTML(path, "PURPORT", "Link to this page:")
}

func getTextFromHTML(path, t1, t2 string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	delim := "NEXTVERSE"

	r := bufio.NewReader(f)
	start := false
	res := ""
	err = nil
	b := []byte{}
	for err != io.EOF {
		b, _, err = r.ReadLine()
		l := string(b)
		if strings.Contains(l, t2) {
			break
		}
		if start {
			if strings.Contains(l, "margin-top") {
				l = delim + "\n" + l
			}
			if strings.Contains(l, ".jpg") {
				continue
			}
			res += l + "\n"
		}
		if strings.Contains(l, t1) {
			start = true
		}
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		return "", err
	}

	return strings.Trim(doc.Selection.Text(), "\n"), nil

}

func getClasses(path, t1, t2 string) (map[string]string, error) {
	m := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		return m, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	start := false
	reg := regexp.MustCompile(`class="(.*?)"`)
	err = nil
	b := []byte{}
	for err != io.EOF {
		b, _, err = r.ReadLine()
		l := string(b)
		if strings.Contains(l, t2) {
			break
		}
		if start {
			if reg.MatchString(l) {
				mat := reg.FindAllString(l, -1)
				for _, mm := range mat {
					m[mm] = path
				}
			}
		}
		if strings.Contains(l, t1) {
			start = true
		}
	}

	// Load the HTML document

	return m, nil

}
