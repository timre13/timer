package common

import (
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/ttf"
    "github.com/veandco/go-sdl2/img"
    "math"
    "fmt"
    "runtime/debug"
    "os"
    "path/filepath"
    "strings"
)

type SessionType int
const (
    SESSION_TYPE_WORK  SessionType = iota
    SESSION_TYPE_BREAK SessionType = iota
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

func GetExeDir() string {
    path, err := os.Executable()
    PANIC_ERR(err)

    dirPath := filepath.Dir(path)
    return dirPath
}

func GetRealPath(exeDir, path string) string {
    if filepath.IsAbs(path) { // Absolute path
        return path
    } else { // Relative path
        return filepath.Join(exeDir, path)
    }
}

func MinsToMillisecs(mins int) int {
    return mins*60*1000
}

func lerpInt(a int, b int, t float32) int {
    return int(math.Round(float64(float32(a) + t * float32(b - a))))
}

func LerpColors(x *sdl.Color, y *sdl.Color, t float32) sdl.Color {
    return sdl.Color{
        R: uint8(lerpInt(int(x.R), int(y.R), t)),
        G: uint8(lerpInt(int(x.G), int(y.G), t)),
        B: uint8(lerpInt(int(x.B), int(y.B), t)),
        A: uint8(lerpInt(int(x.A), int(y.A), t)),
    }
}

func RenderText(rend *sdl.Renderer, font *ttf.Font, str string, color *sdl.Color, x, y int32, isXCent bool, isYCent bool) {
    if str == "" {
        return
    }

    lines := strings.Split(str, "\n")
    var yOffs int32
    for _, line := range lines {
        textSurf, err := font.RenderUTF8Blended(line, *color)
        PANIC_ERR(err)
        textTex, err := rend.CreateTextureFromSurface(textSurf)
        PANIC_ERR(err)

        var rectX, rectY int32
        if isXCent {
            rectX = x-textSurf.W/2
        } else {
            rectX = x
        }
        if isYCent {
            rectY = y-textSurf.H/2
        } else {
            rectY = y
        }

        dstRect := sdl.Rect{
            X: rectX, Y: rectY+yOffs,
            W: textSurf.W, H: textSurf.H}
        rend.Copy(textTex, nil, &dstRect)

        yOffs += textSurf.H

        textSurf.Free()
        textTex.Destroy()
    }
}

func GetTextSize(font *ttf.Font, str string) (int32, int32) {
    lines := strings.Split(str, "\n")
    var outW, outH int
    for _, line := range lines {
        w, h, err := font.SizeUTF8(line)
        PANIC_ERR(err)

        outH += h
        if w > outW { outW = w }
    }
    return int32(outW), int32(outH)
}

func GetTextWidth(font *ttf.Font, text string) int32 {
    w, _ := GetTextSize(font, text)
    return int32(w)
}

func RenderHorizText(rend *sdl.Renderer, font *ttf.Font, str string, color *sdl.Color, x, y int32) {
    if str == "" {
        return
    }

    textSurf, err := font.RenderUTF8Blended(str, *color)
    PANIC_ERR(err)
    textTex, err := rend.CreateTextureFromSurface(textSurf)
    PANIC_ERR(err)

    dstRect := sdl.Rect{
        X: x, Y: y,
        W: textSurf.W, H: textSurf.H}
    rend.CopyEx(textTex, nil, &dstRect, -90, nil, sdl.FLIP_NONE)

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
    fmt.Printf("Loading image: \"%s\"\n", path)
    imgSurf, err := img.Load(path)
    PANIC_ERR(err)

    imgTex, err := rend.CreateTextureFromSurface(imgSurf)
    PANIC_ERR(err)

    width, height := imgSurf.W, imgSurf.H
    imgSurf.Free()

    return Image{Img: imgTex, Width: width, Height: height}
}
