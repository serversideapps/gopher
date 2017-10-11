package main

import(
	"strings"
	"strconv"
)

var(
	FILE_TO_ALGEB 	map[int]string 	=	map[int]string{0:"a",1:"b",2:"c",3:"d",4:"e",5:"f",6:"g",7:"h",8:"i",9:"j",10:"k",11:"l",12:"m",13:"n"}
	ALGEB_TO_FILE 	map[string]int 	=	map[string]int{"a":0,"b":1,"c":2,"d":3,"e":4,"f":5,"g":6,"h":7,"i":8,"j":9,"k":10,"l":11,"m":12,"n":13}
	
	FILE_LETTER		[]string 		=	[]string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n"}
	RANK_LETTER		[]string 		=	[]string{"0","1","2","3","4","5","6","7","8","9"}
)

type 	Square 		struct {
	File	int
	Rank	int
}

type 	Move		struct {
	From 	Square
	To 		Square
	Prom	Piece
}

func (t Tokenizer) Has() bool {
	return len(t.Parts) > 0
}

func (rb RawBoard) MoveFromAlgeb(algeb string) Move {
	m := Move{}

	t := NewTokenizer(algeb)

	fromalgeb := t.Pull(FILE_LETTER) + t.Pull(RANK_LETTER)
	toalgeb := t.Pull(FILE_LETTER) + t.Pull(RANK_LETTER)

	m.From = rb.SquareFromAlgeb(fromalgeb)
	m.To = rb.SquareFromAlgeb(toalgeb)

	m.Prom = Piece{Kind:"-"}

	if t.Has() {
		c := t.PullOne()
		m.Prom = Piece{Kind:c}
	}

	return m
}

func (rb RawBoard) MoveToAlgeb(m Move) string {
	fromalgeb := rb.SquareToAlgeb(m.From)
	toalgeb := rb.SquareToAlgeb(m.To)

	prom := ""

	if m.Prom.Kind != "-" {
		prom = m.Prom.Kind
	}

	return fromalgeb + toalgeb + prom
}

func (rb RawBoard) Rot(sq Square,flip int) Square {
	if flip < 0 {
		flip = 4 + flip
	}
	if flip == 0 {
		return sq
	}
	if flip == 1 {
		return Square{rb.Lastrank - sq.Rank,sq.File}
	}
	if flip == 2 {
		return Square{rb.Lastfile - sq.File,rb.Lastrank - sq.Rank}
	}
	return Square{sq.Rank,rb.Lastfile - sq.File}
}

type 	Tokenizer	struct {
	Parts	[]string
}

func NewTokenizer(str string) *Tokenizer {
	t := Tokenizer{}
	t.Parts = strings.Split(str,"")
	return &t
}

func (t *Tokenizer) PullOne() string {
	if !t.Has() {
		return ""
	}
	defer func(){
		t.Parts = t.Parts[1:]
	}()
	return t.Parts[0]
}

func (t *Tokenizer) Pull(accept []string) string {
	h := make( map[string]bool )
	for _ , a := range accept {
		h[a] = true
	}
	pull := ""
	for ; len(t.Parts) > 0 ; {
		c := t.Parts[0]
		_ , has := h[c]
		if !has {
			return pull
		}
		t.Parts = t.Parts[1:]
		pull += c
	}
	return pull
}

func SquareFromFileRank(f int, r int) Square {
	return Square{File:f, Rank:r}
}

func (rb RawBoard) SquareFromAlgeb(algeb string) Square {
	t := NewTokenizer(algeb)
	fls := t.Pull(FILE_LETTER)
	rls := t.Pull(RANK_LETTER)
	f , _ := ALGEB_TO_FILE[fls]
	rn , _ := strconv.Atoi(rls)
	r := rb.Numranks - rn
	return SquareFromFileRank(f,r)
}

func (rb RawBoard) SquareToAlgeb(sq Square) string {	
	falgeb , _ := FILE_TO_ALGEB[sq.File]
	ralgeb := strconv.Itoa(rb.Numranks - sq.Rank)
	return falgeb + ralgeb
}

func (sq1 Square) Plus(sq2 Square) Square {
	return Square{sq1.File + sq2.File,sq1.Rank + sq2.Rank}
}

func (rb RawBoard) IndexOfSquare(sq Square) int {
	return rb.IndexOfFileRank(sq.File, sq.Rank)
}

func (rb *RawBoard) Put(sq Square, p Piece) {
	index := rb.IndexOfSquare(sq)
	rb.Rep[index]=p
}

func (rb RawBoard) Get(sq Square) Piece {
	index := rb.IndexOfSquare(sq)
	return rb.Rep[index]
}

func (rb *RawBoard) MakeMove(m Move) {
	fromsq := m.From
	tosq := m.To
	
	frompiece := rb.Get(fromsq)

	//topiece := rb.Get(tosq)

	rb.Put(fromsq,Piece{Kind:"-"})

	rb.Put(tosq,frompiece)
}

func (rb RawBoard) IsSquareValid(sq Square) bool {
	if rb.IS_STANDARD() {
		return true
	}
	if rb.IS_FOUR_PLAYER() {
		if (sq.Rank < 3) && (sq.File < 3) {
			return false
		}
		if (sq.Rank < 3) && (sq.File > 10) {
			return false
		}
		if (sq.Rank > 10) && (sq.File > 10) {
			return false
		}
		if (sq.Rank > 10) && (sq.File < 3) {
			return false
		}
		return true
	}
	return false
}