package iwidget

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/ttf"
)

type IWidget interface {
    UpdateMouseState(x, y int32, mouseBtnState uint32, frameTime float32)
    Draw(rend *sdl.Renderer)
    DrawTooltip(rend *sdl.Renderer, font *ttf.Font)
}
