package main

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/gfx"
    "github.com/veandco/go-sdl2/ttf"
    "github.com/gen2brain/beeep"
    "fmt"
    "math"
    "strconv"
    "errors"
    . "timer/consts"
    "timer/common"
    "timer/iwidget"
    "timer/button"
    "timer/entry"
    "timer/label"
    "timer/checkbox"
    "timer/confreader"
    "timer/stats"
)
var PANIC_ERR = common.PANIC_ERR
var WARN_ERR = common.WARN_ERR

func UNUSED(interface{}) {}

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

func getTextWidth(font *ttf.Font, text string) int32 {
    w, _, err := font.SizeUTF8(text)
    PANIC_ERR(err)
    return int32(w)
}

var SESSTYPE_STRS = [...]string{"work", "break"}

func createConfWinWidgets(
        confWinFont *ttf.Font, conf *confreader.Config, win *sdl.Window,
        confWidgetPtrs *[]iwidget.IWidget, focusedConfWidgetPtr *iwidget.IWidget, okBtnImg *common.Image, cancelBtnImg *common.Image) {

    var _widget iwidget.IWidget
    lineI := int32(0)
    addLabelWidget := func(text string, center bool) {
        _widget = &label.Label{Font: confWinFont, XPos: 30, YPos: 30+int32(float32(confWinFont.Height())*1.2)*lineI,
                Text: text, FgColor: &COLOR_FG}
        if center {
            _widget.(*label.Label).XPos = WIN_W/2-getTextWidth(confWinFont, _widget.(*label.Label).Text)/2 // Center label
        }
        *confWidgetPtrs = append(*confWidgetPtrs, _widget)
        lineI++
    }

    addEntryWidget := func(value string) *entry.Entry {
        _widget = &entry.Entry{Font: confWinFont, XPos: 280, YPos: 30+int32(float32(confWinFont.Height())*1.2)*(lineI-1),
                Width: 30, BgColor: &COLOR_BTN, FgColor: &COLOR_FG}
        _widget.(*entry.Entry).Text = value
        _widget.(*entry.Entry).MoveCursToEnd()
        *confWidgetPtrs = append(*confWidgetPtrs, _widget)
        return _widget.(*entry.Entry)
    }

    addCheckboxWidget := func(value bool) *checkbox.CheckBox {
        _widget = &checkbox.CheckBox{XPos: 280, YPos: 30+int32(float32(confWinFont.Height())*1.2)*(lineI-1),
                BgColor: &COLOR_BTN, FgColor: &COLOR_FG, HoverBgColor: &COLOR_BTN_HOVER, HoverBdColor: &COLOR_BTN_HOVER_BD}
        _widget.(*checkbox.CheckBox).Value = value
        *confWidgetPtrs = append(*confWidgetPtrs, _widget)
        return _widget.(*checkbox.CheckBox)
    }


    // TODO: Use frame widget to separate widgets

    addLabelWidget("SETTINGS", true)
    lineI++
    addLabelWidget("Work Session", false)
    addLabelWidget("    Duration:", false)
    workSessDurMinEntry := addEntryWidget(fmt.Sprint(conf.WorkSessDurMin))
    { // Add suffix label
        _widget = &label.Label{Font: confWinFont, XPos: 312, YPos: 30+int32(float32(confWinFont.Height())*1.2)*(lineI-1),
                Text: "minutes", FgColor: &COLOR_FG}
        *confWidgetPtrs = append(*confWidgetPtrs, _widget)
    }
    addLabelWidget("    Auto start:", false)
    autoStartWorkSessCheckb := addCheckboxWidget(conf.AutoStartWorkSess)
    lineI++
    addLabelWidget("Break Session", false)
    addLabelWidget("    Duration:", false)
    { // Add suffix label
        _widget = &label.Label{Font: confWinFont, XPos: 312, YPos: 30+int32(float32(confWinFont.Height())*1.2)*(lineI-1),
                Text: "minutes", FgColor: &COLOR_FG}
        *confWidgetPtrs = append(*confWidgetPtrs, _widget)
    }
    breakSessDurMinEntry := addEntryWidget(fmt.Sprint(conf.BreakSessDurMin))
    addLabelWidget("    Auto start:", false)
    autoStartBreakSessCheckb := addCheckboxWidget(conf.AutoStartBreakSess)
    lineI++
    addLabelWidget( "Misc", false)
    addLabelWidget( "    Show notif.:", false)
    sessEndShowNotifCheckb := addCheckboxWidget(conf.SessEndShowNotif)
    addLabelWidget( "    Show rem. time in title:", false)
    showRemTimeInWinTitleCb := addCheckboxWidget(conf.ShowRemTimeInWinTitle)

    *focusedConfWidgetPtr = workSessDurMinEntry
    workSessDurMinEntry.SetFocused(true)

    okButton := button.Button{CentX: WIN_W/3, CentY: WIN_H-50, Tooltip: "Save settings and close menu", Radius: BTN_SMALL_RAD,
            LabelImg: okBtnImg,
            DefColor: &COLOR_BTN, HoverColor: &COLOR_BTN_HOVER, HoverBdColor: &COLOR_BTN_HOVER_BD}
    okButton.Callback = func(*button.Button) {
        strToInt := func(str string, minVal int, maxVal int) (int, error) {
            intVal, err := strconv.Atoi(str)
            if err != nil {
                sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Error applying settings", "Invalid number: \""+str+"\"", win)
                return 0, err
            }
            if intVal < minVal || intVal > maxVal {
                sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "Error applying settings", "Integer "+fmt.Sprint(intVal)+" out of range", win)
                return 0, errors.New("")
            }
            return intVal, nil
        }
        var err error

        // We will work in a temporary copy
        tconf := *conf

        tconf.WorkSessDurMin, err = strToInt(workSessDurMinEntry.Text, 1, 999)
        if err != nil { goto Error }
        tconf.AutoStartWorkSess = autoStartWorkSessCheckb.Value
        tconf.BreakSessDurMin, err = strToInt(breakSessDurMinEntry.Text, 1, 999)
        if err != nil { goto Error }
        tconf.AutoStartBreakSess = autoStartBreakSessCheckb.Value
        tconf.SessEndShowNotif = sessEndShowNotifCheckb.Value
        tconf.ShowRemTimeInWinTitle = showRemTimeInWinTitleCb.Value

        // Overwrite the old config with the new
        *conf = tconf
        if !conf.ShowRemTimeInWinTitle {
            win.SetTitle(WIN_TITLE)
        }
        sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW))
        *confWidgetPtrs = []iwidget.IWidget{} // Close window
        fmt.Println("Overwrote config")
        Error:
        // We skip the overwriting of the old config if an error happens
        // We also leave the window open
    }
    *confWidgetPtrs = append(*confWidgetPtrs, &okButton)

    cancelButton := button.Button{CentX: WIN_W/3*2, CentY: WIN_H-50, Tooltip: "Close menu without saving", Radius: BTN_SMALL_RAD,
            Callback: nil, LabelImg: cancelBtnImg,
            DefColor: &COLOR_BTN, HoverColor: &COLOR_BTN_HOVER, HoverBdColor: &COLOR_BTN_HOVER_BD}
    cancelButton.Callback = func(*button.Button) {
        sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW))
        *confWidgetPtrs = []iwidget.IWidget{} // Close window
    }
    *confWidgetPtrs = append(*confWidgetPtrs, &cancelButton)
}

