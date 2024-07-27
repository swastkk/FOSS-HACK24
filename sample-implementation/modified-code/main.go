package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"syscall"
    "path/filepath"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "icat [directory]",
	Short: "Kitten to display images in a grid layout in Kitty terminal",
	Run:   session,
}

// Will contain all the window parameters 
type windowParameters struct {
	Row    uint16
	Col    uint16
	xPixel uint16
	yPixel uint16
}

var (
	recursive bool
	maxImages int
    globalWindowParameters windowParameters // Contains Global Level Window Parameters 
)

func init() {
	rootCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Scan directory recursively")
	rootCmd.Flags().IntVarP(&maxImages, "max-images", "n", 100, "Maximum number of images to display")
}

// Gets the window size and modifies the globalWindowParameters (global struct) 
func getWindowSize(window windowParameters) err {
	var err error
	var f *os.File

	// Read the window size from device drivers and print them
	if f, err = os.OpenFile("/dev/tty", unix.O_NOCTTY|unix.O_CLOEXEC|unix.O_NDELAY|unix.O_RDWR, 0666); err == nil {
		var sz *unix.Winsize
		if sz, err = unix.IoctlGetWinsize(int(f.Fd()), unix.TIOCGWINSZ); err == nil {
			fmt.Printf("rows: %v columns: %v width: %v height %v\n", sz.Row, sz.Col, sz.Xpixel, sz.Ypixel)
			return nil 
		}
	}

	fmt.Fprintln(os.Stderr, err)
	// os.Exit(1)
    
    return err 
}

// Function handler for changes in window sizes (will be added to goroutines)
func handleWindowSizeChange() {
	err := getWindowSize(globalWindowParameters)
	if err != nil {
		fmt.Println("Error getting window size:", err)
	}
}

// Checks if a given file is an image 
func isImage(fileName string) bool {
	extensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}
	ext := strings.ToLower(filepath.Ext(fileName))
	for _, e := range extensions {
		if ext == e {
			return true
		}
	}
	return false
}

// findImages recursively searches for image files in the given directory
func discoverImages(dir string) ([]string, error) {
	var images []string

	err := filepath.Walk(dir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isImage(info.Name()) {
			images = append(images, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return images, nil
}

// Routine for session - kitten will run in this space
func session(cmd *cobra.Command, args []string) {

	// Check for Arguements
	if len(args) == 0 {
		fmt.Println("Please specify a directory")
		os.Exit(1)
	}

	// Get directory name and discover images
	dir := args[0]
	images, err := discoverImages(dir)
	if err != nil {
		fmt.Printf("Error discovering images: %v\n", err)
		os.Exit(1)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGWINCH)

	// Get the window size initially when kitten is spawned
	handleWindowSizeChange()

	// Goroutine to listen for window size changes
	go func() {
		for {
			sig := <-sigs
			
            // if window size change syscall is detected, execute the handleWindowSizeChange()
			if sig == syscall.SIGWINCH {
				handleWindowSizeChange()
			}
		}
	}()

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
