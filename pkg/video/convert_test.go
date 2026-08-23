package video

import (
	"bytes"
	"image"
	"image/draw"
	"testing"
)

// TestYCbCrToRGBAParallelMatchesSequential verifies that the banded parallel
// conversion produces pixel-identical output to a plain sequential draw.Draw
// pass. Odd dimensions exercise partial chroma rows on the last bands.
func TestYCbCrToRGBAParallelMatchesSequential(t *testing.T) {
	for _, ratio := range []image.YCbCrSubsampleRatio{
		image.YCbCrSubsampleRatio420,
		image.YCbCrSubsampleRatio422,
		image.YCbCrSubsampleRatio444,
	} {
		ratio := ratio
		t.Run(ratio.String(), func(t *testing.T) {
			w, h := 97, 61
			src := image.NewYCbCr(image.Rect(0, 0, w, h), ratio)
			for i := range src.Y {
				src.Y[i] = byte(i*7 + 1)
			}
			for i := range src.Cb {
				src.Cb[i] = byte(i*13 + 3)
			}
			for i := range src.Cr {
				src.Cr[i] = byte(i*29 + 11)
			}

			want := image.NewRGBA(image.Rect(0, 0, w, h))
			draw.Draw(want, want.Bounds(), src, src.Bounds().Min, draw.Src)
			got := ycbcrToRGBA(src)

			if got.Bounds() != want.Bounds() {
				t.Fatalf("bounds mismatch: got %v want %v", got.Bounds(), want.Bounds())
			}
			if !bytes.Equal(got.Pix, want.Pix) {
				t.Fatal("parallel conversion output differs from sequential reference")
			}
		})
	}
}
