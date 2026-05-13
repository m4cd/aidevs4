// Package maze extracts and solves pipe-rotation puzzles from images.
//
// A puzzle is a 3x3 grid where each cell contains a pipe segment with
// connectors on some subset of its four sides. The "solved" state of the
// puzzle is given as a second image; the task is to determine how many
// clockwise 90° rotations each cell of the unsolved puzzle requires to
// match its counterpart in the solved puzzle.
package maze

import (
	"fmt"
	"image"
	"image/color"
	"os"

	_ "image/jpeg"
	_ "image/png"
)

// darkThreshold: pixels with luminance below this are considered "wall" pixels.
const darkThreshold = 80

// Sides records which of a cell's four edges have pipe connectors.
type Sides struct {
	Top, Right, Bottom, Left bool
}

// String renders as 4 characters in TRBL order; '-' means not connected.
// E.g. "T-BL" means top, bottom, and left are connected; right is not.
func (s Sides) String() string {
	c := func(b bool, ch byte) byte {
		if b {
			return ch
		}
		return '-'
	}
	return string([]byte{c(s.Top, 'T'), c(s.Right, 'R'), c(s.Bottom, 'B'), c(s.Left, 'L')})
}

// RotateCW returns the sides after one clockwise 90° rotation:
// top moves to right, right to bottom, bottom to left, left to top.
func (s Sides) RotateCW() Sides {
	return Sides{Top: s.Left, Right: s.Top, Bottom: s.Right, Left: s.Bottom}
}

// FindRotations returns k in {0,1,2,3} such that rotating `from` k times
// clockwise equals `to`. Returns -1 if no rotation matches.
func FindRotations(from, to Sides) int {
	s := from
	for k := 0; k < 4; k++ {
		if s == to {
			return k
		}
		s = s.RotateCW()
	}
	return -1
}

// ExtractCells loads an image from path, auto-detects the maze grid as the
// largest dark connected component, and returns a 3x3 matrix of cell
// sub-images. cells[row][col]; [0][0] is top-left.
func ExtractCells(path string) ([3][3]image.Image, error) {
	var cells [3][3]image.Image
	img, err := loadImage(path)
	if err != nil {
		return cells, err
	}
	box, err := findMazeBounds(img)
	if err != nil {
		return cells, err
	}
	w, h := box.Dx(), box.Dy()
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			rect := image.Rect(
				box.Min.X+(w*c)/3,
				box.Min.Y+(h*r)/3,
				box.Min.X+(w*(c+1))/3,
				box.Min.Y+(h*(r+1))/3,
			)
			cells[r][c] = subImage(img, rect)
		}
	}
	return cells, nil
}

// DetectSides examines four small regions inside a cell, one just past the
// grid line at the middle of each edge, and reports which sides have a pipe
// connector reaching them.
//
// The key trick is that sampling AT the edge would pick up the grid line
// itself (which spans the full width of every edge of every cell). Sampling
// just inside the cell, at the middle of each side, only sees pixels of a
// pipe — pipes always cross the grid line at the edge's midpoint.
func DetectSides(cell image.Image) Sides {
	b := cell.Bounds()
	w, h := b.Dx(), b.Dy()
	s := minInt(w, h)

	// Geometry:
	//   inset: depth from the edge to skip (covers the grid line)
	//   depth: thickness of the sampling band, perpendicular to the edge
	//   span:  length of the band along the edge, centered on midpoint
	inset := s * 12 / 100
	depth := s * 12 / 100
	span := s * 25 / 100
	if inset < 2 {
		inset = 2
	}
	if depth < 2 {
		depth = 2
	}
	if span < 4 {
		span = 4
	}

	midX := b.Min.X + w/2
	midY := b.Min.Y + h/2

	top := countDark(cell, image.Rect(midX-span/2, b.Min.Y+inset, midX+span/2, b.Min.Y+inset+depth))
	right := countDark(cell, image.Rect(b.Max.X-inset-depth, midY-span/2, b.Max.X-inset, midY+span/2))
	bottom := countDark(cell, image.Rect(midX-span/2, b.Max.Y-inset-depth, midX+span/2, b.Max.Y-inset))
	left := countDark(cell, image.Rect(b.Min.X+inset, midY-span/2, b.Min.X+inset+depth, midY+span/2))

	// A connected pipe of ~4px thickness crossing a `depth`-deep band leaves
	// roughly 4*depth dark pixels. A disconnected side leaves ~0. Threshold
	// halfway between to be robust to slight pipe-thickness variation.
	threshold := depth * 2
	return Sides{
		Top:    top >= threshold,
		Right:  right >= threshold,
		Bottom: bottom >= threshold,
		Left:   left >= threshold,
	}
}

