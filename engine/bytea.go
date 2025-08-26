package engine

import (
	"log"
	"os"
)

// [RU] EncodeImageToBytea читает изображение из файла и возвращает его как []byte <--->
// [ENG] EncodeImageToBytea read images from file and returns it in []byte
func EncodeImageToBytea(filePath string) ([]byte, error) {
	imgData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return imgData, nil
}

// [RU] DecodeByteaToImage сохраняет данные []byte в файл изображения <--->
// [ENG] DecodeByteaToImage saves data in type []byte in image's file
func DecodeByteaToImage(data []byte, outputPath string) error {
	return os.WriteFile(outputPath, data, 0644)
}

var DefaultImage []byte
var defaultImagePath = "C:/Users/Joker/Desktop/GO_Music/GO_Music/static/WhereMyFoto.jpg"

func InitDefaultImage() {
	var err error
	DefaultImage, err = EncodeImageToBytea(defaultImagePath)
	if err != nil {
		log.Fatalf("failed to load default image: %v", err)
	}
}
