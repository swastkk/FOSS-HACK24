package utils

import (
	"image"
	"sync"

	"github.com/disintegration/imaging"
)

type ImageCache struct {
    cache map[string]image.Image
    mutex sync.RWMutex
}

func NewImageCache() *ImageCache {
    return &ImageCache{
        cache: make(map[string]image.Image),
    }
}

func (c *ImageCache) Get(key string) (image.Image, bool) {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    img, ok := c.cache[key]
    return img, ok
}

func (c *ImageCache) Set(key string, img image.Image) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.cache[key] = img
}

func ResizeImage(img image.Image, width, height int) image.Image {
    return imaging.Fit(img, width, height, imaging.Lanczos)
}