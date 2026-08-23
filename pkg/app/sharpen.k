package main

// Amount blends between the original pixel (0) and the fully sharpened
// pixel (1). Values above ~0.8 start to ring on text.
var Amount float32

// Fragment applies a halo-limited five-tap unsharp mask in source-pixel
// space. The result is clamped to the local luminance range so edges gain
// contrast without overshoot halos.
func Fragment(dstPos vec4, srcPos vec2, mask vec4) vec4 {
	size := imageSrc0Size()
	x := srcPos.x
	y := srcPos.y
	c := imageSrc0At(srcPos)
	l := imageSrc0At(vec2(max(x-1, 0), y))
	r := imageSrc0At(vec2(min(x+1, size.x-1), y))
	u := imageSrc0At(vec2(x, max(y-1, 0)))
	d := imageSrc0At(vec2(x, min(y+1, size.y-1)))
	sharp := c*5.0 - (l + r + u + d)
	lo := min(c, min(l, min(r, min(u, d))))
	hi := max(c, max(l, max(r, max(u, d))))
	sharp = min(max(sharp, lo), hi)
	return vec4(
		c.x + Amount*(sharp.x-c.x),
		c.y + Amount*(sharp.y-c.y),
		c.z + Amount*(sharp.z-c.z),
		c.w,
	)
}
