package entry

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/gfx"
    "github.com/veandco/go-sdl2/ttf"
    . "timer/consts"
    "timer/iwidget"
)

type Entry struct {
    XPos, YPos      int32           // Position of the top left corner
    Width           int32           // Width
    BgColor         *sdl.Color      // Background color
    FgColor         *sdl.Color      // Foreground color
    Font            *ttf.Font       // The font that is used to draw the text
    mouseX, mouseY  int32           // The absolute mouse position, set by `UpdateMousePos()`
    mouseBtnState   uint32          // Bitmask of pressed mouse buttons
    isMouseHovered  bool            // Set to true when the mouse is inside the button
}
var _ iwidget.IWidget = (*Entry)(nil)

func (e *Entry) isInside(x, y int32) bool {
    return x >= e.XPos && x < e.XPos+e.Width &&
           y >= e.YPos && y < e.YPos+int32(e.Font.Height())
}

func (e *Entry) UpdateMouseState(x, y int32, mouseBtnState uint32, frameTime float32) {
    e.mouseX = x
    e.mouseY = y
    isHovered := e.isInside(x, y)
    mouseEnteredOrLeft := (e.isMouseHovered != isHovered)
    e.isMouseHovered = isHovered

    // Set the cursor when the mouse enters/leaves the button
    if mouseEnteredOrLeft {
        if e.isMouseHovered {
            sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_HAND))
        } else {
            sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW))
        }
    }

    if e.isMouseHovered &&
    // If the left mouse button has just been pressed
    (e.mouseBtnState & sdl.ButtonLMask()) == 0 && (mouseBtnState & sdl.ButtonLMask()) != 0 {
        // Call the callback if possible
        //if e.Callback != nil {
        //    e.Callback(e)
        //}
    }
    e.mouseBtnState = mouseBtnState
}

func (e *Entry) Draw(rend *sdl.Renderer) {
    //gfx.FilledCircleColor(rend, e.CentX, e.CentY, e.Radius, *e.BgColor)
    gfx.BoxColor(rend, e.XPos, e.YPos, e.XPos+e.Width, e.YPos+int32(e.Font.Height())+4, COLOR_BTN)

    //// Don't draw out of the button and respect image size if smaller than button
    //width := int32(math.Min(float64(e.Radius*2), float64(e.LabelImg.Width)))
    //height := int32(math.Min(float64(e.Radius*2), float64(e.LabelImg.Height)))

    //if e.LabelImg != nil && e.LabelImg.Img != nil {
    //    rend.Copy(e.LabelImg.Img, nil, &sdl.Rect{X: e.CentX-width/2, Y: e.CentY-height/2, W: width, H: height})
    //}
}

func (e *Entry) DrawTooltip(rend *sdl.Renderer, font *ttf.Font) {
}
