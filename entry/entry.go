package entry

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/gfx"
    "github.com/veandco/go-sdl2/ttf"
    "timer/common"
    "timer/iwidget"
)

type Entry struct {
    XPos, YPos      int32           // Position of the top left corner
    Width           int32           // Width
    BgColor         *sdl.Color      // Background color
    FgColor         *sdl.Color      // Foreground color
    Font            *ttf.Font       // The font that is used to draw the text
    Text            string
    mouseX, mouseY  int32           // The absolute mouse position, set by `UpdateMousePos()`
    mouseBtnState   uint32          // Bitmask of pressed mouse buttons
    isMouseHovered  bool            // Set to true when the mouse is inside the button
    cursorCharPos   int
}
var _ iwidget.IWidget = (*Entry)(nil)

func (e *Entry) IsInside(x, y int32) bool {
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
            sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_IBEAM))
        } else {
            sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW))
        }
    }

    if e.isMouseHovered &&
    // If the left mouse button has just been pressed
    (e.mouseBtnState & sdl.ButtonLMask()) == 0 && (mouseBtnState & sdl.ButtonLMask()) != 0 {
    }
    e.mouseBtnState = mouseBtnState
}

func (e *Entry) Draw(rend *sdl.Renderer) {
    gfx.BoxColor(rend, e.XPos, e.YPos, e.XPos+e.Width, e.YPos+int32(e.Font.Height()), *e.BgColor)
    gfx.RectangleColor(rend, e.XPos, e.YPos, e.XPos+e.Width, e.YPos+int32(e.Font.Height()), *e.FgColor)

    cursXOffs := 0
    if e.Text != "" {
        common.RenderText(rend, e.Font, e.Text, e.FgColor, e.XPos+2, e.YPos, false, false)

        var err error
        cursXOffs, _, err = e.Font.SizeUTF8(e.Text[:e.cursorCharPos])
        common.PANIC_ERR(err)
    }
    gfx.BoxColor(rend, e.XPos+int32(cursXOffs)+2, e.YPos, e.XPos+int32(cursXOffs)+2, e.YPos+int32(e.Font.Height()), *e.FgColor)
}

func (e *Entry) DrawTooltip(rend *sdl.Renderer, font *ttf.Font) {
}

func (e *Entry) HandleTextInput(input string) {
    e.Text = e.Text[:e.cursorCharPos] + input + e.Text[e.cursorCharPos:]
    e.cursorCharPos++
}

func (e *Entry) HandleKeyPress(keycode sdl.Keycode) {
    switch keycode {
    case sdl.K_RIGHT:
        e.cursorCharPos++

    case sdl.K_LEFT:
        e.cursorCharPos--

    case sdl.K_BACKSPACE:
        if e.cursorCharPos > 0 {
            e.Text = e.Text[:e.cursorCharPos-1] + e.Text[e.cursorCharPos:]
            e.cursorCharPos--
        }

    case sdl.K_DELETE:
        if e.cursorCharPos < len(e.Text) {
            e.Text = e.Text[:e.cursorCharPos] + e.Text[e.cursorCharPos+1:]
        }
    }

    if e.cursorCharPos < 0 {
        e.cursorCharPos = 0
    } else if e.cursorCharPos > len(e.Text) {
        e.cursorCharPos = len(e.Text)
    }
}
