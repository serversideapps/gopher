package main

import(
	"github.com/gopherjs/gopherjs/js"
	"strconv"
	"strings"
)

const(
	SQUARE_SIZE					float64		=	50.0
	SQUARE_PADDING				float64		=	4.0	
	PIECE_FONT_FACTOR			float64		=	1.5
	PIECE_TOP_FACTOR			float64		=	3.0	
	LIGHT_SQUARE_COLOR			string 		=	"#efefef"
	DARK_SQUARE_COLOR			string 		=	"#7f7f7f"
	SQUARE_Z_INDEX				string 		=	"100"
	PIECE_Z_INDEX				string 		=	"200"
	SQUARE_OPACITY				string 		=	"0.2"
	BOARD_MARGIN				float64		=	10.0
)

var(
	scalefactor		float64		=	0.8

	variant 		string 		=	DEFAULT_VARIANT

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

	PIECE_OPACITIES 		[4]string 	=	[4]string{"1.0","0.7","1.0","1.0"}
	PIECE_FILL_COLORS	 	[4]string 	=	[4]string{"#000000","#ffffff","#ffff00","#ff0000"}
	PIECE_STROKE_COLORS 	[4]string 	=	[4]string{"#ffffff","#afafaf","#afafaf","#afafaf"}
)

func TotalBoardSize(numranks int) float64 {
	return float64(numranks) * SQUARE_SIZE + 2.0 * BOARD_MARGIN
}

func Scalefactor() float64 {
	return scalefactor * TotalBoardSize(14) / TotalBoardSize(rb.Numranks)
}

func Root() *js.Object {
	return Dgebid("root")
}

type Combo struct {
	Selected	string
	Options		map[string]string
}

func CreateCombo(c Combo, handler func(*js.Object)) *js.Object {
	sel := Document().Call("createElement","select")
	for k , v := range c.Options {
		opt := Document().Call("createElement","option")
		opt.Set("id",k)
		opt.Set("name",k)
		opt.Set("value",k)
		opt.Set("innerHTML",v)
		if v == c.Selected {
			opt.Set("selected","true")
		}
		sel.Call("appendChild",opt)
	}
	sel.Call("addEventListener", "change", handler)
	return sel
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
	return ScreenVector{sv.X * Scalefactor(),sv.Y * Scalefactor()}
}

