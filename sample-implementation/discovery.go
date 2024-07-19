package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

type ImageInfo struct {
    Path   string
    Width  int
    Height int
}

func discoverImages(root string, recursive bool) ([]ImageInfo, error) {
    var images []ImageInfo

    walkFn := func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && isImageFile(path) {
            imgInfo, err := getImageInfo(path)
            if err == nil {
                images = append(images, imgInfo)
            }
        }
        if !recursive && info.IsDir() && path != root {
            return filepath.SkipDir
        }
        return nil
    }

    err := filepath.Walk(root, walkFn)
    return images, err
}

func isImageFile(path string) bool {
    ext := filepath.Ext(path)
    return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

func getImageInfo(path string) (ImageInfo, error) {
    file, err := os.Open(path)
    if err != nil {
        return ImageInfo{}, err
    }
    defer file.Close()

    img, _, err := image.DecodeConfig(file)
    if err != nil {
        return ImageInfo{}, err
    }

    return ImageInfo{
        Path:   path,
        Width:  img.Width,
        Height: img.Height,
    }, nil
}