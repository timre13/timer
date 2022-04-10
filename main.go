package main

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/gfx"
    "github.com/veandco/go-sdl2/ttf"
    "github.com/gen2brain/beeep"
    "fmt"
    "math"
    "path/filepath"
    . "timer/consts"
    "timer/common"
    "timer/button"
    "timer/confreader"
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

    gfx.FilledCircleColor(rend, CLOCK_CENT_X, CLOCK_CENT_Y, CLOCK_RAD, COLOR_CLOCK_BG)

    // Draw the arc
    {
        fgColor := common.LerpColors(&COLOR_CLOCK_FG, &COLOR_CLOCK_FG_DONE, elapsedPerc/100)

        const startRad = -math.Pi/2
        endRad := startRad+float64(elapsedPerc/100*math.Pi*2)

        xcoords := []int16{}
        ycoords := []int16{}

        // Calculate arc vertices
        increment := CLOCK_POLY_STEP
        for i := startRad; i <= endRad; i += increment {
            x := math.Cos(i)
            y := math.Sin(i)
            xcoords = append(xcoords, int16(x*CLOCK_RAD)+CLOCK_CENT_X)
            ycoords = append(ycoords, int16(y*CLOCK_RAD)+CLOCK_CENT_Y)
            if endRad-i < increment { // Switch to a smaller incrementation to prevent flickering at the ends
                increment = CLOCK_POLY_STEP_S
            }
        }

        // Add the center coords
        xcoords = append(xcoords, CLOCK_CENT_X)
        ycoords = append(ycoords, CLOCK_CENT_Y)

        gfx.FilledPolygonColor(rend, xcoords, ycoords, fgColor)
    }

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

var SESSTYPE_STRS = [...]string{"work", "break"}

func main() {
    exeDir := common.GetExeDir()

    var confPath string
    if filepath.IsAbs(CONF_PATH) { // Absolute path
        confPath = CONF_PATH
    } else { // Relative path
        confPath = filepath.Join(exeDir, CONF_PATH)
    }
    conf := confreader.LoadConf(confPath)

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

    var fullTimeMs float32 = float32(conf.GetSessLenMs(common.SESSION_TYPE_WORK))
    var elapsedTimeMs float32
    isPaused := false
    sessionType := common.SESSION_TYPE_WORK

    switchSessionType := func() {
        if sessionType == common.SESSION_TYPE_WORK {
            sessionType = common.SESSION_TYPE_BREAK
        } else if sessionType == common.SESSION_TYPE_BREAK {
            sessionType = common.SESSION_TYPE_WORK
        } else {
            panic(sessionType)
        }
    }

    pauseBtnImg := common.LoadImage(rend, filepath.Join(exeDir, "img/pause_btn.png"))
    startBtnImg := common.LoadImage(rend, filepath.Join(exeDir, "img/start_btn.png"))
    workSessionImg := common.LoadImage(rend, filepath.Join(exeDir, "img/work_icon.png"))
    breakSessionImg := common.LoadImage(rend, filepath.Join(exeDir, "img/break_icon.png"))

    pauseBtn := button.Button{
        CentX: CLOCK_CENT_X, CentY: BTN_CENT_Y,
        Radius: BTN_RAD,
        DefColor: &COLOR_BTN, HoverColor: &COLOR_BTN_HOVER, HoverBdColor: &COLOR_BTN_HOVER_BD}

    updatePauseBtnLabel := func() {
        if isPaused {
            pauseBtn.LabelImg = &startBtnImg
            pauseBtn.Tooltip = "Start"
        } else {
            pauseBtn.LabelImg = &pauseBtnImg
            pauseBtn.Tooltip = "Pause"
        }
    }
    updatePauseBtnLabel()

    pauseBtn.Callback = func(b *button.Button) {
        isPaused = !isPaused
        updatePauseBtnLabel()
    }


    sessTypeLabel := button.Button{
        CentX: CLOCK_CENT_X,
        CentY: CLOCK_CENT_Y+140,
        Radius: 24,
        DefColor: &COLOR_TRANSPARENT,
        HoverColor: &COLOR_TRANSPARENT,
        HoverBdColor: &COLOR_TRANSPARENT,
        UseDefCurs: true,
    }

    updateSessTypeLabel := func() {
        if sessionType == common.SESSION_TYPE_WORK {
            sessTypeLabel.LabelImg = &workSessionImg
            sessTypeLabel.Tooltip = "Work Session"
        } else if sessionType == common.SESSION_TYPE_BREAK {
            sessTypeLabel.LabelImg = &breakSessionImg
            sessTypeLabel.Tooltip = "Break Session"
        } else {
            panic(sessionType)
        }
    }
    updateSessTypeLabel()

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

        pauseBtn.UpdateMouseState(mouseX, mouseY, mouseState, fpsMan.RateTicks)
        pauseBtn.Draw(rend)
        pauseBtn.DrawTooltip(rend, tooltipFont)

        sessTypeLabel.UpdateMouseState(mouseX, mouseY, mouseState, fpsMan.RateTicks)
        sessTypeLabel.Draw(rend)
        sessTypeLabel.DrawTooltip(rend, tooltipFont)

        rend.Present()
        if !isPaused {
            elapsedTimeMs += fpsMan.RateTicks
        }
        if remTimeMs <= 0 && !isPaused {
            err = beeep.Notify("Timer", "End of "+SESSTYPE_STRS[sessionType]+" session", "")
            WARN_ERR(err)
            switchSessionType()
            elapsedTimeMs = 0
            fullTimeMs = float32(conf.GetSessLenMs(sessionType))
            // Request the user to click the pause button if it is configured
            switch sessionType {
            case common.SESSION_TYPE_WORK:
                isPaused = !conf.AutoStartWorkSess
            case common.SESSION_TYPE_BREAK:
                isPaused = !conf.AutoStartBreakSess
            }
            updatePauseBtnLabel()
            updateSessTypeLabel()
        }
        gfx.FramerateDelay(&fpsMan)
    }

    remTimeFont.Close()
    tooltipFont.Close()
    ttf.Quit()
    sdl.Quit()

    confreader.WriteConf(confPath, &conf)
}
