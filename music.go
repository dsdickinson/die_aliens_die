package main

import (
	"os"
	"time"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func loadBackgroundMusic(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
//	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
//	speaker.Play(streamer)
//	select {}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()
	music := buffer.Streamer(0, buffer.Len())
	speaker.Play(music)

	return nil
}

