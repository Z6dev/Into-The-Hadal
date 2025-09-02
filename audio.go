package main

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

func PlayAudio(data []byte, ctx *audio.Context) {
	pl, err := ctx.NewPlayer(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	pl.Play()
}
