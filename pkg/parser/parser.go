package parser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Project struct {
	Path    string
	Items   []Item
	Markers []Marker
}

type Marker struct {
	ID       int
	Name     string
	Start    float64
	End      float64
	IsRegion bool
}

type Item struct {
	Position float64
	Length   float64
	Name     string
	Soure    struct {
		Type string
		File string
	}
	Lyrics Lyrics
}

// Checks wether position is more or equal to start time of the item and also
// that is also smaller then the start + duration
func (i *Item) Contains(pos float64) bool {
	if pos < i.Position {
		return false
	}

	if pos > i.Position+i.Length {
		return false
	}

	return true
}

// Returns a Section by name return err if section not found
func (i *Item) GetSection(name string) (Section, error) {

	if name == "" {
		return Section{}, errors.New("section name cannot be empty")
	}

	for _, s := range i.Lyrics.Sections {
		if strings.ToLower(s.Name) == strings.ToLower(name) {
			return s, nil
		}
	}

	return Section{}, errors.New(fmt.Sprintf("section with the name '%s' not found", name))
}

type Lyrics struct {
	Name     string
	Sections []Section
}

type Section struct {
	Name string
	Text []string
}

func NewLyrics() *Lyrics {
	return &Lyrics{
		Name:     "",
		Sections: []Section{},
	}
}

func NewProject() *Project {

	return &Project{}

}

// NOTE: Dirty function to quickly parser RPP file
// should be reaplaced with proper parser
// at this point all I cancserned is itmes in the project
func (p *Project) ParseReaperFile() error {
	var items []Item
	var item *Item

	f, err := os.Open(p.Path)

	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	for s.Scan() {

		//WARNING: Reaper on windows and on linux save region info differently
		// current parser works for files saved under linux

		tokens := strings.Fields(s.Text())

		if tokens[0] == "MARKER" {
			m, err := parseMarker(s.Text())

			if err != nil {
				return err
			}

			//NOTE: In reaper project marker and reagion is the kinda same thing
			// except that marker doesn't have end position, it alwasy 0
			// And to store end position for the region reaper stores separate
			// marker which has last argument == 1. We will not follow same
			// approch and just add the end to previously inserted marker
			// it has potential bug apportunites, but for now we'll roll with it

			if m.End != 0 {
				p.Markers[len(p.Markers)-1].End = m.End
				continue
			}

			p.Markers = append(p.Markers, m)
		}

		if tokens[0] == "<ITEM" {
			item = &Item{}
		}

		if item != nil {
			switch tokens[0] {

			case "NAME":
				item.Name = strings.Join(tokens[1:], " ")
			case "POSITION":
				item.Position, err = strconv.ParseFloat(tokens[1], 0)
				if err != nil {
					return err
				}
			case "LENGTH":
				item.Length, err = strconv.ParseFloat(tokens[1], 0)
				if err != nil {
					return err
				}
			case "<SOURCE":
				item.Soure.Type = tokens[1]
			case "FILE":
				chunks := strings.Split(s.Text(), "\"")
				item.Soure.File = chunks[1]

			case ">":
				items = append(items, *item)
				item = nil
			}
		}

	}

	for i, item := range items {
		path := filepath.Dir(item.Soure.File)
		path = filepath.Join(path, "lyrics.txt")

		l, err := parseLyrics(path)
		if err != nil {
			continue
		}

		items[i].Lyrics = l
	}

	p.Items = items

	return nil
}

// Parses the line in RPP project file that starts with MARKER keyword
// it's a hacky function to get job done before I'll make real parser
func parseMarker(s string) (Marker, error) {

	//NOTE: if name of the region is a single word we can process in a normal
	// way. But if it contains multiple words it will be surrened by boublequotes
	// and most lickly there will be spaces in the name and it will mess
	// everything up. There four we search for double quotes and if the are
	// present we extract the name and re split the string so that it will be
	// same as if the name is single word.

	var m Marker

	tokens := strings.Fields(s)

	isEnd, err := strconv.ParseBool(tokens[len(tokens)-1])

	if err != nil {
		return m, err
	}

	if isEnd {

		pos, err := strconv.ParseFloat(tokens[2], 10)

		if err != nil {
			return m, err
		}
		m.End = pos

		return m, nil
	}

	start := strings.Index(s, "\"")
	end := strings.LastIndex(s, "\"")

	if start == end && end == -1 || tokens[3] == "" {
		// Normal case
		m.Name = tokens[3]
	} else {
		// Name is a string with multiple words
		m.Name = s[start+1 : end]
		s = fmt.Sprintf("%s name %s", s[:start], s[end+1:])
		tokens = strings.Fields(s)
	}

	id, err := strconv.Atoi(tokens[1])

	if err != nil {
		return m, err
	}

	m.ID = id

	position, err := strconv.ParseFloat(tokens[2], 0)
	if err != nil {
		return m, err
	}

	m.Start = position

	isRegion, err := strconv.ParseBool(tokens[4])

	if err != nil {
		return m, err
	}

	m.IsRegion = isRegion

	return m, nil
}

func parseLyrics(p string) (Lyrics, error) {

	result := NewLyrics()

	f, err := os.Open(p)
	if err != nil {
		return *result, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	result.Name = filepath.Dir(p)

	var name string
	var content []string

	for s.Scan() {
		if s.Text() == "" {
			continue
		}

		if strings.Index(s.Text(), "[") != -1 {
			if name != "" {
				result.Sections = append(result.Sections, Section{
					Name: name,
					Text: content,
				})

				content = []string{}
			}

			name = strings.Trim(s.Text(), "[] ")
			continue
		}

		content = append(content, s.Text())

	}

	result.Sections = append(result.Sections, Section{
		Name: name,
		Text: content,
	})

	return *result, nil
}
