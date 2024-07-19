package main

import "math"

type GridLayout struct {
    Columns    int
    Rows       int
    CellWidth  int
    CellHeight int
}

func calculateGridLayout(imageCount, width, height int) GridLayout {
    aspectRatio := float64(width) / float64(height)
    columns := int(math.Sqrt(float64(imageCount) * aspectRatio))
    rows := int(math.Ceil(float64(imageCount) / float64(columns)))

    // Adjust if too many columns
    if columns > 5 {
        columns = 5
        rows = int(math.Ceil(float64(imageCount) / 5))
    }

    cellWidth := width / columns
    cellHeight := height / rows

    return GridLayout{
        Columns:    columns,
        Rows:       rows,
        CellWidth:  cellWidth,
        CellHeight: cellHeight,
    }
}