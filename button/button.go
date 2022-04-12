package button

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/gfx"
    "github.com/veandco/go-sdl2/ttf"
    "math"
    . "timer/consts"
    "timer/common"
)
var PANIC_ERR = common.PANIC_ERR

type Button struct {
    CentX, CentY    int32           // Center position
    Tooltip         string          // The text to display in a tooltip while hovering
    Radius          int32           // Radius
    Callback        func(*Button)   // The callback that is called when pressing the button
    LabelImg        *common.Image   // The image that is copied on top of the button
    DefColor        *sdl.Color      // Normal color
    HoverColor      *sdl.Color      // Color while hovered
    HoverBdColor    *sdl.Color      // Border color while hovered
    UseDefCurs      bool            // If true, the cursor is NOT set to a hand while hovering
    mouseX, mouseY  int32           // The absolute mouse position, set by `UpdateMousePos()`
    mouseBtnState   uint32          // Bitmask of pressed mouse buttons
    isMouseHovered  bool            // Set to true when the mouse is inside the button
    mouseHoverDurMs float32         // How long the button is hovered for
}

const BTN_TOOLT_DELAY_MS = 500

func (b *Button) isInside(x, y int32) bool {
    xDiff := float64(x-b.CentX)
    yDiff := float64(y-b.CentY)
    dist := math.Sqrt(xDiff*xDiff+yDiff*yDiff)
    return dist < float64(b.Radius)
}

func (b *Button) UpdateMouseState(x, y int32, mouseBtnState uint32, frameTime float32) {
    mouseMoved := (b.mouseX != x) || (b.mouseY != y)
    b.mouseX = x
    b.mouseY = y
    isHovered := b.isInside(x, y)
    mouseEnteredOrLeft := (b.isMouseHovered != isHovered)
    b.isMouseHovered = isHovered

    // Set the cursor when the mouse enters/leaves the button
    if mouseEnteredOrLeft {
        if b.isMouseHovered && !b.UseDefCurs {
            sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_HAND))
        } else {
            sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW))
        }
    }

    // We will show the tooltip with a delay and hide it when the mouse moved
    if b.isMouseHovered && !mouseMoved {
        b.mouseHoverDurMs += frameTime
    } else {
        b.mouseHoverDurMs = 0
    }

    if b.isMouseHovered &&
    // If the left mouse button has just been pressed
    (b.mouseBtnState & sdl.ButtonLMask()) == 0 && (mouseBtnState & sdl.ButtonLMask()) != 0 {
        // Call the callback if possible
        if b.Callback != nil {
            b.Callback(b)
        }
    }
    b.mouseBtnState = mouseBtnState
}

func (b *Button) Draw(rend *sdl.Renderer) {
    // Draw hover border
    if b.isMouseHovered {
        gfx.AACircleColor(rend, b.CentX, b.CentY, b.Radius+1, *b.HoverBdColor)
    }

    // Fill with different color when hovered
    if b.isMouseHovered {
        gfx.FilledCircleColor(rend, b.CentX, b.CentY, b.Radius, *b.HoverColor)
    } else {
        gfx.FilledCircleColor(rend, b.CentX, b.CentY, b.Radius, *b.DefColor)
    }

    // Don't draw out of the button and respect image size if smaller than button
    width := int32(math.Min(float64(b.Radius*2), float64(b.LabelImg.Width)))
    height := int32(math.Min(float64(b.Radius*2), float64(b.LabelImg.Height)))

    if b.LabelImg != nil && b.LabelImg.Img != nil {
        rend.Copy(b.LabelImg.Img, nil, &sdl.Rect{X: b.CentX-width/2, Y: b.CentY-height/2, W: width, H: height})
    }
}

func limit(x, min, max int) int {
    if x < min {
        return min
    } else if x > max {
        return max
    }
    return x
}

func (b *Button) DrawTooltip(rend *sdl.Renderer, font *ttf.Font) {
    if b.isMouseHovered && b.mouseHoverDurMs >= BTN_TOOLT_DELAY_MS {
        tooltipW, tooltipH, err := font.SizeUTF8(b.Tooltip)
        PANIC_ERR(err)
        tooltipX := int32(limit(int(b.mouseX)+20, 0, WIN_W-tooltipW))
        tooltipY := int32(limit(int(b.mouseY)+10, 0, WIN_H-tooltipH))

        gfx.RoundedBoxColor(rend, tooltipX, tooltipY, tooltipX+int32(tooltipW), tooltipY+int32(tooltipH), 2, COLOR_TOOLTIP_BG)
        PANIC_ERR(err)
        common.RenderText(rend, font, b.Tooltip, &COLOR_FG, tooltipX, tooltipY, false, false)
    }
}
