package main

// StrokeSource represents a input device to provide strokes.
type StrokeSource interface {
	Position() (int, int)
	IsJustReleased() bool
}