func (sv ScreenVector) Unscaled() ScreenVector {
	return ScreenVector{sv.X / Scalefactor(),sv.Y / Scalefactor()}
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
	return coord * Scalefactor()
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
	if dragunderway {
		dragunderway = false	

		dsq := rb.ScaledScreenVectorToSquare(dragd.Correct(HALF_SQUARE_SIZE_SCREENVECTOR.Scaled()))

		fromalgeb := IdPart(draggedid, 1)

		fromsqorig := rb.Rot(rb.SquareFromAlgeb(fromalgeb),flip)

		tosq := rb.Rot(fromsqorig.Plus(dsq),-flip)

		toalgeb := rb.SquareToAlgeb(tosq)

		algeb := fromalgeb + toalgeb

		m := rb.MoveFromAlgeb(algeb)

		rb.MakeMove(m)

		dsv := rb.SquareToScaledScreenVector(dsq)

		nsv := dragstartst.Plus(dsv)

		st := NewStyleFromId(draggedid)

		st.SetTopLeft(nsv)

		SetStyleOfId(draggedid,*st)

		DrawBoard()
	}
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

func VariantComboHandler(event *js.Object) {			
	target := event.Get("target")
	key := target.Get("selectedOptions").Get("0").Get("value").String()
	variant , _ = SUPPORTED_VARIANT_KEYS[key]
	rb = NewRawBoard(variant)
	rb.SetFromStartrawfen()
	DrawBoard()
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


func GrowButtonHandler(event *js.Object) {			
	scalefactor*=1.1
	DrawBoard()
}

func ShrinkButtonHandler(event *js.Object) {			
	scalefactor/=1.1
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
	containerdiv := CreateDiv("boardcontainer")	
	st := NewStyle("position:relative;background:url(assets/images/backgrounds/wood.jpg);")
	st.Set("width",Scaledpx(float64(rb.Numfiles)*SQUARE_SIZE+2.0*BOARD_MARGIN))
	st.Set("height",Scaledpx(float64(rb.Numranks)*SQUARE_SIZE+2.0*BOARD_MARGIN))
	containerdiv.Set("style",st.Report())
	div := CreateDiv("board")	
	div.Call("addEventListener", "mouseup", BoardMouseUpHandler)
	div.Call("addEventListener", "mousemove", BoardMouseMoveHandler)
	st = NewStyle("position:absolute;background:url(assets/images/backgrounds/wood.jpg);")
	st.Set("width",Scaledpx(float64(rb.Numfiles)*SQUARE_SIZE))
	st.Set("height",Scaledpx(float64(rb.Numranks)*SQUARE_SIZE))
	st.Set("top",Scaledpx(BOARD_MARGIN))
	st.Set("left",Scaledpx(BOARD_MARGIN))
	div.Set("style",st.Report())
	for nf := 0 ; nf < rb.Numfiles ; nf++ {
		for nr := 0 ; nr < rb.Numranks ; nr++ {
			bcol := LIGHT_SQUARE_COLOR			
			if ((nr+nf)%2)==1 {
				bcol = DARK_SQUARE_COLOR			
			}			
			sq := SquareFromFileRank(nf,nr)
			algeb := rb.SquareToAlgeb(sq)
			rotsq := rb.Rot(sq,flip)
			f := rotsq.File
			r := rotsq.Rank

			squarediv := CreateDiv("square_" + algeb)			
			style := NewStyle("position:absolute;")			
			style.Set("z-index",SQUARE_Z_INDEX)
			style.Set("width",Scaledpx(SQUARE_SIZE))
			style.Set("height",Scaledpx(SQUARE_SIZE))
			style.Set("background-color",bcol)
			style.Set("opacity",SQUARE_OPACITY)
			style.Set("top",Scaledpx(float64(r)*SQUARE_SIZE))
			style.Set("left",Scaledpx(float64(f)*SQUARE_SIZE))
			squarediv.Set("style",style.Report())

			if rb.IsSquareValid(sq) {
				div.Call("appendChild",squarediv)
			}

			p := rb.PieceAtFileRank(nf,nr)
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
			style.Set("opacity",PIECE_OPACITIES[p.Color])
			fillcol := PIECE_FILL_COLORS[p.Color]
			strokecol := PIECE_STROKE_COLORS[p.Color]
			svg := pieces[p.Kind]			
			svgc := strings.Replace(svg,"fill=\"#101010\"","fill=\""+fillcol+"\"",-1)
			svgc = strings.Replace(svgc,"fill:#ececec","fill:"+strokecol,-1)
			svgc = strings.Replace(svgc,"stroke:#101010","stroke:"+fillcol,-1)
			piecediv.Set("innerHTML",svgc)
			piecediv.Set("style",style.Report())
			if p.Kind!="-" {
				div.Call("appendChild",piecediv)
			}
		}
	}
	containerdiv.Call("appendChild",div)
	return containerdiv
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

	controlsdiv := CreateDiv("controls")

	variantcombo := CreateCombo(Combo{variant,SUPPORTED_VARIANT_KEYS},VariantComboHandler)

	flipbutton := CreateButton("Flip",FlipButtonHandler)

	resetbutton := CreateButton("Reset",ResetButtonHandler)

	growbutton := CreateButton("+",GrowButtonHandler)

	shrinkbutton := CreateButton("-",ShrinkButtonHandler)

	controlsdiv.Call("appendChild",variantcombo)
	controlsdiv.Call("appendChild",flipbutton)
	controlsdiv.Call("appendChild",resetbutton)
	controlsdiv.Call("appendChild",growbutton)
	controlsdiv.Call("appendChild",shrinkbutton)

	Root().Call("appendChild",board)

	Root().Call("appendChild",controlsdiv)
}

func main() {

	DocumentElement().Set("style","background-color:#afafaf;")

	rb = NewRawBoard(variant)
	rb.SetFromStartrawfen()

	DrawBoard()
}