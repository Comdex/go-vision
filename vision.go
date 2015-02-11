/*
	Package vision provides some basic operation on computer vision.
*/
package vision

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	//	"errors"
	//	"path/filepath"

	//	"reflect"
)

type Size struct {
	Width  int
	Height int
}

func (sz Size) Area() int {
	return sz.Width * sz.Height
}

type GrayImage struct {
	Size
	Pixels []byte
}

func (m *GrayImage) String() string {
	return fmt.Sprintf("[gray-image]%dx%d", m.Width, m.Height)
}

func ImageFromFile(fn string) (image.Image, error) {
	// Open the file.
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the image.
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	//fmt.Println("type:", reflect.TypeOf(img))

	return img, nil
}

func (m *GrayImage) SetImage(img image.Image) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	m.Width, m.Height = w, h
	l := m.Area()
	if len(m.Pixels) < l {
		m.Pixels = make([]byte, l)
	} else if len(m.Pixels) > l {
		m.Pixels = m.Pixels[:l]
	}

	switch img := img.(type) {
	case *image.Gray:
		if img.Stride == w {
			copy(m.Pixels, img.Pix)
		} else {
			for y := 0; y < h; y++ {
				copy(m.Pixels[y*w:], img.Pix[y*img.Stride:y*img.Stride+w])
			}
		}

		//	case *image.RGBA:
		//		_ = m
		//		return nil, errors.New("TODO RGBA")

	case *image.NRGBA:
		idx := 0
		for y := 0; y < h; y++ {
			s := y * img.Stride
			for x := 0; x < w; x++ {
				r, g, b, a := img.Pix[s], img.Pix[s+1], img.Pix[s+2], img.Pix[s+3]
				switch a {
				case 0:
					m.Pixels[idx] = 0
				case 0xff:
					m.Pixels[idx] = byte((19595*uint32(r) + 38470*uint32(g) + 7471*uint32(b) + 1<<15) >> 16)
				default:
					m.Pixels[idx] = byte(((19595*uint32(r)+38470*uint32(g)+7471*uint32(b))*uint32(a) + (1<<15)*255) / ((1 << 16) * 255))
				}
				idx++
				s += 4
			}
		}

	case *image.Paletted:
		pal := make([]byte, len(img.Palette))
		for i, c := range img.Palette {
			r, g, b, _ := c.RGBA()
			pal[i] = byte((19595*uint32(r) + 38470*uint32(g) + 7471*uint32(b) + 1<<15) >> 24)
		}
		idx := 0
		for y := 0; y < h; y++ {
			l := img.Pix[y*img.Stride:]
			for x := 0; x < w; x++ {
				m.Pixels[idx] = pal[l[x]]
				idx++
			}
		}

	case *image.YCbCr:
		if img.YStride == w {
			copy(m.Pixels, img.Y)
		} else {
			for y := 0; y < h; y++ {
				copy(m.Pixels[y*w:], img.Y[y*img.YStride:y*img.YStride+w])
			}
		}

	default:
		idx := 0
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				c := img.At(x, y)
				r, g, b, _ := c.RGBA()
				y := (19595*r + 38470*g + 7471*b + 1<<15) >> 24
				m.Pixels[idx] = byte(y)
				idx++
			}
		}
	}
}

func (m *GrayImage) AsImage() image.Image {
	gray := image.NewGray(image.Rect(0, 0, m.Width, m.Height))
	if gray.Stride == m.Width {
		copy(gray.Pix, m.Pixels)
	} else {
		for y := 0; y < m.Height; y++ {
			copy(gray.Pix[y*gray.Stride:y*gray.Stride+m.Width], m.Pixels[y*m.Width:])
		}
	}
	return gray
}

func SaveImageAsJpeg(m image.Image, quality int) error {
	return nil
}

func SaveImageAsPng(m image.Image, fn string) error {
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, m)
}

type RGB [3]byte

type RGBInt [3]int

type RGBImage struct {
	Size
	Pixels []RGB
}

func (m *RGBImage) SetImage(img image.Image) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	m.Width, m.Height = w, h
	l := m.Area()
	if len(m.Pixels) < l {
		m.Pixels = make([]RGB, l)
	} else if len(m.Pixels) > l {
		m.Pixels = m.Pixels[:l]
	}

	switch img := img.(type) {
	case *image.Gray:
		idx := 0
		for y := 0; y < h; y++ {
			s := y * img.Stride
			for x := 0; x < w; x++ {
				c := img.Pix[s]
				m.Pixels[idx][0], m.Pixels[idx][1], m.Pixels[idx][2] = c, c, c
				idx++
				s++
			}
		}

	case *image.NRGBA:
		idx := 0
		for y := 0; y < h; y++ {
			s := y * img.Stride
			for x := 0; x < w; x++ {
				r, g, b, a := img.Pix[s], img.Pix[s+1], img.Pix[s+2], img.Pix[s+3]
				switch a {
				case 0:
					m.Pixels[idx][0], m.Pixels[idx][1], m.Pixels[idx][2] = 0, 0, 0
				case 0xff:
					m.Pixels[idx][0], m.Pixels[idx][1], m.Pixels[idx][2] = r, g, b
				default:
					m.Pixels[idx][0], m.Pixels[idx][1], m.Pixels[idx][2] =
						byte(uint32(r)*uint32(a)/255), byte(uint32(g)*uint32(a)/255), byte(uint32(b)*uint32(a)/255)
				}
				idx++
				s += 4
			}
		}

	case *image.Paletted:
		pal := make([]RGB, len(img.Palette))
		for i, c := range img.Palette {
			r, g, b, _ := c.RGBA()
			pal[i][0], pal[i][1], pal[i][2] = byte(r>>8), byte(g>>8), byte(b>>8)
		}
		idx := 0
		for y := 0; y < h; y++ {
			s := y * img.Stride
			for x := 0; x < w; x++ {
				m.Pixels[idx] = pal[img.Pix[s]]
				idx++
				s++
			}
		}
		/*
			case *image.YCbCr:
				if img.YStride == w {
					copy(m.Pixels, img.Y)
				} else {
					for y := 0; y < h; y++ {
						copy(m.Pixels[y*w:], img.Y[y*img.YStride:y*img.YStride+w])
					}
				}
		*/
	default:
		idx := 0
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				c := img.At(x, y)
				r, g, b, _ := c.RGBA()
				m.Pixels[idx][0], m.Pixels[idx][1], m.Pixels[idx][2] = byte(r>>8), byte(g>>8), byte(b>>8)
				idx++
			}
		}
	}
}

func (m RGBImage) AsImage() image.Image {
	rgb := image.NewRGBA(image.Rect(0, 0, m.Width, m.Height))
	idx := 0
	for y := 0; y < m.Height; y++ {
		p := y * rgb.Stride
		for x := 0; x < m.Width; x++ {
			rgb.Pix[p], rgb.Pix[p + 1], rgb.Pix[p + 2], rgb.Pix[p + 3] = m.Pixels[idx][0], m.Pixels[idx][1], m.Pixels[idx][2], 255
			idx++
			p += 4
		}
	}
	return rgb
}
