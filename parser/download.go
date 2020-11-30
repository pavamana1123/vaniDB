package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type settings struct {
	Links    []string `json:"links"`
	SavePath string   `json:"savePath"`
	i        int
}

type book struct {
	name   string
	levels []string
	next   string
	page   string
}

var (
	set *settings
)

const (
	linkRegex = `https://prabhupadabooks.com/(\w+)((/(\d+|\w+))+(/(\w+|(\d+(\w+|)-\d+(\w+|)))))\?d=1`
)

func init() {
	dat, err := ioutil.ReadFile("./settings.json")
	if err != nil {
		log.Fatal(`Settings read error: `, err)
	}

	set = &settings{}
	err = json.Unmarshal(dat, set)
	if err != nil {
		log.Fatal(`sunmarshal`, err)
	}
}

func (s *settings) update(link string) {
	s.Links[s.i] = link

	dat, err := json.Marshal(set)
	if err != nil {
		log.Fatal(`up sunmarshal`, err)
	}

	err = ioutil.WriteFile("./settings.json", dat, 0666)
	if err != nil {
		log.Fatal(`Settings read error: `, err)
	}

}

func (b *book) download() {
	fmt.Println("Downloading ", b.next)

	resp, err := http.Get(b.next)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusInternalServerError {
		p, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		b.page = string(p)
	} else {
		log.Fatal("Not 200: ", resp.StatusCode)
	}

	b.savePage()
	b.findNext()

	return
}

func (b *book) findNext() {
	matches := regexp.MustCompile(fmt.Sprintf(`Next: <a href="(%s)">`, linkRegex)).
		FindAllStringSubmatch(b.page, -1)

	if len(matches) == 0 {
		b.next = "END"
		return
	}
	b.next = matches[0][1]
	nb := newBook(matches[0][1])
	if b.name != nb.name {
		b.next = "END"
	}
	b.levels = nb.levels
	set.update(b.next)
}

func (b *book) savePage() {

	path := fmt.Sprintf("%s/%s/%s", set.SavePath, b.name, strings.Join(b.levels[:len(b.levels)-1], "/"))

	err := os.MkdirAll(path, 0666)
	if err != nil {
		log.Fatal(`mkdir`, err)
	}

	f, err := os.Create(fmt.Sprintf("%s/%s.html", path, b.levels[len(b.levels)-1]))
	if err != nil {
		log.Fatal(`new file error: `, err)
	}
	defer f.Close()

	_, err = f.WriteString(b.page)
	if err != nil {
		log.Fatal(`write to file error: `, err)
	}
}

func newBook(link string) *book {
	matches := regexp.MustCompile(linkRegex).FindAllStringSubmatch(link, -1)
	if len(matches) == 0 {
		log.Fatal("Incoherent link: ", link, len(matches))
	}
	return &book{
		name:   matches[0][1],
		next:   link,
		levels: strings.Split(matches[0][2], "/")[1:],
	}
}

func downloadAll() {
	for _, link := range set.Links {
		newBook(link).download()
	}
}
