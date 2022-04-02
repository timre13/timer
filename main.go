package main

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/gfx"
    "github.com/veandco/go-sdl2/ttf"
    "math"
    "fmt"
)

func CHECK_ERR(err error) {
    if err != nil {
        sdl.Quit()
        panic(err)
    }
}

const TARGET_FPS        = 30

const WIN_W             = 500
const WIN_H             = 600
const CLOCK_CENT_X      = 250
const CLOCK_CENT_Y      = 250
const CLOCK_RAD         = 230
const CLOCK_INS_RAD     = 170
const BTN_RAD           = 40
const BTN_CENT_Y        = WIN_H-BTN_RAD-20

const FONT_PATH             = "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf"
const REM_TIME_FONT_SIZE    = 36
const TOOLTIP_FONT_SIZE     = 18

var COLOR_BG            = sdl.Color{R:   4, G:  13, B:  35, A: 255}
var COLOR_FG            = sdl.Color{R: 247, G: 247, B: 255, A: 255}
var COLOR_CLOCK_BG      = sdl.Color{R:  87, G: 115, B: 153, A: 255}
var COLOR_CLOCK_FG      = sdl.Color{R: 254, G:  95, B:  85, A: 255}
var COLOR_CLOCK_FG_DONE = sdl.Color{R:  86, G: 227, B: 159, A: 255}
var COLOR_BTN           = sdl.Color{R:  37, G:  49, B:  65, A: 255}
var COLOR_BTN_HOVER     = sdl.Color{R:  52, G:  69, B:  91, A: 255}
var COLOR_BTN_HOVER_BD  = sdl.Color{R: 244, G: 185, B:  66, A: 255}
var COLOR_TOOLTIP_BG    = sdl.Color{R:  82, G: 108, B: 142, A: 200}

func lerpInt(a int, b int, t float32) int {
    return int(math.Round(float64(float32(a) + t * float32(b - a))))
}

func lerpColors(x *sdl.Color, y *sdl.Color, t float32) sdl.Color {
    return sdl.Color{
        R: uint8(lerpInt(int(x.R), int(y.R), t)),
        G: uint8(lerpInt(int(x.G), int(y.G), t)),
        B: uint8(lerpInt(int(x.B), int(y.B), t)),
        A: 255,
    }
}

func drawClock(rend *sdl.Renderer, elapsedPerc float32) {
    if elapsedPerc < 0 {
        elapsedPerc = 0
    }
    if elapsedPerc > 100 {
        elapsedPerc = 100
    }
    fgColor := lerpColors(&COLOR_CLOCK_FG, &COLOR_CLOCK_FG_DONE, elapsedPerc/100)

    gfx.FilledCircleColor(rend, CLOCK_CENT_X, CLOCK_CENT_Y, CLOCK_RAD, COLOR_CLOCK_BG)
    gfx.FilledPieColor(rend, CLOCK_CENT_X, CLOCK_CENT_Y, CLOCK_RAD, -90, -90+int32(elapsedPerc/100*360), fgColor)
    gfx.FilledCircleColor(rend, CLOCK_CENT_X, CLOCK_CENT_Y, CLOCK_INS_RAD, COLOR_BG)
}

func renderText(rend *sdl.Renderer, font *ttf.Font, str string, color *sdl.Color, x, y int32, areCoordsCent bool) {
    textSurf, err := font.RenderUTF8Blended(str, *color)
    CHECK_ERR(err)
    textTex, err := rend.CreateTextureFromSurface(textSurf)
    CHECK_ERR(err)

    var rectX, rectY int32
    if areCoordsCent {
        rectX = x-textSurf.W/2
        rectY = y-textSurf.H/2
    } else {
        rectX = x
        rectY = y
    }

    dstRect := sdl.Rect{
        X: rectX, Y: rectY,
        W: textSurf.W, H: textSurf.H}
    rend.Copy(textTex, nil, &dstRect)

    textSurf.Free()
    textTex.Destroy()
}

func drawRemTime(rend *sdl.Renderer, font *ttf.Font, remTimeMs int) {
    var remTimeStr string
    if remTimeMs <= 0 {
        remTimeStr = "--:--"
    } else {
        remTimeStr = fmt.Sprintf("%02d:%02d", remTimeMs/1000/60, remTimeMs/1000%60)
    }
    renderText(rend, font, remTimeStr, &COLOR_FG, CLOCK_CENT_X, CLOCK_CENT_Y, true)
}

type Button struct {
    // TODO: Label image
    CentX, CentY    int32           // Center position
    Tooltip         string          // The text to display in a tooltip while hovering
    Radius          int32           // Radius
    Callback        func(*Button)   // The callback that is called when pressing the button
    DefColor        *sdl.Color      // Normal color
    HoverColor      *sdl.Color      // Color while hovered
    HoverBdColor    *sdl.Color      // Border color while hovered
    mouseX, mouseY  int32           // The absolute mouse position, set by `UpdateMousePos()`
    mouseBtnState   uint32          // Bitmask of pressed mouse buttons
    isMouseHovered  bool            // Set to true when the mouse is inside the button
}

