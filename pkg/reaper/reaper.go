package reaper

import (
	"log"
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

//TODO: widthdraw all functionality related to tracking changes into it's onw pkg
//probobly all we want to store is a list of files found. Or maybe not even that.
//maybe to pass any file found to teh watcher pkg immidiately

//TODO: find functionality that I do not use and remove it.

//TODO: Mesure how much network I ocupy? Is it me who ruined network on sunday?

func NewProject() *Project {
	log.Println("NewProject invoked")
	config, err := NewDefauiltConfig()

	if err != nil {
		panic("NewProject failed to create a config: " + err.Error())
	}

	return &Project{
		Project:    parser.Project{},
		Path:       config.Path,
		logger:     config.Logger,
		Cursor:     0,
		TrackID:    config.TrackID,
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
	log.Println("NewDefaultConfig invoked")
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
	log.Println("NewProjectWithConfig invoked")
	p := NewProject()
	log.Println("NewProejct returned", p.Stats == nil)

	if cfg == nil {
		p.logger.Error("pkg/reaper/reaper.go NewProjectWithConfig()", "error", "cfg == nil")
	} else {
		p = &Project{
			Path:    cfg.Path,
			logger:  cfg.Logger,
			TrackID: cfg.TrackID,
		}
	}

	subprojects := []*Project{}

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

	for _, sub := range subprojects {
		p.Lines = append(p.Lines, sub.Lines...)
		p.Sections = append(p.Sections, sub.Sections...)
		p.Watchlist = append(p.Watchlist, sub.Watchlist...)
		p.Comments = append(p.Comments, sub.Comments...)
	}

	p.Update(p.Cursor)

	return p

}

// parses .rpp project specified in the confin. Also, will look for the notes,
// chords, lyrics and helper .rpp files. In the project directory.
func NewSubProjectWithConfig(cfg *ProjectConfig) *Project {

	p := NewProject()

	if cfg == nil {
		p.logger.Error("pkg/reaper/reaper.go NewProjectWithConfig()", "error", "cfg == nil")
		panic("NewSubProject invoked with nil config!")
	}

	p = &Project{
		Path:    cfg.Path,
		logger:  cfg.Logger,
		TrackID: cfg.TrackID,
	}

	proj, err := reaper.Parse(cfg.Path)

	if err != nil {
		p.logger.Error("SONG PROJECT: failed to parse project file", "error", err, "path", cfg.Path)
		return p
	}

	p.Project = proj

	//NOTE: Getting helper project, notes, lyrics, chords

	projHelper := GetHelperProjectInProjectDir(cfg.Path)

	if projHelper.Path != "" {
		p.Watchlist = append(p.Watchlist, projHelper.Path)
	}

	projLyrics := GetLyricsFromProjectDir(cfg.Path)
	if projLyrics.Path != "" {
		p.Watchlist = append(p.Watchlist, projLyrics.Path)
	}

	projectNotes := GetNotesFromProjectDir(cfg.Path)

	if projectNotes.Path != "" {
		p.Watchlist = append(p.Watchlist, projectNotes.Path)
	}
	p.Comments = append(p.Comments, &projectNotes)

	projectChords := GetChordsFromProjectDir(cfg.Path)
	if projectChords.Path != "" {
		p.Watchlist = append(p.Watchlist, projectChords.Path)
	}

	// END

	//NOTE: Start mergins all separate pices of data into one data structure
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

			p.Lines = append(p.Lines, ll...)
			sec.Lines = ll
		}

		//if section doesn't contain any text. we will add one empty line that will be whole section
		//long.

		if sec.Lines.Empty() {
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
			sec.Lines = append(sec.Lines, ll...)
			p.Lines = append(p.Lines, ll...)
		}

		p.Sections = append(p.Sections, sec)
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
	p.Cursor = pos

	return true
}

// update cursor position and triggers the update of current item and marker and line
func (p *Project) Update(pos float64) {

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

	for i := len(p.Lines) - 1; i >= 0; i-- {
		line := p.Lines[i]
		if p.Cursor >= line.Position && p.Cursor < line.Position+line.Length {
			p.Line = line
			return true
		}
	}

	//due to small imperfections in positioning sometimes even though cursor is right on the region
	//it's not exactly precisly there. So, we have to do this additional check. But it happens
	//rarly, and i don't need to do this check every loop. only if no lines found, just to make sure
	for _, line := range p.Lines {
		if line.Position-p.Cursor < 0.001 && line.Position-p.Cursor > -0.001 {
			p.Line = line
			return true
		}
	}
	p.Line = lyrics.LyricLine{ID: -1}
	return true
}

// assignes the first section with positions <= and position+length > cursor to a Project.section
// if cursor still in current section no op
// if no section found then empty section is assigned
func (p *Project) UpdateSection() bool {

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
	//if cursor still within the current item we do nothing
	if p.Cursor >= p.Item.Position && p.Cursor <= p.Item.Position+p.Item.Length {
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
