package reaper

import (
	"fmt"
	"log"
	"log/slog"
	lp "lp"
	"lp/pkg/lyrics"
	"os"
	"path"
	reaper "rpp-parser"
	"rpp-parser/pkg/parser"
	"strconv"
	"strings"
	"time"

	"sse"
)

// FIXME: refactoring Project.
// I want to make a list of all lines with all id's in it and then render them out.
// Also, I would relly want to switch to HTMX sse processing. So that I can send over a state and
// renew view on the client. So, that in case I'm doint any changes in reaper they would
// immidiatly show throw on the client

// TODO: Is it possible to have dynamic targets? So that I can replace only one line of text, not
// the whole page.

type Project struct {
	parser.Project
	Stats  os.FileInfo
	Path   string
	logger *slog.Logger

	//NOTE:Refactoring section

	Cursor     float64
	TrackID    int
	Track      parser.Track
	Item       parser.Item
	PeekItem   parser.Item
	Marker     parser.Marker
	PeekMarker parser.Marker
	Section    Section
	Line       LyricLine
	Lines      Lines
	Sections   []Section
}

type LyricLine struct {
	ID       int
	Position float64
	Length   float64
	Text     string
}

type Lines []LyricLine

type ProjectConfig struct {
	Path    string
	TrackID int
	Logger  *slog.Logger
}

