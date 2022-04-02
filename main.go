package main

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/gfx"
    "github.com/veandco/go-sdl2/ttf"
    "fmt"
    . "timer/consts"
    "timer/common"
    "timer/button"
)
var CHECK_ERR = common.CHECK_ERR

func drawClock(rend *sdl.Renderer, elapsedPerc float32) {
    if elapsedPerc < 0 {
        elapsedPerc = 0
    }
    if elapsedPerc > 100 {
        elapsedPerc = 100
    }
    fgColor := common.LerpColors(&COLOR_CLOCK_FG, &COLOR_CLOCK_FG_DONE, elapsedPerc/100)

    gfx.FilledCircleColor(rend, CLOCK_CENT_X, CLOCK_CENT_Y, CLOCK_RAD, COLOR_CLOCK_BG)
    gfx.FilledPieColor(rend, CLOCK_CENT_X, CLOCK_CENT_Y, CLOCK_RAD, -90, -90+int32(elapsedPerc/100*360), fgColor)
    gfx.FilledCircleColor(rend, CLOCK_CENT_X, CLOCK_CENT_Y, CLOCK_INS_RAD, COLOR_BG)
}

func drawRemTime(rend *sdl.Renderer, font *ttf.Font, remTimeMs int) {
    var remTimeStr string
    if remTimeMs <= 0 {
        remTimeStr = "--:--"
    } else {
        remTimeStr = fmt.Sprintf("%02d:%02d", remTimeMs/1000/60, remTimeMs/1000%60)
    }
    common.RenderText(rend, font, remTimeStr, &COLOR_FG, CLOCK_CENT_X, CLOCK_CENT_Y, true)
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

    var fullTimeMs float32
    var elapsedTimeMs float32
    fullTimeMs = 50000
    var isPaused bool

    // TODO: Set label image
    pauseBtn := button.Button{
        CentX: CLOCK_CENT_X, CentY: BTN_CENT_Y,
        Tooltip: "Pause",
        Callback: func(b *button.Button) {
            isPaused = !isPaused
        },
        Radius: BTN_RAD,
        DefColor: &COLOR_BTN, HoverColor: &COLOR_BTN_HOVER, HoverBdColor: &COLOR_BTN_HOVER_BD}

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

        remTimeMs := int(fullTimeMs-elapsedTimeMs)
        drawClock(rend, elapsedTimeMs/fullTimeMs*100.0)
        drawRemTime(rend, remTimeFont, remTimeMs)

        pauseBtn.UpdateMouseState(mouseX, mouseY, mouseState)
        pauseBtn.Draw(rend)
        pauseBtn.DrawTooltip(rend, tooltipFont)

        rend.Present()
        if !isPaused {
            elapsedTimeMs += fpsMan.RateTicks
        }
        if remTimeMs <= 0 {
            isPaused = true
        }
        gfx.FramerateDelay(&fpsMan)
    }

    remTimeFont.Close()
    tooltipFont.Close()
    ttf.Quit()
    sdl.Quit()
}
