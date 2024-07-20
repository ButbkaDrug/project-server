package reaper

import (
	"fmt"
	"strconv"
	"strings"
)

type Name struct {
	Name   string
	Number int
	Lang   string
}

func (n Name) String() string {

	text := n.Name

	if n.Number > 0 {
		text += fmt.Sprintf(" %d", n.Number)
	}

	if n.Lang != "" {
		text += fmt.Sprintf(" %s", n.Lang)
	}

	return text

}

func parseSectionName(name string) Name {
	tokens := strings.Split(name, " ")
	n := Name{}

	for _, t := range tokens {
		if len(t) > 2 {
			if n.Name == "" {
				n.Name = t
				continue
			}

			n.Name += " " + t
			continue
		}

		num, err := strconv.Atoi(t)

		if err != nil && n.Lang == "" {
			n.Lang = t
			continue
		}

		n.Number = num
	}

	return n

}

func (n *Name) GetFullName() string {
	name := n.Name
	if n.Number > 0 {
		name += fmt.Sprintf(" %d", n.Number)
	}

	if n.Lang != "" {
		name += fmt.Sprintf(" %s", n.Lang)
	}

	return name
}
