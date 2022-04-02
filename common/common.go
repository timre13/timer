package common

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/ttf"
    "github.com/veandco/go-sdl2/img"
    "math"
    "fmt"
    "runtime/debug"
)

func PANIC_ERR(err error) {
    if err != nil {
        sdl.Quit()
        panic(err)
    }
}

func WARN_ERR(err error) {
    if err != nil {
        fmt.Print("Warning: ", err)
        fmt.Println("\n----- Stack trace -----")
        debug.PrintStack()
        fmt.Println("--------- end ---------")
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
    PANIC_ERR(err)
    textTex, err := rend.CreateTextureFromSurface(textSurf)
    PANIC_ERR(err)

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

type Image struct {
    Img *sdl.Texture
    Width, Height int32
}

/*
 * Returns:
 *     Texture
 *     Width
 *     Height
*/
func LoadImage(rend *sdl.Renderer, path string) Image {
    imgSurf, err := img.Load(path)
    PANIC_ERR(err)

    imgTex, err := rend.CreateTextureFromSurface(imgSurf)
    PANIC_ERR(err)

    width, height := imgSurf.W, imgSurf.H
    imgSurf.Free()

    return Image{Img: imgTex, Width: width, Height: height}
}
