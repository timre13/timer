package main

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/gfx"
    "github.com/veandco/go-sdl2/ttf"
    "github.com/gen2brain/beeep"
    "fmt"
    . "timer/consts"
    "timer/common"
    "timer/button"
)
var PANIC_ERR = common.PANIC_ERR
var WARN_ERR = common.WARN_ERR

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

type SessionType int
const (
    SESSION_TYPE_WORK  SessionType = iota
    SESSION_TYPE_BREAK SessionType = iota
)

var SESSTYPE_STRS = [...]string{"work", "break"}
// Specifies how long the different types of sessions are
var SESS_DUR_MS = [...]float32{25*60*1000, 5*60*1000}

func main() {
    err := sdl.Init(sdl.INIT_VIDEO)
    PANIC_ERR(err)
    err = ttf.Init()
    PANIC_ERR(err)

    window, err := sdl.CreateWindow("Timer", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WIN_W, WIN_H, 0)
    PANIC_ERR(err)

    rend, err := sdl.CreateRenderer(window, 0, 0)
    PANIC_ERR(err)
    err = rend.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
    PANIC_ERR(err)

    remTimeFont, err := ttf.OpenFont(FONT_PATH, REM_TIME_FONT_SIZE)
    PANIC_ERR(err)
    tooltipFont, err := ttf.OpenFont(FONT_PATH, TOOLTIP_FONT_SIZE)
    PANIC_ERR(err)

    var fullTimeMs float32 = SESS_DUR_MS[SESSION_TYPE_WORK]
    var elapsedTimeMs float32
    var isPaused bool
    sessionType := SESSION_TYPE_WORK

    switchSessionType := func() {
        if sessionType == SESSION_TYPE_WORK {
            sessionType = SESSION_TYPE_BREAK
        } else if sessionType == SESSION_TYPE_BREAK {
            sessionType = SESSION_TYPE_WORK
        } else {
            panic(sessionType)
        }
    }

    pauseBtnImg := common.LoadImage(rend, "img/pause_btn.png")
    startBtnImg := common.LoadImage(rend, "img/start_btn.png")
    workSessionImg := common.LoadImage(rend, "img/work_icon.png")
    breakSessionImg := common.LoadImage(rend, "img/break_icon.png")

    pauseBtn := button.Button{
        CentX: CLOCK_CENT_X, CentY: BTN_CENT_Y,
        Tooltip: "Pause",
        LabelImg: &pauseBtnImg,
        Radius: BTN_RAD,
        DefColor: &COLOR_BTN, HoverColor: &COLOR_BTN_HOVER, HoverBdColor: &COLOR_BTN_HOVER_BD}
    pauseBtn.Callback = func(b *button.Button) {
        isPaused = !isPaused
        if isPaused {
            pauseBtn.LabelImg = &startBtnImg
        } else {
            pauseBtn.LabelImg = &pauseBtnImg
        }
    }

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
        var sessImg *common.Image
        if sessionType == SESSION_TYPE_WORK {
            sessImg = &workSessionImg
        } else if sessionType == SESSION_TYPE_BREAK {
            sessImg = &breakSessionImg
        } else {
            panic(sessionType)
        }
        rend.Copy(sessImg.Img, nil, &sdl.Rect{
            X: CLOCK_CENT_X-sessImg.Width/2,
            Y: CLOCK_CENT_Y+140-sessImg.Height,
            W: sessImg.Width,
            H: sessImg.Height})

        pauseBtn.UpdateMouseState(mouseX, mouseY, mouseState)
        pauseBtn.Draw(rend)
        pauseBtn.DrawTooltip(rend, tooltipFont)

        rend.Present()
        if !isPaused {
            elapsedTimeMs += fpsMan.RateTicks
        }
        if remTimeMs <= 0 && !isPaused {
            err = beeep.Notify("Timer", "End of "+SESSTYPE_STRS[sessionType]+" session", "")
            WARN_ERR(err)
            switchSessionType()
            elapsedTimeMs = 0
            fullTimeMs = SESS_DUR_MS[sessionType]
            // TODO: Wait for user to continue
        }
        gfx.FramerateDelay(&fpsMan)
    }

    remTimeFont.Close()
    tooltipFont.Close()
    ttf.Quit()
    sdl.Quit()
}
