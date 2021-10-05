package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// MouseStrokeSource is a StrokeSource implementation of mouse.
type MouseStrokeSource struct{}

func (m *MouseStrokeSource) Position() (int, int) {
	return ebiten.CursorPosition()
}

func (m *MouseStrokeSource) IsJustReleased() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}
