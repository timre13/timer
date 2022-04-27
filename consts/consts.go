package consts

import "github.com/veandco/go-sdl2/sdl"

const CONF_PATH         = "config.json"
const STAT_PATH         = "stats.json"

const TARGET_FPS        = 30

const WIN_W             = 500
const WIN_H             = 600
const STATS_WIN_W       = 1000
const STATS_WIN_H       = 800
const CLOCK_CENT_X      = 250
const CLOCK_CENT_Y      = 250
const CLOCK_RAD         = 230
const CLOCK_INS_RAD     = 170
const CLOCK_POLY_STEP   = 0.02 // The less the step, the more polygons the clock will have
const CLOCK_POLY_STEP_S = 0.001 // The less the small step, the less flickering the clock will have at the arc's ends
const BTN_RAD           = 40
const BTN_CENT_Y        = WIN_H-BTN_RAD-20
const BTN_SMALL_RAD     = 16
const WIN_TITLE         = "Timer"

const FONT_PATH             = "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf"
const REM_TIME_FONT_SIZE    = 36
const TOOLTIP_FONT_SIZE     = 18
const CONFWIN_FONT_SIZE     = 14

var COLOR_BG            = sdl.Color{R:   4, G:  13, B:  35, A: 255}
var COLOR_FG            = sdl.Color{R: 247, G: 247, B: 255, A: 255}
var COLOR_CLOCK_BG      = sdl.Color{R:  87, G: 115, B: 153, A: 255}
var COLOR_CLOCK_FG      = sdl.Color{R: 254, G:  95, B:  85, A: 255}
var COLOR_CLOCK_FG_DONE = sdl.Color{R:  86, G: 227, B: 159, A: 255}
var COLOR_BTN           = sdl.Color{R:  37, G:  49, B:  65, A: 255}
var COLOR_BTN_HOVER     = sdl.Color{R:  52, G:  69, B:  91, A: 255}
var COLOR_BTN_HOVER_BD  = sdl.Color{R: 244, G: 185, B:  66, A: 255}
var COLOR_TOOLTIP_BG    = sdl.Color{R:  82, G: 108, B: 142, A: 200}

var COLOR_TRANSPARENT   = sdl.Color{R: 0, G: 0, B: 0, A: 0}
