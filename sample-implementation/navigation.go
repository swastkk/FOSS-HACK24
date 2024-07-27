package main

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

func handleNavigation(images []ImageInfo, layout GridLayout) error {
    if err := keyboard.Open(); err != nil {
        return err
    }
    defer keyboard.Close()

    currentIndex := 0
    for {
        highlightImage(images[currentIndex], layout)

        _, key, err := keyboard.GetKey()
        if err != nil {
            return err
        }

        switch key {
        case keyboard.KeyArrowDown:
            currentIndex = (currentIndex + layout.Columns) % len(images)
        case keyboard.KeyArrowUp:
            currentIndex = (currentIndex - layout.Columns + len(images)) % len(images)
        case keyboard.KeyArrowLeft:
            currentIndex = (currentIndex - 1 + len(images)) % len(images)
        case keyboard.KeyArrowRight:
            currentIndex = (currentIndex + 1) % len(images)
        case keyboard.KeyEsc:
            return nil
        }
    }
}

func highlightImage(img ImageInfo, layout GridLayout) {
    // Implement highlighting logic here
    // This could involve re-rendering the image with a border or overlay
    fmt.Printf("Highlighted: %s\n", img.Path)
}
