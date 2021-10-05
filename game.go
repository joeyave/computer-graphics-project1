package main

import (
	"embed"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"image/png"
	"log"
	"math/rand"
)

type Game struct {
	touchIDs []ebiten.TouchID
	strokes  map[*Stroke]struct{}
	sprites  []*Sprite
}

//go:embed images/*
var imagesFS embed.FS

var images []*image.Image

func init() {
	imagesDir, err := imagesFS.ReadDir("images")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range imagesDir {
		if file.IsDir() {
			continue
		}

		img, err := imagesFS.Open("images/" + file.Name())
		if err != nil {
			continue
		}

		pngImg, err := png.Decode(img)
		if err != nil {
			continue
		}
		images = append(images, &pngImg)
	}
}

func NewGame() *Game {
	// Initialize the sprites.
	var sprites []*Sprite

	for _, img := range images {
		ebitenImage := ebiten.NewImageFromImage(*img)
		w, h := ebitenImage.Size()
		for i := 0; i < 3; i++ {
			s := &Sprite{
				image: ebitenImage,
				x:     rand.Intn(screenWidth - w),
				y:     rand.Intn(screenHeight - h),
			}
			sprites = append(sprites, s)
		}
	}

	// Initialize the game.
	return &Game{
		strokes: map[*Stroke]struct{}{},
		sprites: sprites,
	}
}

func (g *Game) spriteAt(x, y int) *Sprite {
	// As the sprites are ordered from back to front,
	// search the clicked/touched sprite in reverse order.
	for i := len(g.sprites) - 1; i >= 0; i-- {
		s := g.sprites[i]
		if s.In(x, y) {
			return s
		}
	}
	return nil
}

func (g *Game) updateStroke(stroke *Stroke) {
	stroke.Update()
	if !stroke.IsReleased() {
		return
	}

	s := stroke.DraggingObject().(*Sprite)
	if s == nil {
		return
	}

	s.MoveBy(stroke.PositionDiff())

	index := -1
	for i, ss := range g.sprites {
		if ss == s {
			index = i
			break
		}
	}

	// Move the dragged sprite to the front.
	g.sprites = append(g.sprites[:index], g.sprites[index+1:]...)
	g.sprites = append(g.sprites, s)

	stroke.SetDraggingObject(nil)
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s := NewStroke(&MouseStrokeSource{})
		s.SetDraggingObject(g.spriteAt(s.Position()))
		g.strokes[s] = struct{}{}
	}
	for _, id := range g.touchIDs {
		s := NewStroke(&TouchStrokeSource{id})
		s.SetDraggingObject(g.spriteAt(s.Position()))
		g.strokes[s] = struct{}{}
	}

	for s := range g.strokes {
		g.updateStroke(s)
		if s.IsReleased() {
			delete(g.strokes, s)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	draggingSprites := map[*Sprite]struct{}{}
	for s := range g.strokes {
		if sprite := s.DraggingObject().(*Sprite); sprite != nil {
			draggingSprites[sprite] = struct{}{}
		}
	}

	for _, s := range g.sprites {
		if _, ok := draggingSprites[s]; ok {
			continue
		}
		s.Draw(screen, 0, 0, 1)
	}
	for s := range g.strokes {
		dx, dy := s.PositionDiff()
		if sprite := s.DraggingObject().(*Sprite); sprite != nil {
			sprite.Draw(screen, dx, dy, 0.5)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
