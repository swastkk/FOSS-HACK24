package main

import (
    "encoding/base64"
    "fmt"
    "image"
    "image/png"
    "os"
    "strings"

    "github.com/nfnt/resize"
    "golang.org/x/term"
)

func loadImage(filePath string) (image.Image, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        return nil, err
    }
    return img, nil
}

func resizeImage(img image.Image, width, height uint) image.Image {
    return resize.Resize(width, height, img, resize.Lanczos3)
}

func imageToBase64(img image.Image) (string, error) {
    buf := new(strings.Builder)
    err := png.Encode(buf, img)
    if err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString([]byte(buf.String())), nil
}

func printImageToKitty(imgBase64 string, width, height int) {
    fmt.Printf("\x1b_Gf=1,t=%d,%d;x=%s\x1b\\", width, height, imgBase64)
}

func main() {
    // Define the image paths and terminal layout
    imagePaths := []string{"A_Simple_Podcast.png", "arch_girl.png", "1331008.png", "anime_wall2.png", "anime_wall3.png", "1330654.png"}

    // Get terminal size
    width, height, err := term.GetSize(0)
    if err != nil {
        fmt.Println("Error getting terminal size:", err)
        return
    }
    
    // Terminal layout
    cols, rows := 3, 2
    imgWidth, imgHeight := width/cols, height/rows

    // Print images in a 3x2 grid
    for i, path := range imagePaths {
        img, err := loadImage(path)
        if err != nil {
            fmt.Println("Error loading image:", err)
            return
        }

        resizedImg := resizeImage(img, uint(imgWidth), uint(imgHeight))
        imgBase64, err := imageToBase64(resizedImg)
        if err != nil {
            fmt.Println("Error encoding image to base64:", err)
            return
        }

        printImageToKitty(imgBase64, imgWidth, imgHeight)

        if (i+1)%cols == 0 {
            fmt.Println() // New line after each row of images
        }
    }
}
