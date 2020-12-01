package main

import (
	"log"
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

func TestGetClassPurport(t *testing.T) {
	n, err := getPurport("../books/html/sb/1/1/1.html")
	if err != nil {
		log.Fatal(err)
	}
	pretty.Println(n)
}
