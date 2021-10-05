package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TouchStrokeSource is a StrokeSource implementation of touch.
type TouchStrokeSource struct {
	ID ebiten.TouchID
}

func (t *TouchStrokeSource) Position() (int, int) {
	return ebiten.TouchPosition(t.ID)
}

func (t *TouchStrokeSource) IsJustReleased() bool {
	return inpututil.IsTouchJustReleased(t.ID)
}
