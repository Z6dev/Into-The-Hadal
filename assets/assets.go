package assets

import (
	_ "embed"
)

//go:embed images/sandtile.png
var Sandtile_png []byte

//go:embed audio/brick-on-metal.mp3
var MetalThud_mp3 []byte

//go:embed audio/thud.mp3
var MetalCrash_mp3 []byte
