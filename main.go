package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/bmp"
)

var (
	T       = 4
	N       = 5
	TimeTag = false
)

type ImgWithID struct {
	image.Image
	ID   int
	Path string
}

func main() {
	flag.IntVar(&T, "t", T, "T")
	flag.IntVar(&N, "n", N, "N")
	flag.BoolVar(&TimeTag, "tag", TimeTag, "time tag on decrypted image")
	flag.Parse()
	if T <= 0 || N <= 0 || T > N {
		fmt.Println("t must be between 1 and n")
		return
	}
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("no input file")
		return
	}
	if len(args) == 1 {
		img, format, err := loadImage(args[0])
		if err != nil {
			panic(err)
		}
		fmt.Printf("loaded %s image from %s\n", format, args[0])
		generateImage(img)
		fmt.Println("done")
		return
	}
	imgs := []ImgWithID{}
	for _, path := range args {
		i, err := strconv.Atoi(path)
		if err != nil {
			panic(err)
		}
		path = fmt.Sprintf("image_%d.bmp", i)
		img, format, err := loadImage(path)
		if err != nil {
			panic(err)
		}
		if format != "bmp" {
			fmt.Printf("%s is not a bmp file\n", path)
			return
		}
		fmt.Printf("loaded %s image from %s\n", format, path)
		imgs = append(imgs, ImgWithID{
			Image: img,
			ID:    i,
			Path:  path,
		})
	}
	decryptImage(imgs...)
	fmt.Println("done")
}

func loadImage(path string) (image.Image, string, error) {
	reader, err := os.Open(path)
	defer reader.Close()
	if err != nil {
		return nil, "", fmt.Errorf("failed to load %s: %v", path, err)
	}
	return image.Decode(reader)
}

func saveImage(img image.Image, path string) {
	buffer := bytes.NewBuffer(nil)
	if err := bmp.Encode(buffer, img); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path, buffer.Bytes(), 0644); err != nil {
		panic(err)
	}
}

// Image

func generateImage(secretImg image.Image) {
	// Decode secret image
	w, h := secretImg.Bounds().Dx(), secretImg.Bounds().Dy()
	originImg := image.NewGray(image.Rect(0, 0, w, h))
	secretImgData := []uint8{}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			oldColor := secretImg.At(x, y)
			grayColor := color.GrayModel.Convert(oldColor)
			originImg.Set(x, y, grayColor)
			valueY := grayColor.(color.Gray).Y + 1
			if valueY >= 250 {
				secretImgData = append(secretImgData, 250)
				valueY -= 250
				valueY += 1
			}
			secretImgData = append(secretImgData, uint8(valueY))
		}
	}
	saveImage(originImg, "image_origin.bmp")
	grayImgData := make([][]uint8, N)
	for i := 0; i < len(secretImgData); i += T {
		result := generateShares(T, N, secretImgData[i:i+T], 251)
		for j := 0; j < N; j++ {
			grayImgData[j] = append(grayImgData[j], result[j])
		}
	}
	for i, data := range grayImgData {
		h := (len(data)-1)/w + 1
		grayImg := image.NewGray(image.Rect(0, 0, w, h))
		for i := 0; i < len(data); i++ {
			grayImg.Set(i%w, i/w, color.Gray{Y: data[i]})
		}
		fileName := fmt.Sprintf("image_%d.bmp", i+1)
		saveImage(grayImg, fileName)
	}
}

func generateShares(t, n int, origin []uint8, p int) (points []uint8) {
	poly := origin[:]

	termAt := func(n int) int {
		ans := 0
		for i, coe := range poly {
			x := pow(n, i, p) * int(coe) % p
			ans = (ans + x) % p
		}
		return ans
	}

	for i := 0; i < n; i++ {
		points = append(points, uint8(termAt(i+1)))
	}
	return
}

func decryptImage(imgs ...ImgWithID) {
	grayImgData := make([][]uint8, len(imgs))
	w, h := imgs[0].Bounds().Dx(), imgs[0].Bounds().Dy()
	for i, img := range imgs {
		if img.Bounds().Dx() != w || img.Bounds().Dy() != h {
			panic(fmt.Errorf("image %s is not the same size as the first image\n", img.Path))
		}
		grayImgData[i] = make([]uint8, w*h)
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				oldColor := img.At(x, y)
				r, g, b, a := oldColor.RGBA()
				if r != g || r != b || a != 0xffff {
					panic(fmt.Errorf("bmp %s is not a gray image", img.Path))
				}
				valueY := uint8(r >> 8)
				grayImgData[i][x+y*w] = valueY
			}
		}
	}
	secretImgData := []uint8{}
	extBit := false
	for i := 0; i < w*h; i++ {
		points := [][2]int{}
		for j := 0; j < len(imgs); j++ {
			points = append(points, [2]int{imgs[j].ID, int(grayImgData[j][i])})
		}
		result := lagrangeCoefficients(points, 251)
		for _, coe := range result {
			if coe == 0 {
				continue
			}
			if extBit {
				secretImgData = append(secretImgData, uint8(250+coe-2))
				extBit = false
			} else {
				if coe == 250 {
					extBit = true
				} else {
					secretImgData = append(secretImgData, uint8(coe-1))
				}
			}
		}
	}
	fullH := (len(secretImgData)-1)/w + 1
	secretImg := image.NewGray(image.Rect(0, 0, w, fullH))
	for i := 0; i < len(secretImgData); i++ {
		secretImg.SetGray(i%w, i/w, color.Gray{Y: secretImgData[i]})
	}
	fileName := "decrypted.bmp"
	if TimeTag {
		fileName = fmt.Sprintf("decrypted_%s.bmp", time.Now().Format("2006-01-02_15-04-05"))
	}
	saveImage(secretImg, fileName)
}

// utils

func pow(a, b, p int) int {
	ans := 1
	base := a
	for i := b; i > 0; i >>= 1 {
		if i&1 == 1 {
			ans = ans * base % p
		}
		base = base * base % p
	}
	return ans
}

func lagrangeCoefficients(points [][2]int, p int) []int {
	len := len(points)
	ans := make([]int, len)

	for i, point := range points {
		coe := make([]int, len)
		coe[0] = point[1]
		down := 1
		for j, otherPoint := range points {
			if i == j {
				continue
			}
			down *= point[0] - otherPoint[0]
			pre := 0
			for k, v := range coe {
				coe[k] = v*(-otherPoint[0]) + pre
				coe[k] %= p
				pre = v
			}
		}
		// 费马小定理求逆元
		down = pow(down, p-2, p)
		for j := range ans {
			ans[j] += coe[j] * down
			ans[j] = (ans[j]%p + p) % p
		}
	}
	return ans
}