type Section struct {
	ID       int
	Position float64
	Lenght   float64
	Name     Name
	Lines    Lines
	Comments []string
	//data needed to renderer
	Elapsed bool
	Active  bool
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

type Name struct {
	Name   string
	Number int
	Lang   string
}

type Offset struct {
	Offset float64
	Length float64
}

type OffsetMap []Offset

func (o OffsetMap) Empty() bool {

	for _, offset := range o {
		if offset.Offset != offset.Length {
			return false
		}
	}

	return true
}

func NewOffsetMap(s lyrics.Section) OffsetMap {
	offsets := make([]Offset, len(s.Text))

	for i, line := range s.Text {
		offsets[i].Offset = line.Offset
		offsets[i].Length = line.Length
	}

	return offsets
}

func GetOffsets(l lyrics.Lyrics, s lyrics.Section) OffsetMap {
	offsets := NewOffsetMap(s)

	if !offsets.Empty() {
		return offsets
	}

	alt := GetAlternativeSection(l, s)

	offsets = NewOffsetMap(alt)
	if !offsets.Empty() {
		return offsets
	}

	rel := GetRelativeSection(l, s)
	offsets = NewOffsetMap(rel)
	if !offsets.Empty() {
		return offsets
	}

	relalt := GetRelativeAlternativeSection(l, s)
	offsets = NewOffsetMap(relalt)
	if !offsets.Empty() {
		return offsets
	}

	all := GetAllRelativeSections(l, s)

	for _, relative := range all {
		offsets := NewOffsetMap(relative)

		if !offsets.Empty() {
			return offsets
		}
	}

	return offsets

}

// looking for the first section with the same name and number, but not the language and returns it
// if no matches found returns original section
func GetAlternativeSection(l lyrics.Lyrics, s lyrics.Section) lyrics.Section {
	for _, alt := range l.Sections {
		if s.Name == alt.Name && s.Number == alt.Number && s.Lang != alt.Lang {
			return alt
		}
	}

	return s
}

// looking for the first section with the same name and language but not number returns it
// if no matches found returns original section
func GetRelativeSection(l lyrics.Lyrics, s lyrics.Section) lyrics.Section {
	for _, alt := range l.Sections {
		if s.Name == alt.Name && s.Number != alt.Number && s.Lang == alt.Lang {
			return alt
		}
	}

	return s
}

// lokking for a section with a same name but not the number and not the language
func GetRelativeAlternativeSection(l lyrics.Lyrics, s lyrics.Section) lyrics.Section {
	for _, alt := range l.Sections {
		if s.Name == alt.Name && s.Number != alt.Number && s.Lang != alt.Lang {
			return alt
		}
	}

	return s
}

// looking for all section with a matching name
func GetAllRelativeSections(l lyrics.Lyrics, s lyrics.Section) []lyrics.Section {
	result := []lyrics.Section{}
	for _, section := range l.Sections {
		if section.Name == s.Name {
			result = append(result, section)
		}
	}

	return result
}

func CalculateOffsets(text []lyrics.TextLine, markers []parser.Marker, i int) OffsetMap {
	offsets := make(OffsetMap, len(text))

	for k := range text {
		m := markers[i]

		length := (markers[i+1].Position - m.Position) / float64(len(text))
		offset := float64(k) * length

		offsets[k] = Offset{offset, length}

	}

	return offsets
}
func ApplyOffsets(section []lyrics.TextLine, offsets OffsetMap) []lyrics.TextLine {
	for i := range section {

		if i >= len(offsets) {
			continue
		}

		section[i].Offset = offsets[i].Offset
		section[i].Length = offsets[i].Length

	}

	return section
}

func ConverToLyricLynes(lines Lines, section []lyrics.TextLine, m parser.Marker) []LyricLine {
	result := make([]LyricLine, len(section))
	for i, line := range section {

		result[i] = LyricLine{
			ID:       len(lines) + i,
			Position: m.Position + line.Offset,
			Length:   line.Length,
			Text:     line.Text,
		}
	}
	return result
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

func GetLyricsForItem(i parser.Item) lyrics.Lyrics {
	return GetLyricsFromProjectDir(i.SourceSection.Source.File)
}

func GetNotesForItem(i parser.Item) lyrics.Lyrics {
	return GetNotesFromProjectDir(i.SourceSection.Source.File)
}

func GetLyricsFromProjectDir(p string) lyrics.Lyrics {
	dir := path.Dir(p)
	l, err := getLyrics(dir)

	if err != nil {
		slog.Error("GetLyricsFromProjectDir(): failed to parse lyrics file", "path", dir, "error", err.Error())
	}

	return l
}

func GetNotesFromProjectDir(p string) lyrics.Lyrics {
	dir := path.Dir(p)
	l, err := getNotes(dir)

	if err != nil {
		slog.Error("GetLyricsFromProjectDir(): failed to parse lyrics file", "path", dir, "error", err.Error())
	}

	return l
}

func getLyrics(dirPath string) (lyrics.Lyrics, error) {
	var err error
	files := []string{"song.txt", "lyrics.txt", "lyrics-backup"}

	for _, file := range files {
		l, e := lp.Parse(path.Join(dirPath, file))

		if e != nil {
			err = e

			continue
		}

		return l, nil

	}

	return lyrics.Lyrics{}, err
}

func getNotes(dirPath string) (lyrics.Lyrics, error) {
	var err error
	files := []string{"notes.txt", "notes"}

	for _, file := range files {
		l, e := lp.Parse(path.Join(dirPath, file))

		if e != nil {
			err = e

			continue
		}

		return l, nil

	}

	return lyrics.Lyrics{}, err
}

func NewSectionFromMarker(m parser.Marker, end parser.Marker) Section {

	return Section{
		ID:       m.ID,
		Position: m.Position,
		Lenght:   end.Position - m.Position,
		Name:     parseSectionName(m.Name),
	}
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

// NOTE: At the moment i assume that i would always get a host project, meaning the project, that
// contains subprojects with actual songs. for the ease of use i want to be able to open a particular
// song and have it behave normally. I would need to lookup words in current directory.
func NewProject() *Project {
	return NewProjectWithConfig(&ProjectConfig{})
}

func NewProjectWithConfig(cfg *ProjectConfig) *Project {
	//TODO: refactor this function into smaller chunks where possible

	var lines Lines
	var sections []Section

	if cfg == nil {
		cfg = &ProjectConfig{}
	}

	logger := cfg.Logger

	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	}

	proj, err := reaper.Parse(cfg.Path)

	if err != nil {
		logger.Error("PROJECT: failed to parse project file", "error", err, "path", cfg.Path)
		return &Project{
			logger: logger,
		}
	}

	track := parser.Track{}
	if cfg.TrackID >= len(proj.Tracks) {
		logger.Error("PROJECT: failed to select track. TrackID is too high", "TrackID", cfg.TrackID, "Number of tracks", len(proj.Tracks))
	} else {
		track = proj.Tracks[cfg.TrackID]
	}

	itemsLyrics := make([]lyrics.Lyrics, len(track.Items))

	for i, item := range track.Items {
		itemsLyrics[i] = GetLyricsForItem(item)
	}

	//NOTE: I want to build fuill map right at init. Every time project file changes, we will
	//reparse it. No need to calculate what section we at. We should be able to use timestamp to
	//to find active line and the just scroll to it. I should assing id's to the line right at this
	//stage so that rendered will only have to render, and endpoint would only have to serve the
	//data

	for i, m := range proj.Markers {
		if !m.IsRegion {
			continue
		}

		idx := -1
		// Now I need to get item's id. So that i can match names and get the lyrics
		for j, item := range track.Items {
			if m.Position >= item.Position && m.Position < item.Position+item.Length {
				idx = j
				break
			}
		}

		// if we didn't match any item. idex will remain negative
		if idx < 0 {
			continue
		}

		//NOTE: Maybe during the section constructon I should place native and alternative text into
		//the section, so that if i'm changing the language on the fly, I don't have to re calculate
		//anythins???
		sec := NewSectionFromMarker(m, proj.Markers[i+1])

		for _, section := range itemsLyrics[idx].Sections {
			if section.GetFullName() != m.Name {
				continue
			}

			sec.Comments = section.Comments

			offsets := GetOffsets(itemsLyrics[idx], section)

			if offsets.Empty() {
				offsets = CalculateOffsets(section.Text, proj.Markers, i)
			}

			section.Text = ApplyOffsets(section.Text, offsets)

			ll := ConverToLyricLynes(lines, section.Text, m)

			lines = append(lines, ll...)
			sec.Lines = ll
		}

		//if section doesn't contain any text. we will add one empty line that will be whole section
		//long.
		if len(sec.Lines) == 0 {
			ll := LyricLine{
				ID:       len(lines),
				Position: m.Position,
				Length:   proj.Markers[i+1].Position - m.Position,
				Text:     "",
			}
			sec.Lines = append(sec.Lines, ll)
			lines = append(lines, ll)
		}

		sections = append(sections, sec)
	}

	p := &Project{
		Path:     cfg.Path,
		TrackID:  cfg.TrackID,
		Project:  proj,
		Track:    track,
		Lines:    lines,
		Sections: sections,
		logger:   logger,
	}

	stats, _ := p.getStats()

	p.Stats = stats
	p.Update(proj.Cursor)

	return p

}

//FIXME:get notes with specific ID

func GetNotesWithID(notes lyrics.Lyrics, ID int) []string {
	var result []string
	for _, note := range notes.Sections {
		if note.Number == ID {
			for _, line := range note.Text {
				result = append(result, line.Text)
			}
		}
	}

	return result
}

func NewSongProjectWithConfig(cfg *ProjectConfig) *Project {

	var lines Lines
	var sections []Section

	if cfg == nil {
		cfg = &ProjectConfig{}
	}

	logger := cfg.Logger

	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	}

	proj, err := reaper.Parse(cfg.Path)

	if err != nil {
		logger.Error("SONG PROJECT: failed to parse project file", "error", err, "path", cfg.Path)
		return &Project{
			logger: logger,
		}
	}

	projLyrics := GetLyricsFromProjectDir(cfg.Path)
	projectNotes := GetNotesFromProjectDir(cfg.Path)

	logger.Info("Notes loaded", "notes", projectNotes)

	for i, m := range proj.Markers {
		if !m.IsRegion {
			continue
		}

		sec := NewSectionFromMarker(m, proj.Markers[i+1])
		sec.Comments = GetNotesWithID(projectNotes, m.ID)

		for _, section := range projLyrics.Sections {
			if section.GetFullName() != m.Name {
				continue
			}

			offsets := GetOffsets(projLyrics, section)

			if offsets.Empty() {
				offsets = CalculateOffsets(section.Text, proj.Markers, i)
			}

			section.Text = ApplyOffsets(section.Text, offsets)

			ll := ConverToLyricLynes(lines, section.Text, m)

			lines = append(lines, ll...)
			sec.Lines = ll
		}

		//if section doesn't contain any text. we will add one empty line that will be whole section
		//long.
		if len(sec.Lines) == 0 {
			ll := LyricLine{
				ID:       len(lines),
				Position: m.Position,
				Length:   proj.Markers[i+1].Position - m.Position,
				Text:     "",
			}
			sec.Lines = append(sec.Lines, ll)
			lines = append(lines, ll)
		}

		sections = append(sections, sec)
	}

	p := &Project{
		Path:     cfg.Path,
		TrackID:  cfg.TrackID,
		Project:  proj,
		Track:    parser.Track{},
		Lines:    lines,
		Sections: sections,
		logger:   logger,
	}

	stats, _ := p.getStats()

	p.Stats = stats
	p.Update(proj.Cursor)

	return p
}

func (p *Project) GetMapForCurrentItem() []Section {
	var result []Section

	for _, s := range p.Sections {
		if s.Position >= p.Item.Position && s.Position < p.Item.Position+p.Item.Length {
			result = append(result, s)
		}
	}

	return result
}

func (p *Project) UpdateCursor(pos float64) bool {
	p.logger.Debug("UpdateCurosr()", "cursor", pos)
	p.Cursor = pos

	return true
}

// update cursor position and triggers the update of current item and marker and line
func (p *Project) Update(pos float64) {
	p.logger.Debug("UPDATE()", "new position", pos)

	//order is important!
	for i, s := range p.Sections {
		p.Sections[i] = s.Update(pos)
	}

	p.UpdateMarker()
	p.UpdatePeekMarker()

	if !p.UpdateCursor(pos) {
		return
	}

	if !p.UpdateLine() {
		return
	}
	if !p.UpdateSection() {
		return
	}

	if !p.UpdateItem() {
		return
	}
	p.UpdatePeekItem()
}

func (p *Project) UpdateLine() bool {
	p.logger.Debug("UpdateLine()")

	if p.Cursor >= p.Line.Position && p.Cursor < p.Line.Position+p.Line.Length {
		return false
	}

	for i := len(p.Lines) - 1; i >= 0; i-- {
		line := p.Lines[i]
		if p.Cursor >= line.Position && p.Cursor < line.Position+line.Length {
			p.Line = line
			p.logger.Debug("UpdateLine()", "cursor", p.Cursor, "currentline", p.Line)
			return true
		}
	}

	//due to small imperfections in positioning sometimes even though cursor is right on the region
	//it's not exactly precisly there. So, we have to do this additional check. But it happens
	//rarly, and i don't need to do this check every loop. only if no lines found, just to make sure
	for _, line := range p.Lines {
		if line.Position-p.Cursor < 0.001 && line.Position-p.Cursor > -0.001 {
			p.logger.Debug("found close enough line", "cursor", p.Cursor, "lineID", line.ID, "Position", line.Position, "differance", line.Position-p.Cursor)

			p.Line = line
			return true
		}
	}

	p.logger.Debug("line is not found", "cursor", p.Cursor)
	p.Line = LyricLine{ID: -1}

	return true
}

// assignes the first section with positions <= and position+length > cursor to a Project.section
// if cursor still in current section no op
// if no section found then empty section is assigned
func (p *Project) UpdateSection() bool {
	p.logger.Debug("UpdateSection()")

	if p.Cursor >= p.Section.Position && p.Cursor < p.Section.Position+p.Section.Lenght {
		return false
	}

	for _, s := range p.Sections {
		if p.Cursor >= s.Position && p.Cursor < s.Position+s.Lenght {
			p.Section = s
			return true
		}
	}

	p.Section = Section{}

	return true
}

// selects the track with given index. if track not exists no op
func (p *Project) UpdateTrack(idx int) {

	if idx < 0 || idx >= len(p.Tracks) {
		return
	}

	p.TrackID = idx

}

// if cursor is still in range of current item and peekItem no op. otherwise searches for the first item to
// contain the cursor and "selects" it
func (p *Project) UpdateItem() bool {
	p.logger.Debug("UpdateItem()")

	//if cursor still within the current item we do nothing
	if p.Cursor >= p.Item.Position && p.Cursor <= p.Item.Position+p.Item.Length {
		p.logger.Debug("UpdateItem() item not changed")
		return false
	}

	for _, item := range p.Track.Items {
		//if cursor past the item starting position and didn't reach the end position we make it
		//an active item
		if p.Cursor >= item.Position && p.Cursor < item.Position+item.Length {
			p.Item = item
			return true
		}
	}

	p.Item = parser.Item{}
	return true
}

func (p *Project) UpdatePeekItem() {
	p.logger.Debug("UpdatePeekItem()")

	for _, item := range p.Track.Items {
		if item.Position > p.Cursor {
			p.PeekItem = item
			return
		}
	}
	p.PeekItem = parser.Item{}
}

// looking for a first marker to contain a cursor and selects it. If cursor still in range of the
// current marker no op
func (p *Project) UpdateMarker() {
	p.logger.Debug("UpdateMarker()")

	for i, m := range p.Markers {
		if !m.IsEnd {
			continue
		}

		if p.Cursor < m.Position {
			p.Marker = p.Markers[i-1]
		}
	}
}

func (p *Project) UpdatePeekMarker() {
	p.logger.Debug("UpdatePeekMarker()")
	for _, m := range p.Markers {
		if !m.IsRegion {
			continue
		}

		if m.Position > p.Cursor {
			p.PeekMarker = m
			return
		}
	}

	p.PeekMarker = parser.Marker{}
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
			p.logger.Error("Project.Watch(): failed to get stats of a project file", "error", err)
			return err
		}

		if info.ModTime() != p.Stats.ModTime() && info.Size() != p.Stats.Size() {

			go clients.Broadcast(true)

			proj, err := reaper.Parse(p.Path)

			if err != nil {
				log.Println("parsing error: ", err)
			}
			p.Project = proj

			p.Stats = info
		}

		time.Sleep(time.Second)

	}
}

//WARNING: Reaper on windows and on linux save region info differently
// current parser works for files saved under linux

func GetImageName(n string) string {
	toks := strings.Split(n, "-")
	name := strings.Trim(toks[0], " \t\r\n")
	name = strings.ReplaceAll(name, " ", "-")
	return strings.ToLower(name) + ".jpg"
}
