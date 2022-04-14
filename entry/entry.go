package entry

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/gfx"
    "github.com/veandco/go-sdl2/ttf"
    "timer/common"
    "timer/iwidget"
)

const ENTRY_CURS_BLINK_DELAY_MS = 500

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
    isFocused       bool
    isCursorShown   bool
    untilCursToggle float32
}
var _ iwidget.IWidget = (*Entry)(nil)

func (e *Entry) IsInside(x, y int32) bool {
    return x >= e.XPos && x < e.XPos+e.Width &&
           y >= e.YPos && y < e.YPos+int32(e.Font.Height())
}

func (e *Entry) UpdateMouseState(x, y int32, mouseBtnState uint32, frameTime float32) {
    e.mouseX = x
    e.mouseY = y
    isHovered := e.IsInside(x, y)
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

    e.mouseBtnState = mouseBtnState
}

func (e *Entry) Draw(rend *sdl.Renderer) {
    // Draw an extra outline if the widget has keyboard focus
    if e.isFocused {
        gfx.BoxColor(rend, e.XPos-1, e.YPos-1, e.XPos+e.Width+1, e.YPos+int32(e.Font.Height())+1, sdl.Color{R: 200, G: 200, B: 255, A: 255})
    }

    gfx.BoxColor(rend, e.XPos, e.YPos, e.XPos+e.Width, e.YPos+int32(e.Font.Height()), *e.BgColor)
    gfx.RectangleColor(rend, e.XPos, e.YPos, e.XPos+e.Width, e.YPos+int32(e.Font.Height()), *e.FgColor)

    cursXOffs := 0
    if e.Text != "" {
        common.RenderText(rend, e.Font, e.Text, e.FgColor, e.XPos+2, e.YPos, false, false)

        var err error
        cursXOffs, _, err = e.Font.SizeUTF8(e.Text[:e.cursorCharPos])
        common.PANIC_ERR(err)
    }
    if e.isFocused && e.isCursorShown {
        // Draw cursor
        gfx.BoxColor(rend, e.XPos+int32(cursXOffs)+2, e.YPos, e.XPos+int32(cursXOffs)+2, e.YPos+int32(e.Font.Height()), *e.FgColor)
    }
}

func (e *Entry) DrawTooltip(rend *sdl.Renderer, font *ttf.Font) {
}

func (e *Entry) HandleTextInput(input string) {
    if !e.isFocused {
        panic("Not focused Entry got text input")
    }

    e.Text = e.Text[:e.cursorCharPos] + input + e.Text[e.cursorCharPos:]
    e.cursorCharPos++
    e.isCursorShown = true
    e.untilCursToggle = ENTRY_CURS_BLINK_DELAY_MS
}

func (e *Entry) HandleKeyPress(keycode sdl.Keycode) {
    if !e.isFocused {
        panic("Not focused Entry got keyboard input")
    }

    showCursor := func() {
        e.isCursorShown = true
        e.untilCursToggle = ENTRY_CURS_BLINK_DELAY_MS
    }

    switch keycode {
    case sdl.K_RIGHT:
        e.cursorCharPos++
        showCursor()

    case sdl.K_LEFT:
        e.cursorCharPos--
        showCursor()

    case sdl.K_BACKSPACE:
        if e.cursorCharPos > 0 {
            e.Text = e.Text[:e.cursorCharPos-1] + e.Text[e.cursorCharPos:]
            e.cursorCharPos--
        }
        showCursor()

    case sdl.K_DELETE:
        if e.cursorCharPos < len(e.Text) {
            e.Text = e.Text[:e.cursorCharPos] + e.Text[e.cursorCharPos+1:]
        }
        showCursor()
    }

    if e.cursorCharPos < 0 {
        e.cursorCharPos = 0
    } else if e.cursorCharPos > len(e.Text) {
        e.cursorCharPos = len(e.Text)
    }
}

func (e *Entry) Tick(frameTime float32) {
    e.untilCursToggle -= frameTime
    if e.untilCursToggle < 0 {
        e.isCursorShown = !e.isCursorShown
        e.untilCursToggle = ENTRY_CURS_BLINK_DELAY_MS
    }
}

func (e *Entry) MoveCursToEnd() {
    e.cursorCharPos = len(e.Text)
}

func (e *Entry) SetFocused(focused bool) {
    e.isFocused = focused
    e.isCursorShown = true
    e.untilCursToggle = ENTRY_CURS_BLINK_DELAY_MS
}