// SolveRotations computes how many clockwise 90° rotations each cell of the
// unsolved puzzle requires to match the solved puzzle. Returns a 3x3 matrix
// where element [r][c] is in {0,1,2,3}, or -1 if no rotation matches
// (indicates either a detection error or that the inputs don't correspond).
func SolveRotations(unsolvedPath, solvedPath string) ([3][3]int, error) {
	var result [3][3]int

	unsolved, err := ExtractCells(unsolvedPath)
	if err != nil {
		return result, fmt.Errorf("unsolved: %w", err)
	}
	solved, err := ExtractCells(solvedPath)
	if err != nil {
		return result, fmt.Errorf("solved: %w", err)
	}

	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			from := DetectSides(unsolved[r][c])
			to := DetectSides(solved[r][c])
			result[r][c] = FindRotations(from, to)
		}
	}
	return result, nil
}

// --- internals ---

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	return img, nil
}

// findMazeBounds returns the bounding box of the largest 4-connected
// component of dark pixels. For these puzzle images the maze grid is always
// at least an order of magnitude larger than any other dark region (title
// text, factory icons, etc.) so it's the unambiguous winner.
func findMazeBounds(img image.Image) (image.Rectangle, error) {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	dark := make([]bool, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if luma(img.At(b.Min.X+x, b.Min.Y+y)) < darkThreshold {
				dark[y*w+x] = true
			}
		}
	}

	visited := make([]bool, w*h)
	queue := make([]int, 0, 1024)
	var best image.Rectangle
	bestSize := 0

	for start := 0; start < w*h; start++ {
		if !dark[start] || visited[start] {
			continue
		}
		queue = queue[:0]
		queue = append(queue, start)
		visited[start] = true

		minX, minY := start%w, start/w
		maxX, maxY := minX, minY
		size := 0

		for head := 0; head < len(queue); head++ {
			p := queue[head]
			size++
			px, py := p%w, p/w
			if px < minX {
				minX = px
			}
			if py < minY {
				minY = py
			}
			if px > maxX {
				maxX = px
			}
			if py > maxY {
				maxY = py
			}
			// 4-connected neighbors
			if px > 0 {
				if n := p - 1; dark[n] && !visited[n] {
					visited[n] = true
					queue = append(queue, n)
				}
			}
			if px < w-1 {
				if n := p + 1; dark[n] && !visited[n] {
					visited[n] = true
					queue = append(queue, n)
				}
			}
			if py > 0 {
				if n := p - w; dark[n] && !visited[n] {
					visited[n] = true
					queue = append(queue, n)
				}
			}
			if py < h-1 {
				if n := p + w; dark[n] && !visited[n] {
					visited[n] = true
					queue = append(queue, n)
				}
			}
		}

		if size > bestSize {
			bestSize = size
			best = image.Rect(
				b.Min.X+minX, b.Min.Y+minY,
				b.Min.X+maxX+1, b.Min.Y+maxY+1,
			)
		}
	}

	if bestSize == 0 {
		return image.Rectangle{}, fmt.Errorf("no dark pixels found")
	}
	return best, nil
}

func countDark(img image.Image, r image.Rectangle) int {
	r = r.Intersect(img.Bounds())
	n := 0
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			if luma(img.At(x, y)) < darkThreshold {
				n++
			}
		}
	}
	return n
}

func subImage(img image.Image, r image.Rectangle) image.Image {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	if si, ok := img.(subImager); ok {
		return si.SubImage(r)
	}
	out := image.NewRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			out.Set(x-r.Min.X, y-r.Min.Y, img.At(x, y))
		}
	}
	return out
}

func luma(c color.Color) uint8 {
	r, g, b, _ := c.RGBA()
	y := (299*r + 587*g + 114*b) / 1000
	return uint8(y >> 8)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
