package main

import(
	"github.com/gopherjs/gopherjs/js"
	"strconv"
	"strings"
)

const(
	SQUARE_SIZE			float64		=	50.0
	SQUARE_PADDING		float64		=	10.0	
	PIECE_FONT_FACTOR	float64		=	1.5
	PIECE_TOP_FACTOR	float64		=	3.0
	WHITE_PIECE_COLOR 	string 		=	"#00ff00"
	BLACK_PIECE_COLOR	string 		=	"#ff0000"
	LIGHT_SQUARE_COLOR	string 		=	"#afafaf"
	DARK_SQUARE_COLOR	string 		=	"#dfdfdf"
	SQUARE_Z_INDEX		string 		=	"100"
	PIECE_Z_INDEX		string 		=	"200"
)

var(
	scalefactor		float64		=	1.0

	flip 			int 		=	0

	rb 				*RawBoard

	HALF_SQUARE_SIZE 	float64		=	SQUARE_SIZE / 2.0

	PIECE_SIZE			float64		=	SQUARE_SIZE - 2.0 * SQUARE_PADDING

	dragstart		ScreenVector
	dragstartst		ScreenVector
	dragd			ScreenVector
	
	dragunderway	bool
	draggedid		string	

	HALF_SQUARE_SIZE_SCREENVECTOR	ScreenVector 	=	ScreenVector{HALF_SQUARE_SIZE,HALF_SQUARE_SIZE}
)

func Root() *js.Object {
	return Dgebid("root")
}

func IdPart(id string, i int) string {
	parts := strings.Split(id,"_")
	return parts[i]
}

type ScreenVector	struct {
	X 	float64
	Y 	float64
}

func (sv1 ScreenVector) Plus(sv2 ScreenVector) ScreenVector {
	return ScreenVector{sv1.X + sv2.X,sv1.Y + sv2.Y}
}

func (sv1 ScreenVector) Minus(sv2 ScreenVector) ScreenVector {
	return ScreenVector{sv1.X - sv2.X,sv1.Y - sv2.Y}
}

func (sv1 ScreenVector) Correct(sv2 ScreenVector) ScreenVector {
	x := sv1.X + sv2.X
	y := sv1.Y + sv2.Y
	if sv1.X < 0 {
		x = sv1.X - sv2.X
	}
	if sv1.Y < 0 {
		y = sv1.Y - sv2.Y
	}
	return ScreenVector{x,y}
}

func (sv ScreenVector) Scaled() ScreenVector {
	return ScreenVector{sv.X * scalefactor,sv.Y * scalefactor}
}

func (sv ScreenVector) Unscaled() ScreenVector {
	return ScreenVector{sv.X / scalefactor,sv.Y / scalefactor}
}

type Style 		struct {
	Init 			string
	Properties		map[string]string
}

func (st Style) GetProperty(property string) string {
	return st.Properties[property]
}

func (st Style) GetFloat(property string) float64 {
	f , _ := strconv.ParseFloat(st.GetProperty(property),64)
	return f
}

func (st Style) GetPx(property string) float64 {
	f , _ := strconv.ParseFloat(strings.Replace(st.GetProperty(property),"px","",-1),64)
	return f
}

func NewStyle(init ...string) *Style {	
	st := Style{}
	properties := make( map[string]string )
	if len(init) > 0 {
		parts := strings.Split(init[0],";")
		for _ , part := range parts {
			partnospace := strings.Replace(part," ","",-1)
			subparts := strings.Split(partnospace,":")
			if len(subparts) == 2 {
				property := subparts[0]
				value := subparts[1]
				properties[property] = value
			}
		}
	}		
	st.Properties = properties
	return &st
}

func NewStyleFromId(id string) *Style {
	e := Dgebid(id)
	style := e.Get("style").Get("cssText").String()
	return NewStyle(style)
}

func SetStyleOfId(id string, st Style) {
	Dgebid(id).Set("style",st.Report())
}

func (st *Style) Set(property string, value string) {
	st.Properties[property] = value
}

func (st *Style) SetPx(property string, value float64) {
	st.Properties[property] = strconv.FormatFloat(value, 'g', -1, 64) + "px"
}

