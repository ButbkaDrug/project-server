package reaper

import (
	"log/slog"
	"lp/pkg/lyrics"
	"os"
	reaper "rpp-parser"
	"rpp-parser/pkg/parser"
	"time"

	"github.com/butbkadrug/reaper-lyrics/pkg/generator"
)

//TODO: join specific chords and notes into one file. Because they tent to go together anyway
//TODO: also to keep offsets in a separate file, or parse them on load from the file

//TODO: I fill like it will be better to have offsets in a separate file, or parse them off the porject on load entirely. Having offsets in the file with lyrics clutters it alot and makes it hard to do the chords and etc.

type ProjectConfig struct {
	Path           string
	TrackID        int
	ItemID         int
	IDGen          *generator.Generator
	PositionOffset float64
	Logger         *slog.Logger
}

type Project struct {
	parser.Project
	Stats      map[string]time.Time
	Path       string
	logger     *slog.Logger
	Cursor     float64
	TrackID    int
	Track      parser.Track
	Item       parser.Item
	PeekItem   parser.Item
	Marker     parser.Marker
	PeekMarker parser.Marker
	Section    Section
	Line       lyrics.LyricLine
	Lines      lyrics.Lines
	Sections   []Section
	Watchlist  []string
}

func NewProject() *Project {
	return NewProjectWithConfig(&ProjectConfig{})
}

func NewProjectWithConfig(cfg *ProjectConfig) *Project {
	//TODO: refactor this function into smaller chunks where possible

	var lines lyrics.Lines
	var sections []Section

	if cfg == nil {
		cfg = &ProjectConfig{}
	}

	logger := cfg.Logger

	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
		logger.Info("logger is nil: output to stdout")

	}

	g := generator.NewGenerator()
	cfg.IDGen = g

	logger.Debug("NewProjectWithConfig()", "path", cfg.Path)

	watchlist := []string{}
	watchlist = append(watchlist, cfg.Path)
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

	subprojects := []*Project{}

	logger.Debug("subprojects", "total", len(subprojects))

	for i, item := range track.Items {

		if source := item.SourceSection.Source; source.Format == "RPP_PROJECT" {

			watchlist = append(watchlist, source.File)

			cfg.Path = source.File
			cfg.PositionOffset = item.Position
			cfg.ItemID = i

			subProject := NewSubProjectWithConfig(cfg)

			subprojects = append(subprojects, subProject)
			watchlist = append(watchlist, subProject.Watchlist...)

		}
	}

	for _, sub := range subprojects {
		lines = append(lines, sub.Lines...)
		sections = append(sections, sub.Sections...)

	}

	p := &Project{
		Path:      cfg.Path,
		TrackID:   cfg.TrackID,
		Project:   proj,
		Track:     track,
		Lines:     lines,
		Sections:  sections,
		logger:    logger,
		Stats:     make(map[string]time.Time),
		Watchlist: watchlist,
	}

	p.Update(proj.Cursor)

	return p

}

func NewSingleProjectWithConfig(cfg *ProjectConfig) *Project {
	//FIXME: In the mean time I just copied NewProjectWithConfig function but might be worth refactoring repeating stuff into a NewProject function and go from there.
	//TODO: refactor this function into smaller chunks where possible

	var lines lyrics.Lines
	var sections []Section
	var watchlist []string

	if cfg == nil {
		cfg = &ProjectConfig{}
	}

	logger := cfg.Logger

	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
		logger.Info("logger is nil: output to stdout")

	}

	g := generator.NewGenerator()
	cfg.IDGen = g

	logger.Debug("NewProjectWithConfig()", "path", cfg.Path)

	proj := parser.Project{}
	track := parser.Track{}

	subprojects := []*Project{}

	watchlist = append(watchlist, cfg.Path)

	subProject := NewSubProjectWithConfig(&ProjectConfig{
		Path:           cfg.Path,
		IDGen:          g,
		TrackID:        0,
		ItemID:         1,
		Logger:         logger,
		PositionOffset: 0,
	})

	watchlist = append(watchlist, subProject.Watchlist...)

	subprojects = append(subprojects, subProject)

	for _, sub := range subprojects {
		lines = append(lines, sub.Lines...)
		sections = append(sections, sub.Sections...)

	}

	p := &Project{
		Path:      cfg.Path,
		TrackID:   cfg.TrackID,
		Project:   proj,
		Track:     track,
		Lines:     lines,
		Sections:  sections,
		logger:    logger,
		Watchlist: watchlist,
		Stats:     make(map[string]time.Time),
	}

	p.Update(proj.Cursor)

	return p

}

