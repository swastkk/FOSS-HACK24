package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "icat [directory]",
    Short: "Display images in a grid layout in Kitty terminal",
    Run:   run,
}

var (
    recursive bool
    maxImages int
)

func init() {
    rootCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Scan directory recursively")
    rootCmd.Flags().IntVarP(&maxImages, "max-images", "n", 100, "Maximum number of images to display")
}

func run(cmd *cobra.Command, args []string) {
    if len(args) == 0 {
        fmt.Println("Please specify a directory")
        os.Exit(1)
    }

    dir := args[0]
    images, err := discoverImages(dir, recursive)
    if err != nil {
        fmt.Printf("Error discovering images: %v\n", err)
        os.Exit(1)
    }

    if len(images) > maxImages {
        images = images[:maxImages]
    }

    width, height, err := getWindowSize()
    if err != nil {
        fmt.Printf("Error getting window size: %v\n", err)
        os.Exit(1)
    }

    layout := calculateGridLayout(len(images), width, height)

    err = renderImageGrid(images, layout)
    if err != nil {
        fmt.Printf("Error rendering image grid: %v\n", err)
        os.Exit(1)
    }

    err = handleNavigation(images, layout)
    if err != nil {
        fmt.Printf("Error during navigation: %v\n", err)
        os.Exit(1)
    }
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}