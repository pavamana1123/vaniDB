package main

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func getVerse(path string) ([]Verse, error) {
	verse := []Verse{Verse{}}

	raw, err := getVSTFromHTML(path, "TEXT", "SYNONYMS")
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
	return getVSTFromHTML(path, "SYNONYMS", "TRANSLATION")
}

func getTranslation(path string) (string, error) {
	return getVSTFromHTML(path, "TRANSLATION", "PURPORT;Link to this page")
}

func getVSTFromHTML(path, t1, t2 string) (string, error) {
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
		for _, tt2 := range strings.Split(t2, ";") {
			if strings.Contains(l, tt2) {
				goto found
			}
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

found:
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		return "", err
	}

	s := doc.Selection.Text()
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\n")

	return s, nil

}

func getPurport(path string) ([]Paragraph, error) {

	p := []Paragraph{}

	f, err := os.Open(path)
	if err != nil {
		return p, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	start := false
	res := ""
	err = nil
	b := []byte{}
	for err != io.EOF {
		b, _, err = r.ReadLine()
		l := string(b)
		if strings.Contains(l, "Link to this page:") {
			break
		}
		if start {
			res += l + "\n"
		}
		if strings.Contains(l, "PURPORT") {
			start = true
		}
	}

	res = strings.ReplaceAll(res, "<br>", "\n")

	node, err := html.Parse(strings.NewReader(res))
	if err != nil {
		return p, err
	}

	node = node.FirstChild.FirstChild.NextSibling.FirstChild
	for node != nil {
		attr := node.Attr
		class := ""

		for _, a := range attr {
			if a.Key == "class" {
				class = a.Val
			}
		}

		gnode := goquery.NewDocumentFromNode(node)
		text := strings.TrimSpace(gnode.Selection.Text())
		if text != "" {
			p = append(p, Paragraph{
				Type:    getParaType(class),
				Content: text,
			})
		}

		node = node.NextSibling
	}

	return p, nil

}

func getParaType(class string) string {
	verseClasses := []string{
		"VerseRef",
		"Centered-Verse-in-purp",
		"Verse",
		"Verse-in-purp",
		"Verse-Text",
		"Prose-Verse-in-purp",
		"One-line-verse-in-purp",
	}

	for _, c := range verseClasses {
		if c == class {
			return "verse"
		}
	}

	return "normal"
}
