package projects

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/butbkadrug/reaper-lyrics/pkg/parser"
)

type Library struct {
	Path     string
	Projects []Project
}

type Project struct {
	ID     string
	Name   string
	Path   string
	Lyrics parser.Lyrics
}

func NewLibrary(path string) *Library {
	return &Library{
		Path:     path,
		Projects: []Project{},
	}
}

func generateID(name string) string {
	toks := strings.Split(name, "-")
	ID := strings.Trim(toks[0], " -")
	ID = strings.ReplaceAll(ID, " ", "-")
	return strings.ToLower(ID)
}

// Parses the project dir. Directory should contain at least one project that
// matches direcotry name and that one will be considered the actual project
// if no projects matches directory name it's not a valid project dir.
func ParseProjectDir(projectPath string) (Project, error) {

	var proj Project
    var filename string

	foldeName := filepath.Base(projectPath)
	entries, err := os.ReadDir(projectPath)

	if err != nil {
		return Project{}, err
	}

	for _, e := range entries {
		if len(e.Name()) < 4 {
			continue
		}

		if name := e.Name()[:len(e.Name())-4]; name == foldeName {
			proj.Name = name
            filename = e.Name()
		}

		if e.Name() == "lyrics.txt" {

			l, err := parser.ParseLyrics(filepath.Join(projectPath, e.Name()))

			if err != nil {
				log.Println(err)
			}

			proj.Lyrics = l

		}
	}

	if proj.Name == "" {
		return Project{}, errors.New(projectPath + " does not contains proper reaper project")
	}

	proj.ID = generateID(proj.Name)

	proj.Path = filepath.Join(projectPath, filename)

	return proj, nil
}

func (p *Project) GetSection(name string) (parser.Section, error) {

	for _, s := range p.Lyrics.Sections {
		if s.Name == name {
			return s, nil
		}
	}

	return parser.Section{}, errors.New("section not found")
}

func (l *Library) SearchProject(query string) []Project {
	var result []Project

	for _, p := range l.Projects {

		if idx := strings.Index(
			strings.ToLower(p.Name),
			strings.ToLower(query),
		); idx != -1 {

			result = append(result, p)
		}
	}

	return result
}

// In my case library is a directory that contais project directoris. Only when
// dir name matches project filename it is a valid entry.
func (l *Library) ParseLibrary() error {
	var result []Project

	enties, err := os.ReadDir(l.Path)

	if err != nil {
		return err
	}

	for _, entry := range enties {

		if !entry.IsDir() {
			continue
		}

		projectPath := filepath.Join(l.Path, entry.Name())

		proj, err := ParseProjectDir(projectPath)

		if err != nil {
			log.Println(err)
			continue
		}

		result = append(result, proj)
	}

	l.Projects = result

	return nil
}

func (l *Library) GetProject(ID string) (Project, error) {

	for _, proj := range l.Projects {
		if ID == proj.ID {
			return proj, nil
		}
	}

	return Project{}, errors.New("nothing found...")
}
