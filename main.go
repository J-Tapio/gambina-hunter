package main

import (
	"bytes"
	_ "image/png"
	"os"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var img *ebiten.Image
var img2 *ebiten.Image
var gambinaImg *ebiten.Image
var cut1 []byte
var cut2 []byte

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("chr.png")
	if err != nil {
		log.Fatal(err)
	}

	img2, _, err = ebitenutil.NewImageFromFile("chr2.png")
	if err != nil {
		log.Fatal(err)
	}

	gambinaImg, _, err = ebitenutil.NewImageFromFile("gambina-xs.png")
	if err != nil {
		log.Fatal(err)
	}

	cut1, err = os.ReadFile("./cut1.mp3")
	if err != nil {
		log.Fatal(err)
	}
	cut2, err = os.ReadFile("./cut2.mp3")
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	charX      int
	charY      int
	gambinaX int
	gambinaY int
	keys []ebiten.Key
	drinked int
	audioContext *audio.Context
	audioPlayer1 *audio.Player
	audioPlayer2 *audio.Player
	cut1audio *mp3.Stream
	cut2audio *mp3.Stream
}

func (g *Game) Random() {
	newX := rand.Intn(528 - 51)
	newY := rand.Intn(356 - 77)

	g.gambinaX = newX
	g.gambinaY = newY
}

func (g *Game) Update() error {
	if g.audioPlayer1 == nil || g.audioPlayer2 == nil {
		var err error
		g.cut1audio, err = mp3.DecodeWithSampleRate(48000, bytes.NewReader(cut1))
		if err != nil {
			log.Fatal(err)
		}
		g.cut2audio, err = mp3.DecodeWithSampleRate(48000, bytes.NewReader(cut2))
			if err != nil {
			log.Fatal(err)
		}

		g.audioPlayer1, err = g.audioContext.NewPlayer(g.cut1audio)
		if err != nil {
			log.Fatal(err)
		}
		g.audioPlayer1.SetVolume(0.4)

		g.audioPlayer2, err = g.audioContext.NewPlayer(g.cut2audio)
		if err != nil {
			log.Fatal(err)
		}
		g.audioPlayer2.SetVolume(0.4)
	}
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.charX), float64(g.charY))
	screen.DrawImage(img, op)

	gambinaOp := &ebiten.DrawImageOptions{}
	gambinaOp.GeoM.Translate(float64(g.gambinaX), float64(g.gambinaY))
	screen.DrawImage(gambinaImg, gambinaOp)

	for _, key := range g.keys {
		switch key {
		case 31: // UP
			{
				if g.charY > 0 {
					g.charY = g.charY - 3
				}
			}
		case 28: // DOWN
			{
				if g.charY < 356 {
					g.charY = g.charY + 3
				}
			}
		case 29: // LEFT
			{
				if g.charX > -8 {
					g.charX = g.charX - 3
				}
			}
		case 30: // RIGHT
			{
				if g.charX < 528 {
					g.charX = g.charX + 3
				}
				screen.DrawImage(img2, op)
			}
		}
	
		if g.charX > g.gambinaX {
			if (g.charX - g.gambinaX <=66) && (g.charY - g.gambinaY >= -40) && (g.charY - g.gambinaY < 40) {
				g.drinked = g.drinked +1
				g.Random()
				gambinaOp.GeoM.Translate(float64(g.gambinaX), float64(g.gambinaY))
				if err := g.audioPlayer1.Rewind(); err != nil {
					return
				}
				g.audioPlayer1.Play()
			}
		} else {
			if (g.gambinaX - g.charX <= 52) && (g.charY - g.gambinaY >= -40) && (g.charY - g.gambinaY < 40) {
				g.drinked = g.drinked +1
				g.Random()
				gambinaOp.GeoM.Translate(float64(g.gambinaX), float64(g.gambinaY))
				if err := g.audioPlayer2.Rewind(); err != nil {
					return
				}
				g.audioPlayer2.Play()
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Gambina Hunter")
	if err := ebiten.RunGame(&Game{
		charX: 320,
		charY: 320,
		gambinaX: 200,
		gambinaY: 200,
		drinked: 0,
		audioContext: audio.NewContext(48000),
	}); err != nil {
		log.Fatal(err)
	}
}
