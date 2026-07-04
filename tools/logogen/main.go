// logogen recreates the "Technology with Purpose" spiral mark
// programmatically (concentric broken arcs, 3 purple tones), renders a
// preview PNG for visual comparison, and emits the mark as terminal
// half-block ANSI art embedded in Go source.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strings"
)

// ---- palette (sampled from the original) ----
var (
	lavender = color.NRGBA{0xD9, 0xB3, 0xFF, 0xFF}
	violet   = color.NRGBA{0xA3, 0x2C, 0xF2, 0xFF}
	deep     = color.NRGBA{0x6E, 0x1F, 0xA8, 0xFF}
)

// arc is a stroked circle segment. Angles in degrees, 0° = 3 o'clock,
// counter-clockwise positive, drawn with round caps.
type arc struct {
	r     float64 // radius (in supersampled px)
	w     float64 // stroke half-width
	a0    float64 // start angle
	a1    float64 // end angle (a1 > a0, may exceed 360)
	color color.NRGBA
}

const (
	super = 208 // supersampled canvas (8x)
	cell  = 26  // final pixel width/height => 26 cols x 13 rows
	scale = super / cell
)

// The mark, outside in. Tuned by eye against the original.
func arcs() []arc {
	c := float64(super) / 2
	_ = c
	u := float64(super) / 208 // unit
	return []arc{
		// outside accents: the pin tail (bottom-left) and a lavender tick (bottom-right)
		{r: 96 * u, w: 6 * u, a0: 222, a1: 256, color: violet},
		{r: 96 * u, w: 5.5 * u, a0: 288, a1: 306, color: lavender},
		// outer ring: lavender on the right, violet sweeping top->left->bottom
		{r: 82 * u, w: 6 * u, a0: -50, a1: 45, color: lavender},
		{r: 82 * u, w: 7 * u, a0: 62, a1: 255, color: violet},
		// ring 2: deep purple top-left, violet bottom-right
		{r: 62 * u, w: 6.5 * u, a0: 70, a1: 210, color: deep},
		{r: 62 * u, w: 6 * u, a0: 228, a1: 325, color: violet},
		// ring 3: violet top, deep bottom-left
		{r: 43 * u, w: 6 * u, a0: 35, a1: 155, color: violet},
		{r: 43 * u, w: 6 * u, a0: 180, a1: 280, color: deep},
		// center: lavender "C" open to the right
		{r: 22 * u, w: 6 * u, a0: 30, a1: 320, color: lavender},
	}
}

func draw() *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, super, super))
	cx, cy := float64(super)/2, float64(super)/2
	for _, a := range arcs() {
		drawArc(img, cx, cy, a)
	}
	return img
}

func drawArc(img *image.NRGBA, cx, cy float64, a arc) {
	// Bounding box of the arc's ring.
	rOut := a.r + a.w
	x0, x1 := int(cx-rOut)-1, int(cx+rOut)+1
	y0, y1 := int(cy-rOut)-1, int(cy+rOut)+1
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if x < 0 || y < 0 || x >= super || y >= super {
				continue
			}
			dx, dy := float64(x)-cx, cy-float64(y) // math coords (y up)
			dist := math.Hypot(dx, dy)
			if math.Abs(dist-a.r) > a.w {
				if !inCap(dx, dy, a) {
					continue
				}
				img.SetNRGBA(x, y, a.color)
				continue
			}
			ang := math.Atan2(dy, dx) * 180 / math.Pi // -180..180
			if angleIn(ang, a.a0, a.a1) || inCap(dx, dy, a) {
				img.SetNRGBA(x, y, a.color)
			}
		}
	}
}

// inCap checks the round end-caps at both arc ends.
func inCap(dx, dy float64, a arc) bool {
	for _, deg := range []float64{a.a0, a.a1} {
		rad := deg * math.Pi / 180
		ex, ey := a.r*math.Cos(rad), a.r*math.Sin(rad)
		if math.Hypot(dx-ex, dy-ey) <= a.w {
			return true
		}
	}
	return false
}

