package util

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math"
	"os"
)

const DefaultMaxWidth float64 = 100
const DefaultMaxHeight float64 = 100

// 计算图片缩放后的尺寸
func calculateRatioFit(srcWidth, srcHeight int) (int, int) {
	ratio := math.Min(DefaultMaxWidth/float64(srcWidth), DefaultMaxHeight/float64(srcHeight))
	return int(math.Ceil(float64(srcWidth) * ratio)), int(math.Ceil(float64(srcHeight) * ratio))
}

// MakeThumbnail 生成缩略图
func MakeThumbnail(imagePath, savePath string) (string,error)  {
	fmt.Println(imagePath)
	file, _ := os.Open(imagePath)
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "",err
	}

	b := img.Bounds()
	width := b.Max.X
	height := b.Max.Y
	w, h := calculateRatioFit(width, height)

	// 调用resize库进行图片缩放
	m := resize.Resize(uint(w), uint(h), img, resize.Lanczos3)

	// 需要保存的文件
	imgFile, _ := os.Create(savePath)
	defer imgFile.Close()

	// 以PNG格式保存文件
	err = png.Encode(imgFile, m)
	if err != nil {
		return "",err
	}
	return "",nil
}