func (st *Style) SetTopLeft(sv ScreenVector) {
	st.SetPx("top",sv.Y)
	st.SetPx("left",sv.X)
}

func (st Style) Report() string {
	ps := make( []string , 0 )
	ps = append(ps, st.Init)
	for pr , v := range st.Properties {
		ps = append(ps, pr + ": " + v + ";")
	}
	return strings.Join(ps," ")
}

func Document() *js.Object {
	return js.Global.Get("document")
}

func DocumentElement() *js.Object {
	return Document().Get("documentElement")
}

func CreateElement(kind string) *js.Object {
	return Document().Call("createElement",kind)
}

func CreateDiv(id ...string) *js.Object {
	e := CreateElement("div")
	if len(id) > 0 {
		e.Set("id",id[0])
	}
	return e
}

func Dgebid(id string) *js.Object {
	return Document().Call("getElementById",id)
}

func Scaled(coord float64) float64 {
	return coord * scalefactor
}

func Px(coord float64) string {
	return strconv.Itoa(int(coord))+"px;"
}

func Scaledpx(coord float64) string {
	return Px(Scaled(coord))
}

func PieceDragStartHandler(event *js.Object) {	
	event.Call("preventDefault")
	target := event.Get("target")		
	dragstart = ScreenVector{event.Get("clientX").Float(),event.Get("clientY").Float()}
	draggedid = target.Get("id").String()		
	st := NewStyleFromId(draggedid)
	dragstartst = ScreenVector{st.GetPx("left"),st.GetPx("top")}
	dragunderway = true
}

func BoardMouseUpHandler(event *js.Object) {			
	dragunderway = false	

	dsq := rb.ScaledScreenVectorToSquare(dragd.Correct(HALF_SQUARE_SIZE_SCREENVECTOR.Scaled()))

	fromalgeb := IdPart(draggedid, 1)

	fromsqorig := rb.Rot(SquareFromAlgeb(fromalgeb),flip)

	tosq := rb.Rot(fromsqorig.Plus(dsq),-flip)

	toalgeb := tosq.Toalgeb()

	algeb := fromalgeb + toalgeb

	m := MoveFromAlgeb(algeb)

	rb.MakeMove(m)

	println(m.Toalgeb())

	dsv := rb.SquareToScaledScreenVector(dsq)

	nsv := dragstartst.Plus(dsv)

	st := NewStyleFromId(draggedid)

	st.SetTopLeft(nsv)

	SetStyleOfId(draggedid,*st)

	DrawBoard()
}

func BoardMouseMoveHandler(event *js.Object) {		
	if dragunderway {
		client := ScreenVector{event.Get("clientX").Float(),event.Get("clientY").Float()}		

		dragd = client.Minus(dragstart)

		st := NewStyleFromId(draggedid)		

		nsv := dragstartst.Plus(dragd)

		st.SetTopLeft(nsv)

		SetStyleOfId(draggedid,*st)
	}
}

func FlipButtonHandler(event *js.Object) {			
	flip += 1
	if flip >3 {
		flip = 0
	}	
	DrawBoard()
}

func ResetButtonHandler(event *js.Object) {			
	rb.SetFromStartrawfen()
	DrawBoard()
}

func (rb RawBoard) ScreenVectorToSquare(sv ScreenVector) Square {
	f := int(sv.X / SQUARE_SIZE)
	r := int(sv.Y / SQUARE_SIZE)
	return SquareFromFileRank(f,r)
}

func (rb RawBoard) ScaledScreenVectorToSquare(sv ScreenVector) Square {
	return rb.ScreenVectorToSquare(sv.Unscaled())
}

func (rb RawBoard) SquareToScreenVector(sq Square) ScreenVector {
	x := float64(sq.File) * SQUARE_SIZE
	y := float64(sq.Rank) * SQUARE_SIZE
	return ScreenVector{x,y}
}

func (rb RawBoard) SquareToScaledScreenVector(sq Square) ScreenVector {	
	return rb.SquareToScreenVector(sq).Scaled()
}

