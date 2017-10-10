package main

import(
	"github.com/gopherjs/gopherjs/js"
	"strconv"
)

func Document() *js.Object {
	return js.Global.Get("document")
}

func DocumentElement() *js.Object {
	return Document().Get("documentElement")
}

func CreateElement(kind string) *js.Object {
	return Document().Call("createElement",kind)
}

func CreateDiv() *js.Object {
	return CreateElement("div")
}

func (rb RawBoard) RawBoardDiv() *js.Object {
	div := CreateDiv()
	div.Set("style","position:absolute;top:20px;left:20px;")
	for f := 0 ; f < rb.Numfiles ; f++ {
		for r := 0 ; r < rb.Numranks ; r++ {
			squarediv := CreateDiv()			
			style := "width:50px;height:50px;position:absolute;z-index:100;"
			bcol := "#afafaf"
			if ((r+f)%2)==0 {
				bcol = "#dfdfdf"
			}			
			style += "background-color:"+bcol+";"
			style += "top: " + strconv.Itoa(r*50) + "px;"
			style += "left: " + strconv.Itoa(f*50) + "px;"			
			squarediv.Set("style",style)
			div.Call("appendChild",squarediv)
			piecediv := CreateDiv()
			style = "font-size:20px;width:30px;height:30px;position:absolute;z-index:200;"
			style += "top: " + strconv.Itoa(r*50+10) + "px;"
			style += "left: " + strconv.Itoa(f*50+10) + "px;"
			bcol = "#00ff00"
			p := rb.PieceAtFileRank(f,r)
			if p.Color==0 {
				bcol = "#ff0000"
			}
			style += "background-color:"+bcol+";"
			piecediv.Set("style",style)			
			pieceletterdiv := CreateDiv()
			pieceletterdiv.Set("style","position:absolute;left:10px;")
			pieceletterdiv.Set("innerHTML",p.Kind)
			piecediv.Call("appendChild",pieceletterdiv)
			if p.Kind!="-" {
				div.Call("appendChild",piecediv)
			}
		}
	}
	return div
}

func main() {

	rb := NewRawBoard()
	rb.SetFromStartrawfen()
	
	div := rb.RawBoardDiv()
	DocumentElement().Call("appendChild",div)

}