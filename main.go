package main

import (
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	_ "github.com/spakin/netpbm"
)

var img image.Image

func main() {
	// read command line arguments
	if len(os.Args) != 2 {
		log.Fatal("Usage: ppm <image.ppm>")
	}
	fileName := os.Args[1]

	var err error
	img, err = loadImage(fileName)
	if err != nil {
		log.Fatal(err)
	}

	x, y := img.Bounds().Size().X, img.Bounds().Size().Y
	go func() {
		w := new(app.Window)
		w.Option(app.Title("PPM Viewer"))
		w.Option(app.Size(unit.Dp(x), unit.Dp(y)))
		w.Option(app.MaxSize(unit.Dp(x), unit.Dp(y)))
		w.Option(app.MinSize(unit.Dp(x), unit.Dp(y)))
		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	//theme := material.NewTheme()
	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			// Draw the image.
			drawImage(gtx.Ops, img)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}

func drawImage(ops *op.Ops, img image.Image) {
	imageOp := paint.NewImageOp(img)
	imageOp.Filter = paint.FilterNearest
	imageOp.Add(ops)
	paint.PaintOp{}.Add(ops)
}

func loadImage(path string) (image.Image, error) {
	// Open the file.
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	// Decode ppm image.
	img, _, err := image.Decode(file)
	return img, err
}
