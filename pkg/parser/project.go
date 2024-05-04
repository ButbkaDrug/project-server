package parser

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/butbkadrug/reaper-lyrics/pkg/sse"
)

type Project struct {
	Stats    os.FileInfo
	Path     string
	LastLine string
	Items    []Item
	Markers  []Marker
	System   system
}

type Marker struct {
	ID       int
	Name     string
	Start    float64
	End      float64
	IsRegion bool
}

type Item struct {
	ID       int
	Position float64
	Length   float64
	Name     string
	Soure    struct {
		Type string
		File string
	}
	Lyrics Lyrics
	Image  string
}

type Lyrics struct {
	Name     string
	Sections []Section
}

type Section struct {
	ID   int
	Name string
	Text []string
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

func NewLyrics() *Lyrics {
	return &Lyrics{
		Name:     "",
		Sections: []Section{},
	}
}

func NewProject() *Project {

	return &Project{}

}

// finds first item to contain position
func (p *Project) GetItem(pos float64) (Item, error) {

	for _, item := range p.Items {
		if item.Contains(pos) {
			return item, nil
		}
	}

	return Item{}, errors.New("item not found")
}

// finds firest marker to contain position
func (p *Project) GetMarker(pos float64) (Marker, error) {
	for _, marker := range p.Markers {
		if pos >= marker.Start && pos < marker.End {
			return marker, nil
		}
	}

	return Marker{}, errors.New("marker not found")
}

func (p *Project) getStats() (os.FileInfo, error) {
	info, err := os.Stat(p.Path)

	if err != nil {
		return nil, err
	}

	return info, nil
}

// watching project file for changes
// should be ran in separate go roteen as it contains infinite for loop
func (p *Project) Watch(clients *sse.SSE) error {

	for {

		info, err := p.getStats()

		if err != nil {
			return err
		}

		if info.ModTime() != p.Stats.ModTime() && info.Size() != p.Stats.Size() {

			go clients.Broadcast(true)

			log.Println("File was changed!!")
			if err := p.ParseReaperFile(); err != nil {
				log.Println("parsing error: ", err)
			}

			p.Stats = info
		}

		time.Sleep(time.Second)

	}
}

func (p *Project) readProjectFile() ([]string, error) {
	var lines []string

	f, err := os.Open(p.Path)
	info, err := p.getStats()

	if err != nil {
		return lines, err
	}

	p.Stats = info

	defer f.Close()

	s := bufio.NewScanner(f)

	for s.Scan() {
		lines = append(lines, s.Text())
	}

	return lines, nil

}

func (p *Project) Parse(lines []string) error {
	var items []Item
	var markers []Marker
	var item *Item
	var err error

	for _, line := range lines {

		tokens := strings.Fields(line)

		switch t := tokens[0]; t {
		case "MARKER":

			//NOTE: In reaper project marker and reagion is the kinda same thing
			// except that marker doesn't have end position, it alwasy 0
			// And to store end position for the region reaper stores separate
			// marker which has last argument == 1. We will not follow same
			// approch and just add the end to previously inserted marker
			// it has potential bug apportunites, but for now we'll roll with it

			l := NewLexer(line)

			var values []Token

			for tok := l.NextToken(); tok.Type != EOF; tok = l.NextToken() {
				values = append(values, tok)

			}

			m, err := p.parseMarker(values)

			if err != nil {
				return err
			}

			if m.End != 0 {
				markers[len(markers)-1].End = m.End
			} else {
				markers = append(markers, m)
			}

		case "<ITEM":
			item = &Item{}
		case "NAME":
			if item == nil {
				continue
			}
			item.Name = strings.Join(tokens[1:len(tokens)-1], " ")
			if len(item.Name) > 0 {
				item.Name = item.Name[1:]
			}
			item.Image = GetImageName(item.Name)
		case "POSITION":
			if item == nil {
				continue
			}
			item.Position, err = strconv.ParseFloat(tokens[1], 0)
			if err != nil {
				return err
			}
		case "LENGTH":
			if item == nil {
				continue
			}
			item.Length, err = strconv.ParseFloat(tokens[1], 0)
			if err != nil {
				return err
			}
		case "<SOURCE":
			if item == nil {
				continue
			}
			item.Soure.Type = tokens[1]
		case "FILE":
			if item == nil {
				continue
			}
			chunks := strings.Split(line, "\"")
			item.Soure.File = chunks[1]

		case ">":
			if item == nil {
				continue
			}

			ext := filepath.Ext(item.Soure.File)

			if strings.ToLower(ext) != ".rpp" {
				continue
			}

			item.ID = len(items)
			items = append(items, *item)
			item = nil
		default:
			continue
		}

	}

	for i, item := range items {
		ext := filepath.Ext(item.Soure.File)

		if strings.ToLower(ext) != ".rpp" {
			continue
		}

		dir := filepath.Dir(item.Soure.File)
		path := filepath.Join(dir, "lyrics.txt")

		l, err := ParseLyrics(path)
		if err != nil {
			log.Println(item.Soure.File, err)
			continue
		}

		items[i].Lyrics = l

	}

	p.Items = items
	p.Markers = markers

	return nil
}

func (p *Project) getProjectSystem(s string) system {
	if idx := strings.Index(s, "linux"); idx != -1 {
		return LINUX
	}

	if idx := strings.Index(s, "win"); idx != -1 {
		return WIN
	}

	return UNKNOWN
}

func (p *Project) ParseReaperFile() error {
	//WARNING: Reaper on windows and on linux save region info differently
	// current parser works for files saved under linux

	lines, err := p.readProjectFile()
	if err != nil {
		return err
	}

	p.System = p.getProjectSystem(lines[0])

	switch p.System {
	case LINUX, WIN:
		p.Parse(lines)
	default:
		return errors.New("failed to define project file system")
	}

	return nil
}

// Parses the line in R}PP project file that starts with MARKER keyword
// it's a hacky function to get job done before I'll make real parser
func (p Project) parseMarker(tokens []Token) (Marker, error) {

	//NOTE: if name of the region is a single word we can process in a normal
	// way. But if it contains multiple words it will be surrened by boublequotes
	// and most lickly there will be spaces in the name and it will mess
	// everything up. There four we search for double quotes and if the are
	// present we extract the name and re split the string so that it will be
	// same as if the name is single word.

	var m Marker
	var isEnd, isRegion bool
	//NOTE:  System Dependend filed parsing

	id, err := strconv.Atoi(tokens[1].Value)

	if err != nil {
		return m, err
	}

	position, err := strconv.ParseFloat(tokens[2].Value, 0)
	if err != nil {
		return m, err
	}

	switch p.System {
	case LINUX:
		isEnd, err = strconv.ParseBool(tokens[len(tokens)-1].Value)

		if err != nil {
			return m, err
		}

		isRegion, err = strconv.ParseBool(tokens[4].Value)

		if err != nil {
			return m, err
		}

	case WIN:

		if tokens[7].Value == "R" {
			isRegion = true
		}

		if len(tokens) < 9 {
			isEnd = true
		}
	}

	m.Name = strings.Trim(tokens[3].Value, " \t\n\r")
	m.IsRegion = isRegion
	m.ID = id

	if isEnd {
		m.End = position
	} else {
		m.Start = position
	}

	return m, nil

}

func GetImageName(n string) string {
	toks := strings.Split(n, "-")
	name := strings.Trim(toks[0], " \t\r\n")
	name = strings.ReplaceAll(name, " ", "-")
	return strings.ToLower(name) + ".jpg"
}

func ParseLyrics(p string) (Lyrics, error) {

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

		if strings.Index(s.Text(), "[") != -1 {
			if name != "" {
				result.Sections = append(result.Sections, Section{
					ID:   len(result.Sections),
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
