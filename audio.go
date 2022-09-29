package chimer

import (
	"fmt"
	"os"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
)

func LoadSound(path string) (*beep.Buffer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open sound: %w", err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("failed to decode %s: %w", path, err)
	}

	defer func(streamer beep.StreamSeekCloser) {
		_ = streamer.Close()
	}(streamer)

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	return buffer, nil
}

type AudioCache = Cache[string, *beep.Buffer]

type SoundLoader struct{}

func (a SoundLoader) Get(path string) (*beep.Buffer, error) {
	return LoadSound(path)
}
