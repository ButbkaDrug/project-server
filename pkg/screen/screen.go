package screen

import (
	"github.com/butbkadrug/reaper-lyrics/pkg/sse"
)


type ScreenUpdate []string

type Screen struct {
    Data []string
    Clients *sse.SSE
    Ch *chan []string
}

func NewScreen() *Screen {

    ch := make(chan []string)
    return &Screen{
        Data: []string{"some", "lines", "of", "text", "here"},
        Clients: sse.NewSSE(),
        Ch: &ch,
    }
}

func(s *Screen) Update(data []string) {

    s.Data = data

}
