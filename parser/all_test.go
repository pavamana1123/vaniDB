package main

import (
	"log"
	"path/filepath"
	"testing"
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

func TestGetSections(t *testing.T) {
	n, err := getPurport("../books/bg/7/1.html")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(n)
}

func TestGetClassPurport(t *testing.T) {
	filepath.Walk(`../books/html/`, getKeys)
	for k, v := range m {
		log.Println(k, "\t\t", v)
	}
}
