package main

import(
	"strings"
)

var(
	FILE_TO_ALGEB 	map[int]string 	=	map[int]string{0:"a",1:"b",2:"c",3:"d",4:"e",5:"f",6:"g",7:"h"}
	ALGEB_TO_FILE 	map[string]int 	=	map[string]int{"a":0,"b":1,"c":2,"d":3,"e":4,"f":5,"g":6,"h":7}
	RANK_TO_ALGEB 	map[int]string 	=	map[int]string{0:"8",1:"7",2:"6",3:"5",4:"4",5:"3",6:"2",7:"1"}
	ALGEB_TO_RANK 	map[string]int 	=	map[string]int{"8":0,"7":1,"6":2,"5":3,"4":4,"3":5,"2":6,"1":7}
	FILE_LETTER		[]string 		=	[]string{"a","b","c","d","e","f","g","h"}
	RANK_LETTER		[]string 		=	[]string{"1","2","3","4","5","6","7","8"}
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

func MoveFromAlgeb(algeb string) Move {
	m := Move{}

	t := NewTokenizer(algeb)

	fromalgeb := t.Pull(FILE_LETTER) + t.Pull(RANK_LETTER)
	toalgeb := t.Pull(FILE_LETTER) + t.Pull(RANK_LETTER)

	m.From = SquareFromAlgeb(fromalgeb)
	m.To = SquareFromAlgeb(toalgeb)

	m.Prom = Piece{Kind:"-"}

	if t.Has() {
		c := t.PullOne()
		m.Prom = Piece{Kind:c}
	}

	return m
}

func (m Move) Toalgeb() string {
	fromalgeb := m.From.Toalgeb()
	toalgeb := m.To.Toalgeb()

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

func SquareFromAlgeb(algeb string) Square {
	t := NewTokenizer(algeb)
	fls := t.Pull(FILE_LETTER)
	rls := t.Pull(RANK_LETTER)
	f , _ := ALGEB_TO_FILE[fls]
	r , _ := ALGEB_TO_RANK[rls]
	return SquareFromFileRank(f,r)
}

func (sq Square) Toalgeb() string {
	return FILE_TO_ALGEB[sq.File] + RANK_TO_ALGEB[sq.Rank]
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