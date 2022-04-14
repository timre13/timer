package label

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/gfx"
    "github.com/veandco/go-sdl2/ttf"
    "timer/common"
    "timer/iwidget"
)

type Label struct {
    XPos, YPos      int32           // Position of the top left corner
    BgColor         *sdl.Color      // Background color
    FgColor         *sdl.Color      // Foreground color
    Font            *ttf.Font       // The font that is used to draw the text
    Text            string
}
var _ iwidget.IWidget = (*Label)(nil)

func (e *Label) Draw(rend *sdl.Renderer) {
    if e.BgColor != nil {
        width, _, err := e.Font.SizeUTF8(e.Text)
        common.PANIC_ERR(err)
        gfx.BoxColor(rend, e.XPos, e.YPos, e.XPos+int32(width)+4, e.YPos+int32(e.Font.Height())+4, *e.BgColor)
    }

    common.RenderText(rend, e.Font, e.Text, e.FgColor, e.XPos+2, e.YPos+2, false, false)
}

func (l *Label) UpdateMouseState(x, y int32, mouseBtnState uint32, frameTime float32) {
}

func (l *Label) DrawTooltip(rend *sdl.Renderer, font *ttf.Font) {
}
