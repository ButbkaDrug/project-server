package reaper

import (
	"log/slog"
	"lp/pkg/lyrics"
	"os"
	"path"
	reaper "rpp-parser"
	"rpp-parser/pkg/parser"
	"time"

	"github.com/butbkadrug/reaper-lyrics/pkg/generator"
)

type ProjectConfig struct {
	Path           string
	Single         bool
	TrackID        int
	ItemID         int
	IDGen          *generator.Generator
	PositionOffset float64
	Logger         *slog.Logger
}

type Project struct {
	//parser project.rpp instance
	parser.Project
	//path to a project itself
	Path string
	//point to a logger
	logger *slog.Logger
	//float represent current position of a cursor
	Cursor float64
	//Target track's id which will be used to determen
	//which track we will scan for subprojects
	TrackID int
	//Shortcut to the target track for ease of use
	Track parser.Track
	//Shortcut to a current item under the cursor
	Item parser.Item
	//Shortcut to the upcoming item in the timeline order
	PeekItem parser.Item
	//Current marker under the cursor
	Marker parser.Marker
	//Upcoming marker in the timeline order
	PeekMarker parser.Marker
	//Current section under the cursor
	Section Section
	//Current line uner the cursor
	Line lyrics.LyricLine
	//Slice of lines which will be rendered to a user
	Lines lyrics.Lines
	//slice of all section, of all subprojects
	Sections []Section
	//slice of pointers to a notes structures
	//used for adding/deleting notes from the web interface
	Comments []*lyrics.Lyrics
	//WARN:Watchlist and stats does not belog here

	//Slice of paths which will be tracked for changes
	//for automatic reload on change
	Watchlist []string
	//map of info about files in the watchlist
	Stats map[string]time.Time
}

func NewProject() *Project {
	return NewProjectWithConfig(&ProjectConfig{})
}

// returns new projectConfig with saint default values. Returns an error if fails to get access to the project file
//
// Defauilt values:
//
// Path = ~/Music/REAPER/BaseProject/BaseProject.rpp
//
// TrackID = 1
//
// ItemID = 1
//
// PositionOffset = 0
//
// New Generator and NewLogger that outputs to os.Stdout.
func NewDefauiltConfig() (*ProjectConfig, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return nil, err
	}

	return &ProjectConfig{
		Path:           path.Join(homeDir, "Music", "REAPER", "BaseProject", "BaseProject.rpp"),
		TrackID:        0,
		ItemID:         0,
		PositionOffset: 0,
		IDGen:          generator.New(),
		Logger:         slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	}, nil
}

func ParseHostProject(config *ProjectConfig) ([]*Project, error) {
	subprojects := []*Project{}
	track := parser.Track{}
	proj, err := reaper.Parse(config.Path)

	if err != nil {
		return []*Project{}, err
	}

	if config.TrackID >= len(proj.Tracks) {
		config.Logger.Error("PROJECT: failed to select track. TrackID is too high", "TrackID", config.TrackID, "Number of tracks", len(proj.Tracks))
	} else {
		track = proj.Tracks[config.TrackID]
	}

	for i, item := range track.Items {

		if source := item.SourceSection.Source; source.Format == "RPP_PROJECT" {

			config.Path = source.File
			config.PositionOffset = item.Position
			config.ItemID = i

			subProject := NewSubProjectWithConfig(config)

			subprojects = append(subprojects, subProject)

		}
	}
	return subprojects, nil
}

func NewProjectWithConfig(cfg *ProjectConfig) *Project {
	if cfg == nil {
		conf, err := NewDefauiltConfig()
		if err != nil {
			panic("no config provided and creation of a default config failed")
		}

		cfg = conf
	}

	p := &Project{
		Project:    parser.Project{},
		Path:       cfg.Path,
		logger:     cfg.Logger,
		Cursor:     0,
		TrackID:    cfg.TrackID,
		Track:      parser.Track{},
		Item:       parser.Item{},
		PeekItem:   parser.Item{},
		Marker:     parser.Marker{},
		PeekMarker: parser.Marker{},
		Section:    Section{},
		Line:       lyrics.LyricLine{},
		Lines:      lyrics.Lines{},
		Sections:   []Section{},
		Comments:   []*lyrics.Lyrics{},
		Watchlist:  []string{},
		Stats:      make(map[string]time.Time),
	}

	subprojects := []*Project{}

	cfg.Logger.Debug("reading host project", "path", cfg.Path, "single", cfg.Single)

	if cfg.Single {
		config, err := NewDefauiltConfig()
		if err != nil {
			cfg.Logger.Error("NewProjectWithConfig() failed to creat a default config", "branch", "single", "error", err.Error())
			panic("NewProjectWithConfig failed to creat default project config")
		}

		config.Path = cfg.Path
		config.IDGen = cfg.IDGen
		config.Logger = cfg.Logger

		subProject := NewSubProjectWithConfig(config)
		subprojects = append(subprojects, subProject)

	} else {
		p.Watchlist = append(p.Watchlist, cfg.Path)
		s, err := ParseHostProject(cfg)
		if err != nil {
			slog.Error("NewHostProject() failed to parse host project", "error", err.Error())
			panic("failed to parse host project!")
		}
		subprojects = append(subprojects, s...)
	}

	p.logger.Debug("finished reading subprojects", "len", len(subprojects))

	for _, sub := range subprojects {
		p.Lines = append(p.Lines, sub.Lines...)
		p.Sections = append(p.Sections, sub.Sections...)
		p.Watchlist = append(p.Watchlist, sub.Watchlist...)
		p.Comments = append(p.Comments, sub.Comments...)
	}

	p.Update(p.Cursor)

	return p

}

func NewSubProjectWithConfig(cfg *ProjectConfig) *Project {
	cfg.Logger.Debug("reading subproject", "path", cfg.Path)
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
	comments := []*lyrics.Lyrics{}
	proj, err := reaper.Parse(cfg.Path)

	if err != nil {
		logger.Error("SONG PROJECT: failed to parse project file", "error", err, "path", cfg.Path)
		return &Project{
			logger: logger,
		}
	}

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
	comments = append(comments, &projectNotes)

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

		//PERF: Maybe instead of mergin comments into one slice. I should keep then in their own structure and check if there is comment on render. Then I will preserve a original structure and if I want to update and save it, it shuld be hips easier because structure knows how and where to save itself!

		//TODO: refactor this pice to keep notes structure within section

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
		Comments:  comments,
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