func (rb RawBoard) Js() *js.Object {
	div := CreateDiv("board")	
	div.Call("addEventListener", "mouseup", BoardMouseUpHandler)
	div.Call("addEventListener", "mousemove", BoardMouseMoveHandler)
	st := NewStyle("position:relative; background-color: #00ff00;")
	st.Set("width",Scaledpx(float64(rb.Numfiles)*SQUARE_SIZE))
	st.Set("height",Scaledpx(float64(rb.Numranks)*SQUARE_SIZE))
	div.Set("style",st.Report())
	for nf := 0 ; nf < rb.Numfiles ; nf++ {
		for nr := 0 ; nr < rb.Numranks ; nr++ {
			bcol := LIGHT_SQUARE_COLOR			
			if ((nr+nf)%2)==0 {
				bcol = DARK_SQUARE_COLOR			
			}			
			sq := SquareFromFileRank(nf,nr)
			algeb := sq.Toalgeb()
			rotsq := rb.Rot(sq,flip)
			f := rotsq.File
			r := rotsq.Rank
			squarediv := CreateDiv("square_" + algeb)			
			style := NewStyle("position:absolute;")			
			style.Set("z-index",SQUARE_Z_INDEX)
			style.Set("width",Scaledpx(SQUARE_SIZE))
			style.Set("height",Scaledpx(SQUARE_SIZE))
			style.Set("background-color",bcol)
			style.Set("top",Scaledpx(float64(r)*SQUARE_SIZE))
			style.Set("left",Scaledpx(float64(f)*SQUARE_SIZE))
			squarediv.Set("style",style.Report())
			div.Call("appendChild",squarediv)
			piecediv := CreateDiv("piece_" + algeb)
			piecediv.Set("draggable","true")
			piecediv.Call("addEventListener","dragstart",PieceDragStartHandler)			
			style = NewStyle("position:absolute;")			
			style.Set("z-index",PIECE_Z_INDEX)
			style.Set("font-size",Scaledpx(PIECE_SIZE / PIECE_FONT_FACTOR))
			style.Set("width",Scaledpx(PIECE_SIZE))
			style.Set("height",Scaledpx(PIECE_SIZE))
			style.Set("top",Scaledpx(float64(r)*SQUARE_SIZE+SQUARE_PADDING))
			style.Set("left",Scaledpx(float64(f)*SQUARE_SIZE+SQUARE_PADDING))
			bcol = WHITE_PIECE_COLOR
			p := rb.PieceAtFileRank(nf,nr)
			if p.Color==0 {
				bcol = BLACK_PIECE_COLOR
			}
			style.Set("background-color",bcol)
			piecediv.Set("style",style.Report())			
			pieceletterdiv := CreateDiv()
			style = NewStyle("position:absolute;")
			style.Set("left",Scaledpx(SQUARE_PADDING))
			style.Set("top",Scaledpx(SQUARE_PADDING / PIECE_TOP_FACTOR))
			pieceletterdiv.Set("style",style.Report())
			pieceletterdiv.Set("innerHTML",p.Kind)
			piecediv.Call("appendChild",pieceletterdiv)
			if p.Kind!="-" {
				div.Call("appendChild",piecediv)
			}
		}
	}
	return div
}

func CreateButton(caption string,handler func(*js.Object)) *js.Object {
	button := Document().Call("createElement","input")

	button.Set("type","button")
	button.Set("value",caption)

	button.Call("addEventListener","mousedown",handler)			

	return button
}

func DrawBoard() {	
	Root().Set("innerHTML","")

	board := rb.Js()

	Root().Set("style","padding:30px;")
	
	Root().Call("appendChild",board)

	controlsdiv := CreateDiv("controls")

	flipbutton := CreateButton("Flip",FlipButtonHandler)

	resetbutton := CreateButton("Reset",ResetButtonHandler)

	controlsdiv.Call("appendChild",flipbutton)
	controlsdiv.Call("appendChild",resetbutton)

	Root().Call("appendChild",controlsdiv)
}

func main() {

	rb = NewRawBoard()
	rb.SetFromStartrawfen()

	DrawBoard()
}