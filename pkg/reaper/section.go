package reaper

import (
	"fmt"
	"lp/pkg/lyrics"
)

type Section struct {
	ID       int
	Position float64
	Lenght   float64
	Name     Name
	Lines    lyrics.Lines
	Comments []string
	//data needed to renderer
	ParentName string
	ItemID     int
	Elapsed    bool
	Active     bool
}

func (s *Section) ShortName() string {
	name := fmt.Sprintf("%c", s.Name.Name[0])

	if s.Name.Number > 0 {
		name += fmt.Sprintf("%d", s.Name.Number)
	}

	return name
}

func (s *Section) Update(pos float64) Section {
	s.Active = pos >= s.Position && pos < s.Position+s.Lenght
	s.Elapsed = pos > s.Position+s.Lenght

	return *s
}
