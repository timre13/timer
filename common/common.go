package common

import (
    "math"
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/ttf"
)

func CHECK_ERR(err error) {
    if err != nil {
        sdl.Quit()
        panic(err)
    }
}

func lerpInt(a int, b int, t float32) int {
    return int(math.Round(float64(float32(a) + t * float32(b - a))))
}

func LerpColors(x *sdl.Color, y *sdl.Color, t float32) sdl.Color {
    return sdl.Color{
        R: uint8(lerpInt(int(x.R), int(y.R), t)),
        G: uint8(lerpInt(int(x.G), int(y.G), t)),
        B: uint8(lerpInt(int(x.B), int(y.B), t)),
        A: 255,
    }
}

func RenderText(rend *sdl.Renderer, font *ttf.Font, str string, color *sdl.Color, x, y int32, areCoordsCent bool) {
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

