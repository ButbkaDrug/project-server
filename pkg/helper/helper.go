package helper

import (
	"fmt"
	"log/slog"
	"lp/pkg/chords"
	"lp/pkg/comments"
	"lp/pkg/lyrics"
	"os"
	"path"
	rp "rpp-parser"
	"rpp-parser/pkg/parser"
	"slices"
	"strings"

	"github.com/shopspring/decimal"
)

type Position struct {
	Offset float64
	Length float64
}

type helper struct {
	Path struct {
		Project string
		Base    string
	}
	Proj    *parser.Project
	logfile *os.File
	logger  *slog.Logger
	Lyrics  *lyrics.Lyrics
}

// expects path to a project file to parse
func NewHelper(p string) (*helper, error) {
	f, err := os.Create("/tmp/helper.log")

	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	if p == "" {
		return nil, fmt.Errorf("NewHelper: path must be provided!")
	}

	proj, err := rp.Parse(p)

	if err != nil {
		return nil, err
	}

	return &helper{
		Path: struct {
			Project string
			Base    string
		}{
			Project: p,
			Base:    path.Dir(p),
		},
		Proj:    &proj,
		logfile: f,
		logger:  logger,
		Lyrics:  &lyrics.Lyrics{},
	}, nil

}

func (h *helper) Parse() lyrics.Lyrics {
	h.ParseOffsets()
	h.parseChords()
	h.parseNotes()

	return *h.Lyrics

}

func (h *helper) ParseOffsets() {
	h.logger.Debug("parsing offsets:", "path", h.Path.Base)

	for i, m := range h.Proj.Markers {
		if !m.IsEnd {
			continue
		}

		e := h.Proj.Markers[i-1]
		h.logger.Debug("processing region", "r", e)
		offsets := h.GetOffsetsForRange(e.Position, m.Position)

		h.logger.Debug("looking for long offset", "offsets", offsets, "id", e.ID)

		//NOTE: if no offsets found for the given section we skip to the next section
		if len(offsets) < 1 {
			h.logger.Debug("no offsets... continue to the next region")
			continue
		}

		targetID := -1

		for j, s := range h.Lyrics.Sections {
			//NOTE: Compare section name from rpp file with section name from lyrics file
			//to make sure we update right one
			if e.ID != s.ID {
				continue
			}

			h.logger.Debug("found assosiated section", "section", s, "target", targetID)
			targetID = j

		}

		if targetID < 0 {
			h.Lyrics.Sections = append(h.Lyrics.Sections, lyrics.Section{
				Name: e.Name,
				ID:   e.ID,
			})

			targetID = len(h.Lyrics.Sections) - 1
		}

		text := make([]lyrics.TextLine, len(offsets))

		for k, offset := range offsets {
			text[k].Offset = offset.Offset
			text[k].Length = offset.Length
		}

		h.Lyrics.Sections[targetID].Text = text
	}
}

func (h *helper) parseNotes() {
	notesMarkers := h.FilterMarkersWithPrefix(h.Proj.Markers, "NOTE:")

	if len(notesMarkers) < 1 {
		return
	}

	for _, m := range notesMarkers {
		comment := comments.NewComment()

		note := strings.TrimPrefix(m.Name, "NOTE:")

		noteType, text, typed := strings.Cut(note, ":")

		if !typed {
			comment.Text = note
		} else {
			comment.Type = comments.GetCommentType(noteType)
			text = strings.Trim(text, " ")
			comment.Text = text
		}

		region := h.getAssociatedRegion(m)

		targetID := -1

		for j := range h.Lyrics.Sections {
			if h.Lyrics.Sections[j].ID != region.ID {
				continue
			}

			targetID = j

		}

		if targetID < 0 {
			h.Lyrics.Sections = append(h.Lyrics.Sections, lyrics.Section{
				ID:   region.ID,
				Name: region.Name,
			})

			targetID = len(h.Lyrics.Sections) - 1
		}
		h.Lyrics.Sections[targetID].Comments = append(h.Lyrics.Sections[targetID].Comments, comment)

	}
}

func (h *helper) parseChords() {
	chordsMarkers := h.FilterMarkersWithPrefix(h.Proj.Markers, "CHORDS:")
	if len(chordsMarkers) < 1 {
		h.logger.Info("no chords markers found")
		return
	}

	for _, m := range chordsMarkers {
		//NOTE: parse chords marker.
		//1. Remove prefix
		//2. Split by space
		//3. Further processing as needed

		line := strings.TrimPrefix(m.Name, "CHORDS:")
		line = strings.Trim(line, " ")
		rawChords := strings.Split(line, " ")
		c := make(chords.ChordLine, len(rawChords))

		for i, rc := range rawChords {

			c[i] = chords.NewChord(rc)

		}

		region := h.getAssociatedRegion(m)
		targetID := -1

		for j := range h.Lyrics.Sections {
			if h.Lyrics.Sections[j].ID != region.ID {
				continue
			}

			targetID = j

		}

		if targetID < 0 {
			h.Lyrics.Sections = append(h.Lyrics.Sections, lyrics.Section{
				ID:   region.ID,
				Name: region.Name,
			})
			targetID = len(h.Lyrics.Sections) - 1
		}
		h.Lyrics.Sections[targetID].Chords = append(h.Lyrics.Sections[targetID].Chords, c)
	}

}

