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
	UP direction = iota
	DOWN
	RIGHT
	LEFT
)

type Item struct {
	PosX int
	PosY int
}

type Tail struct {
	PosX int
	PosY int
	Tail *Tail
}

type Player struct {
	PosX int
	PosY int
	Dir  direction
	Tail *Tail
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
				i <- UP
			}
			if event.Key == termbox.KeyArrowDown {
				i <- DOWN
			}
			if event.Key == termbox.KeyArrowRight {
				i <- RIGHT
			}
			if event.Key == termbox.KeyArrowLeft {
				i <- LEFT
			}
			if event.Key == termbox.KeyEsc {
				os.Exit(0)
			}
		}
	}
}

func (p *Player) addTale() {
	if p.Tail == nil {
		p.Tail = &Tail{}
		return
	}
	var tail *Tail
	tail = p.Tail
	for {
		if tail.Tail == nil {
			tail.Tail = &Tail{}
			break
		}
		tail = tail.Tail
	}
}

func update(p *Player, items []Item) {
	var prevX, prevY, tx, ty int
	prevX = p.PosX
	prevY = p.PosY
	switch p.Dir {
	// 座標が左上原点なので、上下の操作がひっくり返る
	case UP:
		p.PosY--
	case DOWN:
		p.PosY++
	case RIGHT:
		p.PosX++
	case LEFT:
		p.PosX--
	}
	for index, item := range items {
		// itemとあたったら取得して、tailを伸ばす
		if item.PosX == p.PosX && item.PosY == p.PosY {
			items = append(items[:index], items[index+1:]...)
			p.addTale()
			break
		}
	}

	// tailのposition更新
	tail := p.Tail
	flag := false
	for {
		if tail == nil {
			break
		}
		// ここでリタイアチェック
		if p.PosX == tail.PosX && p.PosY == tail.PosY {
			flag = true
			break
		}
		tx, ty = tail.PosX, tail.PosY
		tail.PosX = prevX
		tail.PosY = prevY
		prevX, prevY = tx, ty
		tail = tail.Tail
	}
	if flag {
		termbox.Clear(coldef, coldef)
		for i, s := range "finish!!!" {
			termbox.SetCell(20+i, 10, rune(s), coldef, coldef)
		}
		// todo time表示
		termbox.Flush()
		os.Exit(0)
	}
}

func render(width, height int, p *Player, items []Item) {
	termbox.Clear(coldef, coldef)
	// 壁描画
	for y := 0; y < height; y++ {
		termbox.SetCell(0, y, rune('|'), coldef, coldef)
		termbox.SetCell(width, y, rune('|'), coldef, coldef)
	}
	// 床描画
	for x := 0; x < width+1; x++ {
		termbox.SetCell(x, height, rune('-'), coldef, coldef)
	}
	// userpotision反映
	for _, item := range items {
		termbox.SetCell(item.PosX, item.PosY, rune('X'), coldef, coldef)
	}
	termbox.SetCell(p.PosX, p.PosY, rune('O'), coldef, coldef)

	// TODO: update tail
	tail := p.Tail
	for {
		if tail == nil {
			break
		}
		termbox.SetCell(tail.PosX, tail.PosY, rune('o'), coldef, coldef)
		tail = tail.Tail
	}
	termbox.Flush()
}

func main() {
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()
	p := Player{PosX: 25, PosY: 25, Dir: DOWN}
	w, h := 60, 40
	// 適当にいっぱい生成
	items := []Item{}
	for i := 0; i < 100; i++ {
		items = append(items, Item{PosX: rand.Intn(w - 1), PosY: rand.Intn(h - 1)})
	}
	initGame(w, h, &p, items)

	// ユーザ入力取得用のgorutine
	input := make(chan direction)
	defer close(input)
	go takeInput(input)

	// 更新タイマー
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		// key イベント
		case d := <-input:
			p.Dir = d
		// タイマーイベント
		case <-ticker.C:
			update(&p, items)
			render(w, h, &p, items)
		}
	}
}
