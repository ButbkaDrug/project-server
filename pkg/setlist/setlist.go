package setlist

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Setlist struct {
	Path  string
	songs []Song
}

type Song struct {
	ID   int
	Name string
	Path string
}

func NewSong(id int, p string) Song {
	name := filepath.Base(p)
	name = strings.Split(name, "-")[0]
	name = strings.Trim(name, " ")

	return Song{
        ID: id,
		Name: name,
		Path: p,
	}
}

func NewSetlist(p string) *Setlist {
	return &Setlist{
		Path:  p,
		songs: []Song{},
	}
}

func (s *Setlist)Contains(p string) bool {
    for _, v := range s.songs {
        if v.Path == p {
            return true
        }
    }

    return false
}

func (s *Setlist) Read() error {

	f, err := os.Open(s.Path)

	if err != nil {
		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	songs := []Song{}

	for scanner.Scan() {
		songs = append(songs, NewSong(len(songs), scanner.Text()))
	}

	s.songs = songs

	return nil

}

func (s *Setlist) Write() error {
	f, err := os.Create(s.Path)
	if err != nil {
		return err
	}

	defer f.Close()

	var data string

	for _, song := range s.songs {
		data += song.Path + "\n"
	}

	_, err = f.WriteString(data)

	if err != nil {
		return err
	}

	return nil
}

func (s *Setlist) NewSetlist(p string) *Setlist {

	return &Setlist{
		Path:  p,
		songs: []Song{},
	}

}

func (s *Setlist) Append(p string) (Song, error) {

    log.Println("append happened")
    song := NewSong(len(s.songs), p)

    log.Println(song)

	s.songs = append(s.songs, song)

	return song, s.Write()
}

func (s *Setlist) Delete(idx int) error {

	if idx < 0 || idx > len(s.songs) {
		return nil
	}

	songs := []Song{}

	for i, v := range s.songs {
		if i == idx {
			continue
		}

        v.ID = i

		songs = append(songs, v)
	}

	s.songs = songs

	return s.Write()

}

func (s *Setlist) MoveUp(idx int) error {
	if idx <= 0 || idx >= len(s.songs) {
		return nil
	}

	song := s.songs[idx]
    song.ID -= 1

	s.songs[idx] = s.songs[idx-1]
    s.songs[idx].ID += 1

	s.songs[idx-1] = song

	return s.Write()
}

func (s *Setlist) MoveDown(idx int) error {
	if idx < 0 || idx >= len(s.songs)-1 {
		return nil
	}

	song := s.songs[idx]
    song.ID += 1

	s.songs[idx] = s.songs[idx+1]
    s.songs[idx].ID -= 1

	s.songs[idx+1] = song

	return s.Write()
}

func (s *Setlist) GetSongs() []Song {
	return s.songs
}
