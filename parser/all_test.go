package main

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"

	"github.com/kr/pretty"
)

func TestNewBook(t *testing.T) {
	b := newBook("https://prabhupadabooks.com/sb/1/1/1?d=1")
	log.Println(b.name, b.levels)
}

func TestGetAllClass(t *testing.T) {
	c1 = []string{"Textnum"}
	c2 = []string{"Synonyms", "Synonyms-SA", "Synonyms-Section", "Titles", "Verse-Section"}
	getAllClass()
}

func TestGetClassInRange1(t *testing.T) {
	log.Println(getClassInRange("../books/bg/1/1.html", []string{}, []string{}))
}

func TestGetClassInRange(t *testing.T) {
	log.Println(getClassInRange("../books/bg/1/1.html", []string{"Textnum"}, []string{"Synonyms", "Synonyms-SA", "Synonyms-Section"}))
}

func TestParts(t *testing.T) {
	n, err := getPurport("../books/html/cc/antya/6/6.html")
	if err != nil {
		log.Fatal(err)
	}
	pretty.Println(n)
}

func TestGetTextID(t *testing.T) {
	fmt.Println(getTextID("../books/html/sb/1/1/1.html"))
}

func TestParse(t *testing.T) {
	filepath.Walk("../books/html/cc/antya/6/6.html", parse)
}