// serching for a region that will contain niddle. if no region contain niddle, function returns niddle
func (h *helper) getAssociatedRegion(niddle parser.Marker) parser.Marker {
	for i, m := range h.Proj.Markers {
		if !m.IsRegion {
			continue
		}

		if h.isGreaterThenOrEqual(niddle.Position, m.Position) && h.isLessThen(niddle.Position, h.Proj.Markers[i+1].Position) {
			return m
		}
	}

	return niddle
}

// will return slice of offsets and length for markers in sepcified range.
func (h *helper) GetOffsetsForRange(start, stop float64) []Position {
	result := []Position{}

	markers := h.FilterMarkers([]string{"", "START", "EOL"})
	h.logger.Debug("Getting offsets", "found", len(markers))

	idx := 0

	for i, m := range markers {
		if h.isGreaterThenOrEqual(m.Position, start) &&
			h.isLessThen(m.Position, stop) &&
			m.Name != "START" {

			idx = i

			//NOTE: if name of the marker is START. it means that next section lyrics begin at the
			//end of the current section. and thus should have a negative offset. START only should
			//appear as very last marker in the section

			if i > 0 && markers[i-1].Name == "START" {
				startMarker := markers[i-1]
				h.logger.Debug("processing start marker", "marker", startMarker, "start of a region", start)
				offset := h.SubtrackFloats64(startMarker.Position, start)

				h.logger.Debug("calculating negative offset", "pos", m.Position, "startPos", startMarker.Position, "offset", offset)
				result = append(result, Position{Offset: offset})
			}

			offset := h.SubtrackFloats64(m.Position, start)

			//NOTE: if last element has length 0. we will automatically assign it the lengh up to
			//current element
			if len(result) > 0 && result[len(result)-1].Length == 0 {
				pos := result[len(result)-1].Offset
				result[len(result)-1].Length = h.SubtrackFloats64(offset, pos)
			}

			//NOTE: if name of the marker is EOL(End Of Line) it suggests to only update the length
			//of a previous element and to NOT add a new one.
			if m.Name == "EOL" || m.Name == "START" {
				continue
			}

			result = append(result, Position{Offset: offset})
		}
	}

	//NOTE: if last element has 0 lenth we will automaically finish it at the end of the section
	if len(result) > 0 && result[len(result)-1].Length == 0 {
		result[len(result)-1].Length = h.SubtrackFloats64(stop, start+result[len(result)-1].Offset)

		marker := markers[idx]

		if len(markers) <= idx+1 {
			return result
		}

		marker = markers[idx+1]

		if marker.Name == "EOL" {
			addedLength := h.SubtrackFloats64(marker.Position, stop)
			h.logger.Debug("EOL as first marker. Updating last marker", "original", result[len(result)-1].Length, "additional", addedLength)
			// result[len(result)-1].Length = result[len(result)-1].Length + addedLength
			result[len(result)-1].Length = 7777777
		}
	}

	return result
}

// checks if x greater of equal to y wiht presicion of 9 digits
func (h *helper) isGreaterThenOrEqual(x, y float64) bool {
	decX := decimal.NewFromFloatWithExponent(x, -10)
	decY := decimal.NewFromFloatWithExponent(y, -10)

	return decX.GreaterThanOrEqual(decY)
}

func (h *helper) isLessThen(x, y float64) bool {
	decX := decimal.NewFromFloatWithExponent(x, -10)
	decY := decimal.NewFromFloatWithExponent(y, -10)

	return decX.LessThan(decY)
}

// subtrackts y from x with fixed persision of 10 decimals
func (h *helper) SubtrackFloats64(x, y float64) float64 {
	decPos := decimal.NewFromFloatWithExponent(x, -10)
	decStart := decimal.NewFromFloatWithExponent(y, -10)

	return decPos.Sub(decStart).InexactFloat64()
}

// will return only actual markers, not regions and the end markers. Also will filter out results
// we not interested in like markers with undefined name constants
func (h *helper) FilterMarkers(whitelist []string) []parser.Marker {

	var result []parser.Marker

	for _, m := range h.Proj.Markers {
		if !slices.Contains(whitelist, m.Name) {

			continue
		}

		if !m.IsEnd && !m.IsRegion {

			result = append(result, m)
		}
	}

	return result
}

// returns a new slice with markers that has prefix "s". Uses HasPrefix function from strings library
func (h *helper) FilterMarkersWithPrefix(markers []parser.Marker, s string) []parser.Marker {

	var result []parser.Marker

	for _, m := range markers {
		h.logger.Debug("filtering markers with prefix", "prefix", s, "marker", m)
		if !strings.HasPrefix(m.Name, s) {

			continue
		}

		if !m.IsEnd && !m.IsRegion {

			result = append(result, m)
		}
	}

	return result
}
