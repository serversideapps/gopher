package main

import(
	"strconv"
)

const(
	DEFAULT_VARIANT		string	=	"Four Player"	
)

var(
	SUPPORTED_VARIANT_KEYS		map[string]string 			= map[string]string{
		"standard"		:	"Standard",
		"fourplayer"	:	"Four Player",
	}	
	VARIANT_PROPERTIES		map[string]VariantProperty	= map[string]VariantProperty{
			"Standard"		:	VariantProperty{
			Numfiles		:	8,
			Numranks		:	8,
			Startrawfen		:	"r0n0b0q0k0b0n0r0"+
								"p0p0p0p0p0p0p0p0"+
								"-0-0-0-0-0-0-0-0"+
								"-0-0-0-0-0-0-0-0"+
								"-0-0-0-0-0-0-0-0"+
								"-0-0-0-0-0-0-0-0"+
								"p1p1p1p1p1p1p1p1"+
								"r1n1b1q1k1b1n1r1",
		},
		"Four Player"		:	VariantProperty{
			Numfiles		:	14,
			Numranks		:	14,
			Startrawfen		:	"-0-0-0r0n0b0q0k0b0n0r0-0-0-0"+
								"-0-0-0p0p0p0p0p0p0p0p0-0-0-0"+
								"-0-0-0-0-0-0-0-0-0-0-0-0-0-0"+
								"r3p3-0-0-0-0-0-0-0-0-0-0p2r2"+
								"n3p3-0-0-0-0-0-0-0-0-0-0p2n2"+
								"b3p3-0-0-0-0-0-0-0-0-0-0p2b2"+
								"k3p3-0-0-0-0-0-0-0-0-0-0p2q2"+
								"q3p3-0-0-0-0-0-0-0-0-0-0p2k2"+
								"b3p3-0-0-0-0-0-0-0-0-0-0p2b2"+
								"n3p3-0-0-0-0-0-0-0-0-0-0p2n2"+
								"r3p3-0-0-0-0-0-0-0-0-0-0p2r2"+
								"-0-0-0-0-0-0-0-0-0-0-0-0-0-0"+
								"-0-0-0p1p1p1p1p1p1p1p1-0-0-0"+
								"-0-0-0r1n1b1k1q1b1n1r1-0-0-0",
		},
	}
)

func StringAt(str string, index int) string {
	return string([]rune(str)[index])
}

type Piece			struct {
	Kind			string
	Color			int
}

type RawBoard		struct {
	Variant			string
	Numfiles		int
	Lastfile		int
	Numranks		int
	Lastrank		int
	Area			int
	Rep				[]Piece
}

type VariantProperty	struct {
	Numfiles		int
	Numranks		int	
	Startrawfen 	string
}

func NewRawBoard(setvariant ...string) *RawBoard{
	rb := RawBoard{}
	rb.Variant = DEFAULT_VARIANT
	if len(setvariant) > 0 {
		rb.Variant = setvariant[0]
	}
	vp := VARIANT_PROPERTIES[rb.Variant]
	rb.Numfiles = vp.Numfiles
	rb.Lastfile = vp.Numfiles - 1
	rb.Numranks = vp.Numranks
	rb.Lastrank = vp.Numranks - 1
	rb.Area = rb.Numfiles * rb.Numranks
	rb.Rep = make ( []Piece , rb.Area )
	return &rb
}

func (rb RawBoard) IndexOfFileRank(f int, r int) int {
	return f + r * rb.Numfiles
}

func (rb RawBoard) PieceAtFileRank(f int, r int) Piece {
	return rb.Rep[rb.IndexOfFileRank(f,r)]
}

func (rb *RawBoard) SetFromRawFen(rf string) {
	for i := 0 ; i < rb.Area ; i++ {
		index := 2 * i
		kind := StringAt(rf,index)
		color , _ := strconv.Atoi(StringAt(rf,index+1))
		p := Piece{Kind:kind,Color:color}
		rb.Rep[i] = p
	}
}

func (rb *RawBoard) SetFromStartrawfen() {
	startrawfen := VARIANT_PROPERTIES[rb.Variant].Startrawfen
	rb.SetFromRawFen(startrawfen)
}

func (rb RawBoard) IS_STANDARD() bool {
	return rb.Variant == "Standard"
}

func (rb RawBoard) IS_FOUR_PLAYER() bool {
	return rb.Variant == "Four Player"
}