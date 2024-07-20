package reaper

import (
	"log/slog"
	lp "lp"
	"lp/pkg/chords"
	"lp/pkg/lyrics"
	"path"
	"rpp-parser/pkg/parser"

	"github.com/butbkadrug/reaper-lyrics/pkg/generator"
)

// ranges through slice of lines with lyrics and adds chord lines of a same index. then returns new slice of lines
func ApplyChords(ll []lyrics.LyricLine, c chords.Chords) []lyrics.LyricLine {
	for i := range ll {
		if len(c) <= i {
			return ll
		}

		ll[i].Chords = c[i]
	}

	return ll
}

func MarkerToLyricLynes(section lyrics.Section, pos float64, gen *generator.Generator) []lyrics.LyricLine {
	result := make([]lyrics.LyricLine, len(section.Text))
	for i, line := range section.Text {

		result[i] = lyrics.LyricLine{
			ID:       gen.Get(),
			Position: pos + line.Offset,
			Length:   line.Length,
			Text:     line.Text,
		}

		if len(section.Chords) > i {
			result[i].Chords = section.Chords[i]
		}
	}

	return result
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

func GetChordsFromProjectDir(p string) lyrics.Lyrics {
	slog.Debug("loading chords from path", "path", p)
	dir := path.Dir(p)
	l, err := getChordsFromDir(dir)

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

func getChordsFromDir(dirPath string) (lyrics.Lyrics, error) {
	var err error
	files := []string{"chords.txt", "chords"}

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
