package main

import (
	"fmt"
	"os/exec"
)

func renderImageGrid(images []ImageInfo, layout GridLayout) error {
    for i, img := range images {
        row := i / layout.Columns
        col := i % layout.Columns
        x := col * layout.CellWidth
        y := row * layout.CellHeight

        err := renderImage(img, x, y, layout.CellWidth, layout.CellHeight)
        if err != nil {
            return err
        }
    }
    return nil
}

func renderImage(img ImageInfo, x, y, width, height int) error {
    cmd := exec.Command("kitty", "+kitten", "icat",
        "--place", fmt.Sprintf("%dx%d@%dx%d", width, height, x, y),
        img.Path)
    return cmd.Run()
}