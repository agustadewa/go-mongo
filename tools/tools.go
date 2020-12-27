package tools

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"io"
	"log"
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
	pdf.SetFontLocation("./TEMP/FONT")
	pdf.AddFont("ArchivoBlack-Regular", "", "ArchivoBlack-Regular.json")
	pdf.SetFontLocation("./TEMP/FONT")
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
	pdf.SetFontLocation("./TEMP/FONT")
	pdf.AddFont("OrangeTypewriter", "", "OrangeTypewriter.json")
	pdf.SetFontLocation("./TEMP/FONT")
	pdf.AddFont("Kanit-Bold", "", "Kanit-Bold.json")

	pdf.SetHeaderFunc(func() {
		// pdf.Image("./assets/templates/template1.jpg", 0, 0, 297, 200, true, "", 0, "")
		pdf.ImageOptions(templatePath, 0, 0, 297, 210, false, gofpdf.ImageOptions{ImageType: fileType, ReadDpi: true}, 0, "")

		// CALL SIGN
		pdf.SetFont("Kanit-Bold", "", 48)
		pdf.SetXY(247, 80)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(40, 10, callSign, "", 0, "R", false, 0, "")

		// NAME
		pdf.SetFont("Kanit-Bold", "", 18)
		pdf.SetXY(276, 95)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(10, 10, name, "", 0, "R", false, 0, "")

		// FREQUENCY
		pdf.SetFont("OrangeTypewriter", "", 16)
		pdf.SetTextColor(0, 0, 0)
		if band == "40 m" {
			pdf.SetXY(279, 23)
			pdf.CellFormat(10, 10, "7.135 MHz", "", 0, "R", false, 0, "")
		} else if band == "2 m" {
			pdf.SetXY(279, 23)
			pdf.CellFormat(10, 10, "145.240 MHz", "", 0, "R", false, 0, "")
		}
	})

	err := pdf.Output(w)
	if err != nil {
		log.Println("error creating pdf:", err)
	}
	return err
}

// SaveImageFromB64 method
func (tool Tools) SaveImageFromB64(b64 string, filePath string) error {
	unbased, errDecode := base64.StdEncoding.Strict().DecodeString(b64)
	if errDecode != nil {
		fmt.Println(errDecode)
	}

	reader := bytes.NewReader(unbased)

	// Decode JPG
	//img, errDecodeJpeg := jpeg.Decode(reader);
	//if  errDecodeJpeg != nil {
	//	panic("BAD JPG")
	//}

	// Decode PNG
	img, errDecodePng := png.Decode(reader)
	if errDecodePng != nil {
		panic("BAD PNG")
	}

	// Create File
	file, errCreateFile := os.Create(filePath)
	if errCreateFile != nil {
		panic(errCreateFile)
	}
	defer file.Close()

	// Encode JPG
	//jpgOpt := jpeg.Options{Quality: 30}
	//errEncodeJpg = jpeg.Encode(f, img, &jpgOpt)

	// Encode PNG
	errEncodePng := png.Encode(file, img)
	if errEncodePng != nil {
		fmt.Println(errEncodePng)
	}

	//return errEncodeJpg
	return errEncodePng
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
