package main

var(
	FILE_TO_ALGEB 	map[int]string 	=	map[int]string{0:"a",1:"b",2:"c",3:"d",4:"e",5:"f",6:"g",7:"h"}
	RANK_TO_ALGEB 	map[int]string 	=	map[int]string{0:"8",1:"7",2:"6",3:"5",4:"4",5:"3",6:"2",7:"1"}
)

type 	Square 		struct {
	File	int
	Rank	int
}

func SquareFromFileRank(f int, r int) Square {
	return Square{File:f, Rank:r}
}

func (sq Square) Toalgeb() string {
	return FILE_TO_ALGEB[sq.File] + RANK_TO_ALGEB[sq.Rank]
}