func angleIn(ang, a0, a1 float64) bool {
	for ang < a0 {
		ang += 360
	}
	return ang <= a1
}

// downsample box-averages the supersampled canvas to cell x cell,
// carrying alpha so the terminal background shows through.
func downsample(src *image.NRGBA) *image.NRGBA {
	out := image.NewNRGBA(image.Rect(0, 0, cell, cell))
	for cy := 0; cy < cell; cy++ {
		for cx := 0; cx < cell; cx++ {
			var r, g, b, a, n float64
			for y := cy * scale; y < (cy+1)*scale; y++ {
				for x := cx * scale; x < (cx+1)*scale; x++ {
					px := src.NRGBAAt(x, y)
					w := float64(px.A) / 255
					r += float64(px.R) * w
					g += float64(px.G) * w
					b += float64(px.B) * w
					a += w
					n++
				}
			}
			if a == 0 {
				continue
			}
			out.SetNRGBA(cx, cy, color.NRGBA{
				uint8(r / a), uint8(g / a), uint8(b / a), uint8(255 * a / n),
			})
		}
	}
	return out
}

// emitANSI renders the small image as half-block art: each cell is two
// vertical pixels — fg paints the top (▀), bg paints the bottom.
func emitANSI(img *image.NRGBA) string {
	const alphaMin = 60 // below this the pixel is "transparent"
	var b strings.Builder
	for y := 0; y < cell; y += 2 {
		for x := 0; x < cell; x++ {
			top := img.NRGBAAt(x, y)
			bot := img.NRGBAAt(x, y+1)
			topOK, botOK := top.A >= alphaMin, bot.A >= alphaMin
			switch {
			case topOK && botOK:
				fmt.Fprintf(&b, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀\x1b[0m",
					top.R, top.G, top.B, bot.R, bot.G, bot.B)
			case topOK:
				fmt.Fprintf(&b, "\x1b[38;2;%d;%d;%dm▀\x1b[0m", top.R, top.G, top.B)
			case botOK:
				fmt.Fprintf(&b, "\x1b[38;2;%d;%d;%dm▄\x1b[0m", bot.R, bot.G, bot.B)
			default:
				b.WriteString(" ")
			}
		}
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func main() {
	big := draw()
	small := downsample(big)

	// 1. preview PNG (scaled 8x back up, nearest-neighbor, so it's visible)
	preview := image.NewNRGBA(image.Rect(0, 0, cell*8, cell*8))
	for y := 0; y < cell*8; y++ {
		for x := 0; x < cell*8; x++ {
			preview.SetNRGBA(x, y, small.NRGBAAt(x/8, y/8))
		}
	}
	fp, _ := os.Create("preview_small.png")
	png.Encode(fp, preview)
	fp.Close()
	fb, _ := os.Create("preview_big.png")
	png.Encode(fb, big)
	fb.Close()

	// 2. ANSI to stdout (visible in a real terminal)
	fmt.Println(emitANSI(small))

	// 3. Go source with the art, ready to embed
	art := emitANSI(small)
	src := "// Code generated by logogen (scratchpad); see the source image in the repo history. DO NOT EDIT BY HAND.\n\n" +
		"package tui\n\n" +
		"// logoMark is the QASON spiral mark (half-block truecolor ANSI),\n" +
		"// derived from the Technology with Purpose Foundation isotype.\n" +
		"// " + fmt.Sprintf("%d columns x %d rows.", cell, (cell+1)/2) + "\n" +
		"const logoMark = " + fmt.Sprintf("%q", art) + "\n"
	os.WriteFile("logomark_generated.go", []byte(src), 0o644)
	fmt.Fprintln(os.Stderr, "wrote preview_small.png, preview_big.png, logomark_generated.go")
}