func NewSubProjectWithConfig(cfg *ProjectConfig) *Project {
	//TODO: refactor this mess

	var lines lyrics.Lines
	var sections []Section

	if cfg == nil {
		cfg = &ProjectConfig{}
	}

	logger := cfg.Logger

	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
		logger.Info("logger is nil: output to stdout")
	}

	watchlist := []string{}
	proj, err := reaper.Parse(cfg.Path)

	if err != nil {
		logger.Error("SONG PROJECT: failed to parse project file", "error", err, "path", cfg.Path)
		return &Project{
			logger: logger,
		}
	}
	//FIX: we shell check if there is a markers.rpp and parse it for offsets. If there are notes and chords they should have highest priority, then dedicated file with notes/chords. And then general file with all lyrics. Also maybe to put region specific notes and chords in the same place

	//TODO: Get rid of chords and keep everyhting in notes.
	//TODO: Merge projectNotes & projHelper.

	projHelper := GetHelperProjectInProjectDir(cfg.Path)

	if projHelper.Path != "" {
		watchlist = append(watchlist, projHelper.Path)
	}

	projLyrics := GetLyricsFromProjectDir(cfg.Path)
	if projLyrics.Path != "" {
		watchlist = append(watchlist, projLyrics.Path)
	}

	projectNotes := GetNotesFromProjectDir(cfg.Path)
	if projectNotes.Path != "" {
		watchlist = append(watchlist, projectNotes.Path)
	}

	projectChords := GetChordsFromProjectDir(cfg.Path)
	if projectChords.Path != "" {
		watchlist = append(watchlist, projectChords.Path)
	}

	for i, m := range proj.Markers {
		//skip marker if it's not a region
		if !m.IsRegion {
			continue
		}

		//NOTE: I already messed up position calculation couple of times. And at this moment it is imposibble to test it. Therefour we need a testable function for this reason. I should use presicion math function for comparasment.

		sec := Section{
			ID:       m.ID,
			Position: m.Position + cfg.PositionOffset,
			Lenght:   proj.Markers[i+1].Position - m.Position,
			Name:     parseSectionName(m.Name),
		}

		sec.ItemID = cfg.ItemID
		sec.Comments = projectNotes.GetNotesWithNumber(m.ID)
		sec.Comments = append(sec.Comments, projHelper.GetNotesWithID(m.ID)...)

		chords := projectChords.GetChordsWithID(m.ID)

		for _, section := range projLyrics.Sections {
			if section.GetFullName() != m.Name {
				continue
			}

			section.ID = m.ID
			sec.Comments = append(sec.Comments, section.Comments...)

			offsets := projHelper.GetOffsets(section)

			if offsets.Empty() {
				offsets = CalculateOffsets(section.Text, proj.Markers, i)
			}

			section.Text = ApplyOffsets(section.Text, offsets)

			ll := MarkerToLyricLynes(section, m.Position+cfg.PositionOffset, cfg.IDGen)

			if chords.Empty() {
				chords = projLyrics.GetChordsForSection(section)
			}

			//FIXME: there should be a better way to do this
			if chords.Empty() {
				chords = projectChords.GetChordsForSection(section)

			}

			if !chords.Empty() {
				ll = ApplyChords(ll, chords)
			}

			lines = append(lines, ll...)
			sec.Lines = ll
		}

		//if section doesn't contain any text. we will add one empty line that will be whole section
		//long.

		if sec.Lines.Empty() {
			logger.Debug("section has no lyrics", "s", sec.Name)
			ll := []lyrics.LyricLine{}

			if len(chords) < 1 {

				ll = append(ll, lyrics.LyricLine{
					ID:       cfg.IDGen.Get(),
					Position: m.Position + cfg.PositionOffset,
					Length:   proj.Markers[i+1].Position - m.Position,
					Text:     "",
				})
			}
			for _, c := range chords {
				ll = append(ll, lyrics.LyricLine{
					ID:       cfg.IDGen.Get(),
					Position: m.Position + cfg.PositionOffset,
					Length:   proj.Markers[i+1].Position - m.Position,
					Chords:   c,
					Text:     "",
				})

			}
			logger.Debug("appending", "total", len(ll), "lines", ll)
			sec.Lines = append(sec.Lines, ll...)
			lines = append(lines, ll...)
		}

		sections = append(sections, sec)
	}

	p := &Project{
		Path:      cfg.Path,
		TrackID:   cfg.TrackID,
		Project:   proj,
		Track:     parser.Track{},
		Lines:     lines,
		Sections:  sections,
		logger:    logger,
		Watchlist: watchlist,
	}

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

	p.UpdateCursor(pos)
	p.UpdateLine()
	p.UpdateSection()
	if !p.UpdateItem() {
		return
	}
	p.UpdatePeekItem()
	p.UpdateMarker()
	p.UpdatePeekMarker()
}

func (p *Project) UpdateLine() bool {
	p.logger.Debug("UpdateLine() initiated", "cursor", p.Cursor, "current line", p.Line, "total lines", len(p.Lines))
	if next := p.Line.ID + 1; next < len(p.Lines) {
		p.logger.Debug("UpdateLine()", "next", p.Lines[next])
	}

	for i := len(p.Lines) - 1; i >= 0; i-- {
		line := p.Lines[i]
		p.logger.Debug("searching", "line", i, "total", len(p.Lines), "id", line.ID, "cursor", p.Cursor, "pos", line.Position, "len", line.Length)
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
	p.Line = lyrics.LyricLine{ID: -1}

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
