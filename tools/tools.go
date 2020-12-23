package tools

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"math/rand"
	"os"
	"regexp"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// Tools type
type Tools struct{}

// PrintPDF method
func (tool Tools) PrintPDF(name, callSign, band, templatePath, outPath, fileType string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetFontLocation("./public/TEMP/FONT")
	pdf.AddFont("ArchivoBlack-Regular", "", "ArchivoBlack-Regular.json")
	pdf.SetFontLocation("./public/TEMP/FONT")
	pdf.AddFont("ATOMICCLOCKRADIO", "", "ATOMICCLOCKRADIO.json")

	pdf.SetHeaderFunc(func() {
		// pdf.Image("./assets/templates/template1.jpg", 0, 0, 297, 200, true, "", 0, "")
		pdf.ImageOptions(templatePath, 0, 0, 297, 210, false, gofpdf.ImageOptions{ImageType: fileType, ReadDpi: true}, 0, "")

		pdf.SetFont("ArchivoBlack-Regular", "", 47)
		pdf.SetXY(4, 91)
		pdf.SetTextColor(12, 168, 149)
		pdf.Cell(40, 10, callSign)

		pdf.SetFont("ArchivoBlack-Regular", "", 25)
		pdf.SetXY(6, 105)
		pdf.SetTextColor(12, 168, 149)
		pdf.Cell(10, 10, name)

		pdf.SetFont("ATOMICCLOCKRADIO", "", 23)
		pdf.SetTextColor(255, 255, 255)
		if band == "40 m" {
			pdf.SetXY(131, 43)
			pdf.Cell(10, 10, "7.135")
		} else if band == "2 m" {
			pdf.SetXY(119, 43)
			pdf.Cell(10, 10, "145.240")
		}
	})

	err := pdf.OutputFileAndClose(outPath)
	if err != nil {
		fmt.Println("ERRRORRR! ", err)
	}
	return err
}

// PrintPDF method
func (tool Tools) PrintPDFV2(name, callSign, band, templatePath, fileType string, w io.Writer) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetFontLocation("./public/TEMP/FONT")
	pdf.AddFont("ArchivoBlack-Regular", "", "ArchivoBlack-Regular.json")
	pdf.SetFontLocation("./public/TEMP/FONT")
	pdf.AddFont("ATOMICCLOCKRADIO", "", "ATOMICCLOCKRADIO.json")

	pdf.SetHeaderFunc(func() {
		// pdf.Image("./assets/templates/template1.jpg", 0, 0, 297, 200, true, "", 0, "")
		pdf.ImageOptions(templatePath, 0, 0, 297, 210, false, gofpdf.ImageOptions{ImageType: fileType, ReadDpi: true}, 0, "")

		pdf.SetFont("ArchivoBlack-Regular", "", 47)
		pdf.SetXY(4, 91)
		pdf.SetTextColor(12, 168, 149)
		pdf.Cell(40, 10, callSign)

		pdf.SetFont("ArchivoBlack-Regular", "", 25)
		pdf.SetXY(6, 105)
		pdf.SetTextColor(12, 168, 149)
		pdf.Cell(10, 10, name)

		pdf.SetFont("ATOMICCLOCKRADIO", "", 23)
		pdf.SetTextColor(255, 255, 255)
		if band == "40 m" {
			pdf.SetXY(131, 43)
			pdf.Cell(10, 10, "7.135")
		} else if band == "2 m" {
			pdf.SetXY(119, 43)
			pdf.Cell(10, 10, "145.240")
		}
	})

	err := pdf.Output(w)
	if err != nil {
		fmt.Println("ERRRORRR! ", err)
	}
	return err
}

// SaveImageFromB64 method
func (tool Tools) SaveImageFromB64(b64 *string, filePath string) error {
	var err error

	var unbased []byte
	unbased, err = base64.StdEncoding.Strict().DecodeString(*b64)

	reader := bytes.NewReader(unbased)
	var img image.Image
	if img, err = jpeg.Decode(reader); err != nil {
		panic("BAD JPG")
	}

	// if img, err = png.Decode(reader); err != nil {
	// 	panic("BAD PNG")
	// }

	var f *os.File
	if f, err = os.Create(filePath); err != nil {
		panic(err)
	}
	defer f.Close()

	opt := jpeg.Options{Quality: 30}
	err = jpeg.Encode(f, img, &opt)

	// err = png.Encode(f, img)

	return err
}

// ReplaceRegex method
func (tool Tools) ReplaceRegex(b64url *string, b64 *string) {
	var re = regexp.MustCompile(`^[^,]*,`)
	*b64 = re.ReplaceAllString(*b64url, "")
}

// RNDString method
func (tool Tools) RNDString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
