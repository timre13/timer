package checkbox

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/gfx"
    "github.com/veandco/go-sdl2/ttf"
    "timer/iwidget"
)

const CHECKBOX_SIZE = 16

type CheckBox struct {
    XPos, YPos      int32           // Position of the top left corner
    BgColor         *sdl.Color      // Background color
    FgColor         *sdl.Color      // Foreground color
    HoverBgColor    *sdl.Color      // Color while hovered
    HoverBdColor    *sdl.Color      // Border color while hovered
    Value           bool
    mouseX, mouseY  int32           // The absolute mouse position, set by `UpdateMousePos()`
    mouseBtnState   uint32          // Bitmask of pressed mouse buttons
    isMouseHovered  bool            // Set to true when the mouse is inside the button
}
var _ iwidget.IWidget = (*CheckBox)(nil)

func (c *CheckBox) isInside(x, y int32) bool {
    return x >= c.XPos && x < c.XPos+CHECKBOX_SIZE &&
           y >= c.YPos && y < c.YPos+CHECKBOX_SIZE
}

func (c *CheckBox) UpdateMouseState(x, y int32, mouseBtnState uint32, frameTime float32) {
    c.mouseX = x
    c.mouseY = y
    isHovered := c.isInside(x, y)
    mouseEnteredOrLeft := (c.isMouseHovered != isHovered)
    c.isMouseHovered = isHovered

    // Set the cursor when the mouse enters/leaves the button
    if mouseEnteredOrLeft {
        if c.isMouseHovered {
            sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_HAND))
        } else {
            sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW))
        }
    }

    if c.isMouseHovered &&
    // If the left mouse button has just been pressed
    (c.mouseBtnState & sdl.ButtonLMask()) == 0 && (mouseBtnState & sdl.ButtonLMask()) != 0 {
        // Change value
        c.Value = !c.Value
    }
    c.mouseBtnState = mouseBtnState
}

func (c  *CheckBox) Draw(rend *sdl.Renderer) {
    x2 := c.XPos+CHECKBOX_SIZE
    y2 := c.YPos+CHECKBOX_SIZE

    if c.isMouseHovered {
        gfx.BoxColor(rend, c.XPos, c.YPos, x2, y2, *c.HoverBgColor)
    } else {
        gfx.BoxColor(rend, c.XPos, c.YPos, x2, y2, *c.BgColor)
    }

    if c.Value {
        gfx.LineColor(rend, c.XPos+2, c.YPos+2, x2-2, y2-2, *c.FgColor)
        gfx.LineColor(rend, x2-2, c.YPos+2, c.XPos+2, y2-2, *c.FgColor)
    }
}

func (c *CheckBox) DrawTooltip(rend *sdl.Renderer, font *ttf.Font) {
}
