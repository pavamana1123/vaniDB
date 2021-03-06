package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Parser struct {
	HTMLFile string
	JSONFile string
	Hf       *os.File
	Jf       *os.File
}

func NewParser(HTMLFile, JSONFile string) (*Parser, error) {
	p := &Parser{
		HTMLFile: HTMLFile,
		JSONFile: JSONFile,
	}

	os.MkdirAll(filepath.Dir(HTMLFile), 0666)
	os.MkdirAll(filepath.Dir(JSONFile), 0666)

	Hf, err := os.Open(HTMLFile)
	if err != nil {
		return p, err
	}

	Jf, err := os.Create(JSONFile)
	if err != nil {
		return p, err
	}

	p.Hf = Hf
	p.Jf = Jf

	return p, nil
}

func (p *Parser) WriteJSON(t Text) error {
	b, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(p.JSONFile, b, 0666)
}

func (p *Parser) Close() {
	p.Hf.Close()
	p.Jf.Close()
}

func writeAllJSON() {
	// filepath.Walk("../books/html/bg", parse)
	filepath.Walk("../books/html/sb", parse)
	// filepath.Walk("../books/html/cc", parse)
}

func parse(path string, info os.FileInfo, err error) error {

	path, _ = filepath.Abs(path)

	if info.IsDir() {
		if filepath.Join(filepath.Dir(path), info.Name()) == path {
			return nil
		}
		return filepath.Walk(path, parse)
	}

	p, err := NewParser(path, strings.ReplaceAll(path, "html", "json"))
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	t, err := getCompleteText(path)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(t.Info.ID)
	return p.WriteJSON(t)

}

func getCompleteText(path string) (Text, error) {
	t := NewText()
	t.Info.ID = getTextID(path)

	v, err := getVerse(path)
	if err != nil {
		return t, err
	}
	t.Verses = v

	s, err := getSynonyms(path)
	if err != nil {
		return t, err
	}
	t.Synonyms = s

	tr, err := getTranslation(path)
	if err != nil {
		return t, err
	}
	t.Translation = tr

	pr, err := getPurport(path)
	if err != nil {
		return t, err
	}
	t.Purport = pr

	return t, nil
}

func getTextID(path string) string {
	path, _ = filepath.Abs(path)
	t := path
	p := []string{}

	for filepath.Base(t) != "books" {
		p = append([]string{filepath.Base(t)}, p...)
		t = filepath.Dir(t)
	}
	return strings.ReplaceAll(strings.Join(p[1:], "\\"), ".html", "")
}