func (b *Button) isInside(x, y int32) bool {
    xDiff := float64(x-b.CentX)
    yDiff := float64(y-b.CentY)
    dist := math.Sqrt(xDiff*xDiff+yDiff*yDiff)
    return dist < float64(b.Radius)
}

func (b *Button) UpdateMouseState(x, y int32, mouseBtnState uint32) {
    b.mouseX = x
    b.mouseY = y
    if b.isMouseHovered &&
    // If the left mouse button has just been pressed
    (b.mouseBtnState & sdl.ButtonLMask()) == 0 && (mouseBtnState & sdl.ButtonLMask()) != 0 {
        // Call the callback if possible
        if b.Callback != nil {
            b.Callback(b)
        }
    }
    b.mouseBtnState = mouseBtnState
    b.isMouseHovered = b.isInside(x, y)
}

func (b *Button) Draw(rend *sdl.Renderer) {
    if b.isMouseHovered {
        gfx.AACircleColor(rend, b.CentX, b.CentY, b.Radius+1, *b.HoverBdColor)
    }

    if b.isMouseHovered {
        gfx.FilledCircleColor(rend, b.CentX, b.CentY, b.Radius, *b.HoverColor)
    } else {
        gfx.FilledCircleColor(rend, b.CentX, b.CentY, b.Radius, *b.DefColor)
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
    if b.isMouseHovered {
        tooltipW, tooltipH, err := font.SizeUTF8(b.Tooltip)
        CHECK_ERR(err)
        tooltipX := int32(limit(int(b.mouseX)+20, 0, WIN_W-tooltipW))
        tooltipY := int32(limit(int(b.mouseY)+10, 0, WIN_H-tooltipH))

        gfx.RoundedBoxColor(rend, tooltipX, tooltipY, tooltipX+int32(tooltipW), tooltipY+int32(tooltipH), 2, COLOR_TOOLTIP_BG)
        CHECK_ERR(err)
        renderText(rend, font, b.Tooltip, &COLOR_FG, tooltipX, tooltipY, false)
    }
}

func main() {
    err := sdl.Init(sdl.INIT_VIDEO)
    CHECK_ERR(err)
    err = ttf.Init()
    CHECK_ERR(err)

    window, err := sdl.CreateWindow("Timer", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WIN_W, WIN_H, 0)
    CHECK_ERR(err)

    rend, err := sdl.CreateRenderer(window, 0, 0)
    CHECK_ERR(err)
    err = rend.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
    CHECK_ERR(err)

    remTimeFont, err := ttf.OpenFont(FONT_PATH, REM_TIME_FONT_SIZE)
    CHECK_ERR(err)
    tooltipFont, err := ttf.OpenFont(FONT_PATH, TOOLTIP_FONT_SIZE)
    CHECK_ERR(err)

    // TODO: Make it work
    // TODO: Set label image
    pauseBtn := Button{
        CentX: CLOCK_CENT_X, CentY: BTN_CENT_Y,
        Tooltip: "Pause",
        Radius: BTN_RAD,
        DefColor: &COLOR_BTN, HoverColor: &COLOR_BTN_HOVER, HoverBdColor: &COLOR_BTN_HOVER_BD}

    var fullTimeMs float32
    var elapsedTimeMs float32
    fullTimeMs = 50000

    fpsMan := gfx.FPSmanager{}
    gfx.InitFramerate(&fpsMan)
    gfx.SetFramerate(&fpsMan, TARGET_FPS)
    running := true
    for {
        for {
            event := sdl.PollEvent()
            if event == nil {
                break
            }

            switch event.GetType() {
            case sdl.QUIT:
                running = false
            }
        }
        if !running {
            break
        }

        mouseX, mouseY, mouseState := sdl.GetMouseState()

        rend.SetDrawColor(COLOR_BG.R, COLOR_BG.G, COLOR_BG.B, COLOR_BG.A)
        rend.Clear()

        drawClock(rend, elapsedTimeMs/fullTimeMs*100.0)
        drawRemTime(rend, remTimeFont, int(fullTimeMs-elapsedTimeMs))
        pauseBtn.Draw(rend, mouseX, mouseY)
        pauseBtn.DrawTooltip(rend, tooltipFont, mouseX, mouseY)
        pauseBtn.UpdateMouseState(mouseX, mouseY, mouseState)
        pauseBtn.Draw(rend)
        pauseBtn.DrawTooltip(rend, tooltipFont)

        rend.Present()
        elapsedTimeMs += fpsMan.RateTicks
        gfx.FramerateDelay(&fpsMan)
    }

    remTimeFont.Close()
    tooltipFont.Close()
    ttf.Quit()
    sdl.Quit()
}
