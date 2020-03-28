package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

type direction int

const (
	up direction = iota
	down
	right
	left
)

type Item struct {
	PosX int
	PosY int
}

type Tale struct {
	PosX int
	PosY int
	Tale *Tale
}

type Player struct {
	PosX int
	PosY int
	Dir  direction
	Tale *Tale
}

const coldef = termbox.ColorDefault

func clearCheckLoop() {}

func initGame(w, h int, p *Player, items []Item) {
	termbox.HideCursor()
	render(w, h, p, items)
}

func takeInput(i chan<- direction) {
	for {
		switch event := termbox.PollEvent(); event.Type {
		case termbox.EventKey:
			if event.Key == termbox.KeyArrowUp {
				i <- up
			}
			if event.Key == termbox.KeyArrowDown {
				i <- down
			}
			if event.Key == termbox.KeyArrowRight {
				i <- right
			}
			if event.Key == termbox.KeyArrowLeft {
				i <- left
			}
			if event.Key == termbox.KeyEsc {
				os.Exit(0)
			}
		}
	}
}

func update(p *Player, items []Item) {
	for index, item := range items {
		if item.PosX == p.PosX && item.PosY == p.PosY {
			items = append(items[:index], items[index+1:]...)
			break
		}
	}
}

func render(width, height int, p *Player, items []Item) {
	termbox.Clear(coldef, coldef)
	for x := 0; x < 50; x++ {
		for y := 0; y < height; y++ {
			termbox.SetCell(x, y, 'a', coldef, coldef)
		}
	}
	// userpotision反映
	for _, item := range items {
		termbox.SetCell(item.PosX, item.PosY, '★', coldef, coldef)
	}
	termbox.SetCell(p.PosX, p.PosY, '●', coldef, coldef)
	tale := p.Tale
	for {
		if tale == nil {
			break
		}
	}
	termbox.Flush()
}

func main() {
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()
	p := Player{PosX: 0, PosY: 10, Dir: down}
	w, h := termbox.Size()
	items := []Item{}
	for i := 0; i < 30; i++ {
		items = append(items, Item{PosX: rand.Intn(w - 1), PosY: rand.Intn(h - 1)})
	}
	initGame(w, h, &p, items)

	// ユーザ入力取得用のgorutine
	input := make(chan direction)
	defer close(input)
	go takeInput(input)

	ticker := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case d := <-input:
			p.Dir = d
		case <-ticker.C:
			switch p.Dir {
			// 座標が左上原点なので、上下の操作がひっくり返る
			case up:
				p.PosY--
			case down:
				p.PosY++
			case right:
				p.PosX++
			case left:
				p.PosX--
			}
			update(&p, items)
			render(w, h, &p, items)
		}
	}
}