// TODO: Taskbar icon
// TODO: Hiding with taskbar icon

func main() {
    exeDir := common.GetExeDir()

    confPath := common.GetRealPath(exeDir, CONF_PATH)
    conf := confreader.LoadConf(confPath)

    statPath := common.GetRealPath(exeDir, STAT_PATH)
    stat := stats.LoadStats(statPath)

    // FIXME: Handle midnight (the date changes) while the program is running
    currDate := stats.GetCurrentDate()
    fmt.Printf("Today's date: %+v\n", currDate)
    todayStats := stat.GetDay(&currDate)
    fmt.Printf("Today's stats: %+v\n", todayStats)

    err := sdl.Init(sdl.INIT_VIDEO)
    PANIC_ERR(err)
    err = ttf.Init()
    PANIC_ERR(err)

    window, err := sdl.CreateWindow(WIN_TITLE, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WIN_W, WIN_H, 0)
    PANIC_ERR(err)

    rend, err := sdl.CreateRenderer(window, 0, 0)
    PANIC_ERR(err)
    err = rend.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
    PANIC_ERR(err)

    remTimeFont, err := ttf.OpenFont(FONT_PATH, REM_TIME_FONT_SIZE)
    PANIC_ERR(err)
    tooltipFont, err := ttf.OpenFont(FONT_PATH, TOOLTIP_FONT_SIZE)
    PANIC_ERR(err)
    confWinFont, err := ttf.OpenFont(FONT_PATH, CONFWIN_FONT_SIZE)
    PANIC_ERR(err)

    var elapsedTimeMs float32
    isPaused := true
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

    pauseBtnImg         := common.LoadImage(rend, common.GetRealPath(exeDir, "img/pause_btn.png"))
    startBtnImg         := common.LoadImage(rend, common.GetRealPath(exeDir, "img/start_btn.png"))
    settingsBtnImg      := common.LoadImage(rend, common.GetRealPath(exeDir, "img/settings_btn.png"))
    okBtnImg            := common.LoadImage(rend, common.GetRealPath(exeDir, "img/ok_btn.png"))
    cancelBtnImg        := common.LoadImage(rend, common.GetRealPath(exeDir, "img/cancel_btn.png"))
    workSessionImg      := common.LoadImage(rend, common.GetRealPath(exeDir, "img/work_icon.png"))
    breakSessionImg     := common.LoadImage(rend, common.GetRealPath(exeDir, "img/break_icon.png"))

    confWidgetPtrs := []iwidget.IWidget{}
    var focusedConfWidgetPtr iwidget.IWidget

    widgetPtrs := []iwidget.IWidget{}

    pauseBtn := button.Button{
        CentX: CLOCK_CENT_X, CentY: BTN_CENT_Y,
        Radius: BTN_RAD,
        DefColor: &COLOR_BTN, HoverColor: &COLOR_BTN_HOVER, HoverBdColor: &COLOR_BTN_HOVER_BD}
    widgetPtrs = append(widgetPtrs, &pauseBtn)

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


    settingsBtn := button.Button{
        CentX: WIN_W-BTN_SMALL_RAD-4, CentY: BTN_SMALL_RAD+4,
        Radius: BTN_SMALL_RAD,
        Tooltip: "Settings",
        Callback: func(*button.Button) {
            sdl.SetCursor(sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW))
            createConfWinWidgets(confWinFont, &conf, window, &confWidgetPtrs, &focusedConfWidgetPtr, &okBtnImg, &cancelBtnImg)
        },
        LabelImg: &settingsBtnImg,
        DefColor: &COLOR_BTN, HoverColor: &COLOR_BTN_HOVER, HoverBdColor: &COLOR_BTN_HOVER_BD}
    widgetPtrs = append(widgetPtrs, &settingsBtn)


    sessTypeLabel := button.Button{
        CentX: CLOCK_CENT_X,
        CentY: CLOCK_CENT_Y+140,
        Radius: 24,
        DefColor: &COLOR_TRANSPARENT,
        HoverColor: &COLOR_TRANSPARENT,
        HoverBdColor: &COLOR_TRANSPARENT,
        UseDefCurs: true,
    }
    widgetPtrs = append(widgetPtrs, &sessTypeLabel)

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

    remTimeLabelW, remTimeLabelH, err := remTimeFont.SizeUTF8("00:00")
    PANIC_ERR(err)
    remTimeLabel := label.Label{
        Text: "--:--",
        XPos: CLOCK_CENT_X-int32(remTimeLabelW)/2, YPos: CLOCK_CENT_Y-int32(remTimeLabelH)/2,
        Font: remTimeFont, FgColor: &COLOR_FG}
    widgetPtrs = append(widgetPtrs, &remTimeLabel)


    sdl.StartTextInput()
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
                // Update elapsed time counter with the remaining time
                switch sessionType {
                case common.SESSION_TYPE_WORK:
                    todayStats.WorkMs += elapsedTimeMs

                case common.SESSION_TYPE_BREAK:
                    todayStats.BreakMs += elapsedTimeMs
                }
                running = false

            case sdl.TEXTINPUT:
                if len(confWidgetPtrs) != 0 && focusedConfWidgetPtr != nil {
                    focusedConfWidgetPtr.(*entry.Entry).HandleTextInput(event.(*sdl.TextInputEvent).GetText())
                }

            case sdl.KEYDOWN:
                if len(confWidgetPtrs) != 0 && focusedConfWidgetPtr != nil {
                    focusedConfWidgetPtr.(*entry.Entry).HandleKeyPress(event.(*sdl.KeyboardEvent).Keysym.Sym)
                }

            case sdl.MOUSEBUTTONDOWN:
                if len(confWidgetPtrs) != 0 {
                    for _, w := range confWidgetPtrs {
                        switch w.(type) {
                        case *entry.Entry:
                            entryw := w.(*entry.Entry)
                            mx := event.(*sdl.MouseButtonEvent).X
                            my := event.(*sdl.MouseButtonEvent).Y
                            if entryw.IsInside(mx, my) {
                                focusedConfWidgetPtr.(*entry.Entry).SetFocused(false)
                                entryw.SetFocused(true)
                                focusedConfWidgetPtr = w
                            }
                        }
                    }
                }
            }
        }
        if !running {
            break
        }

        mouseX, mouseY, mouseState := sdl.GetMouseState()

        rend.SetDrawColor(COLOR_BG.R, COLOR_BG.G, COLOR_BG.B, COLOR_BG.A)
        rend.Clear()

        fullTimeMs := float32(conf.GetSessLenMs(sessionType))
        remTimeMs := int(fullTimeMs-elapsedTimeMs)

        // The config menu is displayed when `confWidgetPtrs` has widgets
        if len(confWidgetPtrs) == 0 {
            drawClock(rend, elapsedTimeMs/fullTimeMs*100.0)
            remTimeStr := fmt.Sprintf("%02d:%02d", remTimeMs/1000/60, remTimeMs/1000%60)
            remTimeLabel.Text = remTimeStr

            if conf.ShowRemTimeInWinTitle {
                title := WIN_TITLE+" - " + remTimeStr
                if window.GetTitle() != title {
                    window.SetTitle(title)
                }
            }

            for _, w := range widgetPtrs {
                w.UpdateMouseState(mouseX, mouseY, mouseState, fpsMan.RateTicks)
            }
            for _, w := range widgetPtrs {
                w.Draw(rend)
            }
            for _, w := range widgetPtrs {
                w.DrawTooltip(rend, tooltipFont)
            }
        } else {
            gfx.RoundedBoxColor(rend, 20, 20, WIN_W-20, WIN_H-20, 4, sdl.Color{R: COLOR_CLOCK_BG.R, G: COLOR_CLOCK_BG.G, B: COLOR_CLOCK_BG.B, A: 250})
            gfx.RoundedRectangleColor(rend, 20, 20, WIN_W-20, WIN_H-20, 4, sdl.Color{R: 255, G: 255, B: 255, A: 250})

            for _, w := range confWidgetPtrs {
                switch w.(type) {
                case *entry.Entry:
                    // Handle cursor blinking
                    w.(*entry.Entry).Tick(fpsMan.RateTicks)
                }
            }

            for _, w := range confWidgetPtrs {
                w.UpdateMouseState(mouseX, mouseY, mouseState, fpsMan.RateTicks)
            }
            for _, w := range confWidgetPtrs {
                w.Draw(rend)
            }
            for _, w := range confWidgetPtrs {
                w.DrawTooltip(rend, tooltipFont)
            }
        }

        rend.Present()
        if !isPaused {
            elapsedTimeMs += fpsMan.RateTicks
        }
        if remTimeMs <= 0 && !isPaused {
            if conf.SessEndShowNotif {
                err = beeep.Notify(WIN_TITLE, "End of "+SESSTYPE_STRS[sessionType]+" session", "")
                WARN_ERR(err)
            }
            // Update elapsed time counter
            switch sessionType {
            case common.SESSION_TYPE_WORK:
                todayStats.WorkMs += elapsedTimeMs

            case common.SESSION_TYPE_BREAK:
                todayStats.BreakMs += elapsedTimeMs
            }
            switchSessionType()
            elapsedTimeMs = 0
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
    if !todayStats.IsEmpty() {
        stat[currDate] = todayStats
    }
    stats.WriteStats(statPath, &stat)
}
