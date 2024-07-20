package reaper

import (
	"lp/pkg/lyrics"
	"rpp-parser/pkg/parser"
)

func CalculateOffsets(text []lyrics.TextLine, markers []parser.Marker, i int) lyrics.OffsetMap {
	offsets := make(lyrics.OffsetMap, len(text))

	for k := range text {
		m := markers[i]

		length := (markers[i+1].Position - m.Position) / float64(len(text))
		offset := float64(k) * length

		offsets[k] = lyrics.Offset{
			Offset: offset,
			Length: length,
		}

	}

	return offsets
}

func ApplyOffsets(section []lyrics.TextLine, offsets lyrics.OffsetMap) []lyrics.TextLine {
	for i := range section {

		if i >= len(offsets) {
			continue
		}

		section[i].Offset = offsets[i].Offset
		section[i].Length = offsets[i].Length

	}

	return section
}
