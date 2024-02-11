// Copyright 2023 Jens Weißkopf. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package convertis2

import (
	"archive/zip"
	"encoding/binary"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cryptix/wav"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
)

// Audio 784080
// ConvertIS2 converts FLUKE .IS2 files in a infrared picture and a visual picture (.jpg)
func ConvertIS2(filename string, irfilepath string, visfilepath string, bgtemp float64, emission float64, mintemp float64, maxtemp float64) {
	// Fileformat is2
	// 0 unknown
	// 1 Old is2 format (raw,uncompressed,binary)
	// 2 New is2 format (zip based format)
	// fileversion := 0

	err := decodeNewIS2(filename, irfilepath, visfilepath, bgtemp, emission, mintemp, maxtemp)
	if err != nil {
		err2 := decodeOldIS2(filename, irfilepath, visfilepath, bgtemp, emission, mintemp, maxtemp)
		if err2 == nil {
			//fileversion = 1
			log.Println("Fileversion 1 detected.")
		} else {
			log.Fatalln("Unknown is2 format.", err2)
		}
	} else {
		//fileversion = 2
		log.Println("Fileversion 2 detected.")
	}

}

func decodeIS2IRBinaryData2JPG(offset int64, file *os.File, filename string, bgtemp float64, emission float64, mintemp float64, maxtemp float64) error {
	_, err := file.Seek(offset, 0)
	if err != nil {
		return err
	}
	var w uint16
	var r, g, b uint8
	// fluke hot iron palette
	ironpalette := []string{"#00000a", "#000014", "#00001e", "#000025", "#00002a", "#00002e", "#000032", "#000036", "#00003a", "#00003e", "#000042", "#000046", "#00004a", "#00004f", "#000052", "#010055", "#010057", "#020059", "#02005c", "#03005e", "#040061", "#040063", "#050065", "#060067", "#070069", "#08006b", "#09006e", "#0a0070", "#0b0073", "#0c0074", "#0d0075", "#0d0076", "#0e0077", "#100078", "#120079", "#13007b", "#15007c", "#17007d", "#19007e", "#1b0080", "#1c0081", "#1e0083", "#200084", "#220085", "#240086", "#260087", "#280089", "#2a0089", "#2c008a", "#2e008b", "#30008c", "#32008d", "#34008e", "#36008e", "#38008f", "#390090", "#3b0091", "#3c0092", "#3e0093", "#3f0093", "#410094", "#420095", "#440095", "#450096", "#470096", "#490096", "#4a0096", "#4c0097", "#4e0097", "#4f0097", "#510097", "#520098", "#540098", "#560098", "#580099", "#5a0099", "#5c0099", "#5d009a", "#5f009a", "#61009b", "#63009b", "#64009b", "#66009b", "#68009b", "#6a009b", "#6c009c", "#6d009c", "#6f009c", "#70009c", "#71009d", "#73009d", "#75009d", "#77009d", "#78009d", "#7a009d", "#7c009d", "#7e009d", "#7f009d", "#81009d", "#83009d", "#84009d", "#86009d", "#87009d", "#89009d", "#8a009d", "#8b009d", "#8d009d", "#8f009c", "#91009c", "#93009c", "#95009c", "#96009b", "#98009b", "#99009b", "#9b009b", "#9c009b", "#9d009b", "#9f009b", "#a0009b", "#a2009b", "#a3009b", "#a4009b", "#a6009a", "#a7009a", "#a8009a", "#a90099", "#aa0099", "#ab0099", "#ad0099", "#ae0198", "#af0198", "#b00198", "#b00198", "#b10197", "#b20197", "#b30196", "#b40296", "#b50295", "#b60295", "#b70395", "#b80395", "#b90495", "#ba0495", "#ba0494", "#bb0593", "#bc0593", "#bd0593", "#be0692", "#bf0692", "#bf0692", "#c00791", "#c00791", "#c10890", "#c10990", "#c20a8f", "#c30a8e", "#c30b8e", "#c40c8d", "#c50c8c", "#c60d8b", "#c60e8a", "#c70f89", "#c81088", "#c91187", "#ca1286", "#ca1385", "#cb1385", "#cb1484", "#cc1582", "#cd1681", "#ce1780", "#ce187e", "#cf187c", "#cf197b", "#d01a79", "#d11b78", "#d11c76", "#d21c75", "#d21d74", "#d31e72", "#d32071", "#d4216f", "#d4226e", "#d5236b", "#d52469", "#d62567", "#d72665", "#d82764", "#d82862", "#d92a60", "#da2b5e", "#da2c5c", "#db2e5a", "#db2f57", "#dc2f54", "#dd3051", "#dd314e", "#de324a", "#de3347", "#df3444", "#df3541", "#df363d", "#e0373a", "#e03837", "#e03933", "#e13a30", "#e23b2d", "#e23c2a", "#e33d26", "#e33e23", "#e43f20", "#e4411d", "#e4421c", "#e5431b", "#e54419", "#e54518", "#e64616", "#e74715", "#e74814", "#e74913", "#e84a12", "#e84c10", "#e84c0f", "#e94d0e", "#e94d0d", "#ea4e0c", "#ea4f0c", "#eb500b", "#eb510a", "#eb520a", "#eb5309", "#ec5409", "#ec5608", "#ec5708", "#ec5808", "#ed5907", "#ed5a07", "#ed5b06", "#ee5c06", "#ee5c05", "#ee5d05", "#ee5e05", "#ef5f04", "#ef6004", "#ef6104", "#ef6204", "#f06303", "#f06403", "#f06503", "#f16603", "#f16603", "#f16703", "#f16803", "#f16902", "#f16a02", "#f16b02", "#f16b02", "#f26c01", "#f26d01", "#f26e01", "#f36f01", "#f37001", "#f37101", "#f37201", "#f47300", "#f47400", "#f47500", "#f47600", "#f47700", "#f47800", "#f47a00", "#f57b00", "#f57c00", "#f57e00", "#f57f00", "#f68000", "#f68100", "#f68200", "#f78300", "#f78400", "#f78500", "#f78600", "#f88700", "#f88800", "#f88800", "#f88900", "#f88a00", "#f88b00", "#f88c00", "#f98d00", "#f98d00", "#f98e00", "#f98f00", "#f99000", "#f99100", "#f99200", "#f99300", "#fa9400", "#fa9500", "#fa9600", "#fb9800", "#fb9900", "#fb9a00", "#fb9c00", "#fc9d00", "#fc9f00", "#fca000", "#fca100", "#fda200", "#fda300", "#fda400", "#fda600", "#fda700", "#fda800", "#fdaa00", "#fdab00", "#fdac00", "#fdad00", "#fdae00", "#feaf00", "#feb000", "#feb100", "#feb200", "#feb300", "#feb400", "#feb500", "#feb600", "#feb800", "#feb900", "#feb900", "#feba00", "#febb00", "#febc00", "#febd00", "#febe00", "#fec000", "#fec100", "#fec200", "#fec300", "#fec400", "#fec500", "#fec600", "#fec700", "#fec800", "#fec901", "#feca01", "#feca01", "#fecb01", "#fecc02", "#fecd02", "#fece03", "#fecf04", "#fecf04", "#fed005", "#fed106", "#fed308", "#fed409", "#fed50a", "#fed60a", "#fed70b", "#fed80c", "#fed90d", "#ffda0e", "#ffda0e", "#ffdb10", "#ffdc12", "#ffdc14", "#ffdd16", "#ffde19", "#ffde1b", "#ffdf1e", "#ffe020", "#ffe122", "#ffe224", "#ffe226", "#ffe328", "#ffe42b", "#ffe42e", "#ffe531", "#ffe635", "#ffe638", "#ffe73c", "#ffe83f", "#ffe943", "#ffea46", "#ffeb49", "#ffeb4d", "#ffec50", "#ffed54", "#ffee57", "#ffee5b", "#ffee5f", "#ffef63", "#ffef67", "#fff06a", "#fff06e", "#fff172", "#fff177", "#fff17b", "#fff280", "#fff285", "#fff28a", "#fff38e", "#fff492", "#fff496", "#fff49a", "#fff59e", "#fff5a2", "#fff5a6", "#fff6aa", "#fff6af", "#fff7b3", "#fff7b6", "#fff8ba", "#fff8bd", "#fff8c1", "#fff8c4", "#fff9c7", "#fff9ca", "#fff9cd", "#fffad1", "#fffad4", "#fffbd8", "#fffcdb", "#fffcdf", "#fffde2", "#fffde5", "#fffde8", "#fffeeb", "#fffeee", "#fffef1", "#fffef4", "#fffff6"}
	minvalue := uint16(65535)
	maxvalue := uint16(0)
	var mintemppointx int
	var mintemppointy int
	var maxtemppointx int
	var maxtemppointy int
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			err := binary.Read(file, binary.LittleEndian, &w)
			if err != nil {
				return err
			}
			if minvalue > w {
				minvalue = w
				mintemppointx = x
				mintemppointy = y
			}
			if maxvalue < w {
				maxvalue = w
				maxtemppointx = x
				maxtemppointy = y
			}
		}
	}
	mintemperature := raypower2degrees(uint16(float64(minvalue)*0.662+228), bgtemp, emission)
	maxtemperature := raypower2degrees(uint16(float64(maxvalue)*0.662+228), bgtemp, emission)
	log.Printf("Min. and Max. temperature in the file:\n")
	log.Printf("Temperature min=%.2f °C\n", mintemperature)
	log.Printf("Temperature max=%.2f °C\n", maxtemperature)

	mintemperaturescale := mintemperature
	maxtemperaturescale := maxtemperature
	if !(mintemp == 0.0 && maxtemp == 0.0) {
		mintemperaturescale = mintemp
		maxtemperaturescale = maxtemp
		log.Printf("Manuel scale of colortable:\n")
		log.Printf("Temperature min=%.2f °C\n", mintemperaturescale)
		log.Printf("Temperature max=%.2f °C\n", maxtemperaturescale)
	} else {
		log.Printf("Automatic scale of the colortable.\n")
	}
	colorstep := 433 / (maxtemperaturescale - mintemperaturescale)

	log.Printf("Backgroundtemperature=%.2f °C\n", bgtemp)
	log.Printf("Emission factor=%.2f\n", emission)
	irImage := gg.NewContext(390, 240)
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	fontbold, err := truetype.Parse(gobold.TTF)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	face := truetype.NewFace(font, &truetype.Options{Size: 14})
	irImage.SetFontFace(face)
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}
	irImage.SetRGBA(1, 1, 1, 1)
	irImage.Clear()
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			err = binary.Read(file, binary.LittleEndian, &w)
			if err != nil {
				return err
			}
			temperature := raypower2degrees(uint16(float64(w)*0.662+228), bgtemp, emission)
			ci := AbsFloat64(mintemperaturescale-temperature) * colorstep
			if ci > 432 {
				ci = 432
			}
			r, g, b = HTMLColorToRGB(ironpalette[int64(ci)])
			irImage.SetRGB255(int(r), int(g), int(b))
			irImage.SetPixel(x, y)
		}
	}
	colorstep = 433.0 / 221.0
	for y := 0; y < 221; y++ {
		ci := 432 - colorstep*float64(y)
		if ci >= 433 {
			ci = 432
		}
		r, g, b = HTMLColorToRGB(ironpalette[int(ci)])
		irImage.SetLineWidth(1)
		irImage.SetRGB255(int(r), int(g), int(b))
		irImage.DrawLine(320, float64(y), 335, float64(y))
		irImage.Stroke()
	}
	irImage.SetRGB255(0, 0, 0)
	irImage.DrawRectangle(320, 0, 15, 220)
	irImage.Stroke()
	irImage.DrawLine(320, 8, 335, 8)
	irImage.DrawString(fmt.Sprintf("%.1f", maxtemperaturescale), 353, 13)
	irImage.DrawLine(320, 213, 335, 213)
	irImage.DrawString(fmt.Sprintf("%.1f", mintemperaturescale), 353, 219)
	irImage.DrawLine(344, 8, 344, 213)
	irImage.DrawLine(344, 8, 350, 8)
	tempstep := (maxtemperaturescale - mintemperaturescale) / 9
	for i := 24; i < 224; i = i + 25 {
		temp := (tempstep * ((((224 - float64(i)) - 24) / 25) + 1)) + mintemperaturescale
		irImage.DrawLine(344, float64(i), 350, float64(i))
		irImage.DrawString(fmt.Sprintf("%.0f", temp), 353, float64(i)+4)
	}
	irImage.DrawLine(344, 213, 350, 213)
	irImage.DrawString("°C", 346, 234)
	irImage.Stroke()
	irImage.SetRGB255(0, 0, 0)
	face2 := truetype.NewFace(fontbold, &truetype.Options{Size: 13})
	irImage.SetFontFace(face2)
	irImage.SetLineWidth(4)
	irImage.DrawLine(float64(mintemppointx)-4, float64(mintemppointy), float64(mintemppointx)+4, float64(mintemppointy))
	irImage.DrawLine(float64(mintemppointx), float64(mintemppointy)-4, float64(mintemppointx), float64(mintemppointy)+4)
	irImage.DrawString(fmt.Sprintf("%.1f", mintemperature), float64(mintemppointx)-12, float64(mintemppointy)-6)
	irImage.Stroke()
	irImage.DrawLine(float64(maxtemppointx)-4, float64(maxtemppointy), float64(maxtemppointx)+4, float64(maxtemppointy))
	irImage.DrawLine(float64(maxtemppointx), float64(maxtemppointy)-4, float64(maxtemppointx), float64(maxtemppointy)+4)
	irImage.DrawString(fmt.Sprintf("%.1f", maxtemperature), float64(maxtemppointx)-12, float64(maxtemppointy)-6)
	irImage.Stroke()
	irImage.SetRGBA255(200, 200, 255, 230)
	face = truetype.NewFace(font, &truetype.Options{Size: 12})
	irImage.SetFontFace(face)
	irImage.SetLineWidth(1)
	irImage.DrawLine(float64(mintemppointx)-3, float64(mintemppointy), float64(mintemppointx)+3, float64(mintemppointy))
	irImage.DrawLine(float64(mintemppointx), float64(mintemppointy)-3, float64(mintemppointx), float64(mintemppointy)+3)
	irImage.DrawString(fmt.Sprintf("%.1f", mintemperature), float64(mintemppointx)-12, float64(mintemppointy)-6)
	irImage.Stroke()
	irImage.SetRGBA255(255, 200, 200, 230)
	irImage.DrawLine(float64(maxtemppointx)-2, float64(maxtemppointy), float64(maxtemppointx)+2, float64(maxtemppointy))
	irImage.DrawLine(float64(maxtemppointx), float64(maxtemppointy)-2, float64(maxtemppointx), float64(maxtemppointy)+2)
	irImage.DrawString(fmt.Sprintf("%.1f", maxtemperature), float64(maxtemppointx)-12, float64(maxtemppointy)-6)
	irImage.Stroke()
	outFile, _ := os.Create(filename)
	defer outFile.Close()
	jpegerr := jpeg.Encode(outFile, irImage.Image(), &jpeg.Options{Quality: 100})
	if jpegerr != nil {
		log.Fatalln("Can't encode jpeg format.", jpegerr, filename)
	}
	return nil
}

func decodeNEWIS2IRBinaryData2JPG(offset int64, file *os.File, filename string, bgtemp float64, emission float64, mintemp float64, maxtemp float64) error {
	_, err := file.Seek(offset, 0)
	if err != nil {
		return err
	}
	var w uint16
	var r, g, b uint8
	// fluke hot iron palette
	ironpalette := []string{"#00000a", "#000014", "#00001e", "#000025", "#00002a", "#00002e", "#000032", "#000036", "#00003a", "#00003e", "#000042", "#000046", "#00004a", "#00004f", "#000052", "#010055", "#010057", "#020059", "#02005c", "#03005e", "#040061", "#040063", "#050065", "#060067", "#070069", "#08006b", "#09006e", "#0a0070", "#0b0073", "#0c0074", "#0d0075", "#0d0076", "#0e0077", "#100078", "#120079", "#13007b", "#15007c", "#17007d", "#19007e", "#1b0080", "#1c0081", "#1e0083", "#200084", "#220085", "#240086", "#260087", "#280089", "#2a0089", "#2c008a", "#2e008b", "#30008c", "#32008d", "#34008e", "#36008e", "#38008f", "#390090", "#3b0091", "#3c0092", "#3e0093", "#3f0093", "#410094", "#420095", "#440095", "#450096", "#470096", "#490096", "#4a0096", "#4c0097", "#4e0097", "#4f0097", "#510097", "#520098", "#540098", "#560098", "#580099", "#5a0099", "#5c0099", "#5d009a", "#5f009a", "#61009b", "#63009b", "#64009b", "#66009b", "#68009b", "#6a009b", "#6c009c", "#6d009c", "#6f009c", "#70009c", "#71009d", "#73009d", "#75009d", "#77009d", "#78009d", "#7a009d", "#7c009d", "#7e009d", "#7f009d", "#81009d", "#83009d", "#84009d", "#86009d", "#87009d", "#89009d", "#8a009d", "#8b009d", "#8d009d", "#8f009c", "#91009c", "#93009c", "#95009c", "#96009b", "#98009b", "#99009b", "#9b009b", "#9c009b", "#9d009b", "#9f009b", "#a0009b", "#a2009b", "#a3009b", "#a4009b", "#a6009a", "#a7009a", "#a8009a", "#a90099", "#aa0099", "#ab0099", "#ad0099", "#ae0198", "#af0198", "#b00198", "#b00198", "#b10197", "#b20197", "#b30196", "#b40296", "#b50295", "#b60295", "#b70395", "#b80395", "#b90495", "#ba0495", "#ba0494", "#bb0593", "#bc0593", "#bd0593", "#be0692", "#bf0692", "#bf0692", "#c00791", "#c00791", "#c10890", "#c10990", "#c20a8f", "#c30a8e", "#c30b8e", "#c40c8d", "#c50c8c", "#c60d8b", "#c60e8a", "#c70f89", "#c81088", "#c91187", "#ca1286", "#ca1385", "#cb1385", "#cb1484", "#cc1582", "#cd1681", "#ce1780", "#ce187e", "#cf187c", "#cf197b", "#d01a79", "#d11b78", "#d11c76", "#d21c75", "#d21d74", "#d31e72", "#d32071", "#d4216f", "#d4226e", "#d5236b", "#d52469", "#d62567", "#d72665", "#d82764", "#d82862", "#d92a60", "#da2b5e", "#da2c5c", "#db2e5a", "#db2f57", "#dc2f54", "#dd3051", "#dd314e", "#de324a", "#de3347", "#df3444", "#df3541", "#df363d", "#e0373a", "#e03837", "#e03933", "#e13a30", "#e23b2d", "#e23c2a", "#e33d26", "#e33e23", "#e43f20", "#e4411d", "#e4421c", "#e5431b", "#e54419", "#e54518", "#e64616", "#e74715", "#e74814", "#e74913", "#e84a12", "#e84c10", "#e84c0f", "#e94d0e", "#e94d0d", "#ea4e0c", "#ea4f0c", "#eb500b", "#eb510a", "#eb520a", "#eb5309", "#ec5409", "#ec5608", "#ec5708", "#ec5808", "#ed5907", "#ed5a07", "#ed5b06", "#ee5c06", "#ee5c05", "#ee5d05", "#ee5e05", "#ef5f04", "#ef6004", "#ef6104", "#ef6204", "#f06303", "#f06403", "#f06503", "#f16603", "#f16603", "#f16703", "#f16803", "#f16902", "#f16a02", "#f16b02", "#f16b02", "#f26c01", "#f26d01", "#f26e01", "#f36f01", "#f37001", "#f37101", "#f37201", "#f47300", "#f47400", "#f47500", "#f47600", "#f47700", "#f47800", "#f47a00", "#f57b00", "#f57c00", "#f57e00", "#f57f00", "#f68000", "#f68100", "#f68200", "#f78300", "#f78400", "#f78500", "#f78600", "#f88700", "#f88800", "#f88800", "#f88900", "#f88a00", "#f88b00", "#f88c00", "#f98d00", "#f98d00", "#f98e00", "#f98f00", "#f99000", "#f99100", "#f99200", "#f99300", "#fa9400", "#fa9500", "#fa9600", "#fb9800", "#fb9900", "#fb9a00", "#fb9c00", "#fc9d00", "#fc9f00", "#fca000", "#fca100", "#fda200", "#fda300", "#fda400", "#fda600", "#fda700", "#fda800", "#fdaa00", "#fdab00", "#fdac00", "#fdad00", "#fdae00", "#feaf00", "#feb000", "#feb100", "#feb200", "#feb300", "#feb400", "#feb500", "#feb600", "#feb800", "#feb900", "#feb900", "#feba00", "#febb00", "#febc00", "#febd00", "#febe00", "#fec000", "#fec100", "#fec200", "#fec300", "#fec400", "#fec500", "#fec600", "#fec700", "#fec800", "#fec901", "#feca01", "#feca01", "#fecb01", "#fecc02", "#fecd02", "#fece03", "#fecf04", "#fecf04", "#fed005", "#fed106", "#fed308", "#fed409", "#fed50a", "#fed60a", "#fed70b", "#fed80c", "#fed90d", "#ffda0e", "#ffda0e", "#ffdb10", "#ffdc12", "#ffdc14", "#ffdd16", "#ffde19", "#ffde1b", "#ffdf1e", "#ffe020", "#ffe122", "#ffe224", "#ffe226", "#ffe328", "#ffe42b", "#ffe42e", "#ffe531", "#ffe635", "#ffe638", "#ffe73c", "#ffe83f", "#ffe943", "#ffea46", "#ffeb49", "#ffeb4d", "#ffec50", "#ffed54", "#ffee57", "#ffee5b", "#ffee5f", "#ffef63", "#ffef67", "#fff06a", "#fff06e", "#fff172", "#fff177", "#fff17b", "#fff280", "#fff285", "#fff28a", "#fff38e", "#fff492", "#fff496", "#fff49a", "#fff59e", "#fff5a2", "#fff5a6", "#fff6aa", "#fff6af", "#fff7b3", "#fff7b6", "#fff8ba", "#fff8bd", "#fff8c1", "#fff8c4", "#fff9c7", "#fff9ca", "#fff9cd", "#fffad1", "#fffad4", "#fffbd8", "#fffcdb", "#fffcdf", "#fffde2", "#fffde5", "#fffde8", "#fffeeb", "#fffeee", "#fffef1", "#fffef4", "#fffff6"}
	minvalue := uint16(65535)
	maxvalue := uint16(0)
	var mintemppointx int
	var mintemppointy int
	var maxtemppointx int
	var maxtemppointy int
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			err := binary.Read(file, binary.LittleEndian, &w)
			if err != nil {
				return err
			}
			if minvalue > w {
				minvalue = w
				mintemppointx = x
				mintemppointy = y
			}
			if maxvalue < w {
				maxvalue = w
				maxtemppointx = x
				maxtemppointy = y
			}
		}
	}
	//
	// DEBUG
	// For calibrationa a picture with min temp of 9.9 and a max temp of 18.9
	// For example:
	// calibrate(9.9, 18.9, float64(minvalue), float64(maxvalue), bgtemp, emission)
	// Result = 0.20100000000000015 154.0349999999324
	// The resulting function with calibrated values:
	// raypower2degrees(uint16(float64(minvalue)*0.201+154.035), bgtemp, emission)
	//
	mintemperature := raypower2degrees(uint16(float64(minvalue)*0.201+154.035), bgtemp, emission)
	maxtemperature := raypower2degrees(uint16(float64(maxvalue)*0.201+154.035), bgtemp, emission)
	log.Printf("Min. and Max. temperature in the file:\n")
	log.Printf("Temperature min=%.2f °C\n", mintemperature)
	log.Printf("Temperature max=%.2f °C\n", maxtemperature)

	mintemperaturescale := mintemperature
	maxtemperaturescale := maxtemperature
	if !(mintemp == 0.0 && maxtemp == 0.0) {
		mintemperaturescale = mintemp
		maxtemperaturescale = maxtemp
		log.Printf("Manual scale of the colortable:\n")
		log.Printf("Temperature min=%.2f °C\n", mintemperaturescale)
		log.Printf("Temperature max=%.2f °C\n", maxtemperaturescale)
	} else {
		log.Printf("Automatic scale of the colortable.\n")
	}
	colorstep := 433 / (maxtemperaturescale - mintemperaturescale)

	log.Printf("Backgroundtemperature=%.2f °C\n", bgtemp)
	log.Printf("Emission factor=%.2f\n", emission)

	/*mintemperature := raypower2degrees(uint16(float64(minvalue)*0.201+154.035), bgtemp, emission)
	maxtemperature := raypower2degrees(uint16(float64(maxvalue)*0.201+154.035), bgtemp, emission)
	mintemperaturescale := mintemperature
	maxtemperaturescale := maxtemperature
	if !(mintemp == 0.0 && maxtemp == 0.0) {
		mintemperaturescale = mintemp
		maxtemperaturescale = maxtemp
	}
	colorstep := 433 / (maxtemperaturescale - mintemperaturescale)
	log.Printf("Temperature min=%.2f °C\n", mintemperature)
	log.Printf("Temperature max=%.2f °C\n", maxtemperature)
	log.Printf("Hintergrundtemperatur=%.2f °C\n", bgtemp)
	log.Printf("Emissionsgrad=%.2f\n", emission)*/

	irImage := gg.NewContext(390, 240)
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return err
	}
	fontbold, err := truetype.Parse(gobold.TTF)
	if err != nil {
		return err
	}
	face := truetype.NewFace(font, &truetype.Options{Size: 14})
	irImage.SetFontFace(face)
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}
	irImage.SetRGBA(1, 1, 1, 1)
	irImage.Clear()
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			err = binary.Read(file, binary.LittleEndian, &w)
			if err != nil {
				return err
			}
			temperature := raypower2degrees(uint16(float64(w)*0.201+154.035), bgtemp, emission)
			ci := AbsFloat64(mintemperaturescale-temperature) * colorstep
			if ci > 432 {
				ci = 432
			}
			r, g, b = HTMLColorToRGB(ironpalette[int64(ci)])
			irImage.SetRGB255(int(r), int(g), int(b))
			irImage.SetPixel(x, y)
		}
	}
	colorstep = 433.0 / 221.0
	for y := 0; y < 221; y++ {
		ci := 432 - colorstep*float64(y)
		if ci >= 433 {
			ci = 432
		}
		r, g, b = HTMLColorToRGB(ironpalette[int(ci)])
		irImage.SetLineWidth(1)
		irImage.SetRGB255(int(r), int(g), int(b))
		irImage.DrawLine(320, float64(y), 335, float64(y))
		irImage.Stroke()
	}
	irImage.SetRGB255(0, 0, 0)
	irImage.DrawRectangle(320, 0, 15, 220)
	irImage.Stroke()
	irImage.DrawLine(320, 8, 335, 8)
	irImage.DrawString(fmt.Sprintf("%.1f", maxtemperaturescale), 353, 13)
	irImage.DrawLine(320, 213, 335, 213)
	irImage.DrawString(fmt.Sprintf("%.1f", mintemperaturescale), 353, 219)
	irImage.DrawLine(344, 8, 344, 213)
	irImage.DrawLine(344, 8, 350, 8)
	tempstep := (maxtemperaturescale - mintemperaturescale) / 9
	for i := 24; i < 224; i = i + 25 {
		temp := (tempstep * ((((224 - float64(i)) - 24) / 25) + 1)) + mintemperaturescale
		irImage.DrawLine(344, float64(i), 350, float64(i))
		irImage.DrawString(fmt.Sprintf("%.0f", temp), 353, float64(i)+4)
	}
	irImage.DrawLine(344, 213, 350, 213)
	irImage.DrawString("°C", 346, 234)
	irImage.Stroke()
	irImage.SetRGB255(0, 0, 0)
	face2 := truetype.NewFace(fontbold, &truetype.Options{Size: 13})
	irImage.SetFontFace(face2)
	irImage.SetLineWidth(4)
	irImage.DrawLine(float64(mintemppointx)-4, float64(mintemppointy), float64(mintemppointx)+4, float64(mintemppointy))
	irImage.DrawLine(float64(mintemppointx), float64(mintemppointy)-4, float64(mintemppointx), float64(mintemppointy)+4)
	irImage.DrawString(fmt.Sprintf("%.1f", mintemperature), float64(mintemppointx)-12, float64(mintemppointy)-6)
	irImage.Stroke()
	irImage.DrawLine(float64(maxtemppointx)-4, float64(maxtemppointy), float64(maxtemppointx)+4, float64(maxtemppointy))
	irImage.DrawLine(float64(maxtemppointx), float64(maxtemppointy)-4, float64(maxtemppointx), float64(maxtemppointy)+4)
	irImage.DrawString(fmt.Sprintf("%.1f", maxtemperature), float64(maxtemppointx)-12, float64(maxtemppointy)-6)
	irImage.Stroke()
	irImage.SetRGBA255(200, 200, 255, 230)
	face = truetype.NewFace(font, &truetype.Options{Size: 12})
	irImage.SetFontFace(face)
	irImage.SetLineWidth(1)
	irImage.DrawLine(float64(mintemppointx)-3, float64(mintemppointy), float64(mintemppointx)+3, float64(mintemppointy))
	irImage.DrawLine(float64(mintemppointx), float64(mintemppointy)-3, float64(mintemppointx), float64(mintemppointy)+3)
	irImage.DrawString(fmt.Sprintf("%.1f", mintemperature), float64(mintemppointx)-12, float64(mintemppointy)-6)
	irImage.Stroke()
	irImage.SetRGBA255(255, 200, 200, 230)
	irImage.DrawLine(float64(maxtemppointx)-2, float64(maxtemppointy), float64(maxtemppointx)+2, float64(maxtemppointy))
	irImage.DrawLine(float64(maxtemppointx), float64(maxtemppointy)-2, float64(maxtemppointx), float64(maxtemppointy)+2)
	irImage.DrawString(fmt.Sprintf("%.1f", maxtemperature), float64(maxtemppointx)-12, float64(maxtemppointy)-6)
	irImage.Stroke()
	outFile, _ := os.Create(filename)
	defer outFile.Close()
	jpegerr := jpeg.Encode(outFile, irImage.Image(), &jpeg.Options{Quality: 100})
	if jpegerr != nil {
		log.Fatalln("Can't encode jpeg data.", jpegerr, filename)
	}
	return nil
}

// Decode old fileformat
func decodeOldIS2(filename string, irfilepath string, visfilepath string, bgtemp float64, emission float64, mintemp float64, maxtemp float64) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error while opening file.", err)
		return err
	}
	defer file.Close()
	var bval uint8
	c := 0
	i := 0
	for {
		i = i + 1
		err = binary.Read(file, binary.LittleEndian, &bval)
		if err != nil {
			log.Println(err)
			return err
		}
		if bval == 0xFF {
			c = c + 1
		} else {
			c = 0
		}
		if c >= 20 {
			break
		}
		if i == 1000 {
			return fmt.Errorf("%s Offset not found. Unknown file structure!", filename)
		}
	}
	offset := int64(i)
	if irfilepath != "" {
		fi, err := os.Stat(irfilepath)
		if err != nil {
			log.Println("Create:", irfilepath)
		} else {
			switch mode := fi.Mode(); {
			case mode.IsDir():
				log.Println(irfilepath, "is a directory.")
				friendlyname := strings.Replace(filepath.Base(filename), ".IS2", ".jpg", -1)
				friendlyname = strings.Replace(friendlyname, ".is2", ".jpg", -1)
				irfilepath = irfilepath + "/" + friendlyname
				log.Println(irfilepath)
			case mode.IsRegular():
				log.Println("Overwrite:", irfilepath)
			}
		}
		errdecir := decodeIS2IRBinaryData2JPG(offset+15828, file, irfilepath, bgtemp, emission, mintemp, maxtemp)
		if errdecir != nil {
			log.Fatalln("Can't encode infrared data.", errdecir, irfilepath)
		}
	}
	if visfilepath != "" {
		var w uint16
		var r, g, b uint8
		fi, err := os.Stat(visfilepath)
		if err != nil {
			log.Println("Create:", visfilepath)
		} else {
			switch mode := fi.Mode(); {
			case mode.IsDir():
				log.Println(visfilepath, "is a directory.")
				visfilepath = visfilepath + "/" + strings.Replace(filepath.Base(filename), ".IS2", ".jpg", -1)
			case mode.IsRegular():
				log.Println("Overwrite:", visfilepath)
			}
		}
		_, seekerr := file.Seek(offset+169484, 0)
		if seekerr != nil {
			log.Fatalln("File corrupt. Position [offset+169484 vissual picture data] not found.", seekerr, filepath.Base(filename))
		}
		visImage := gg.NewContext(640, 480)
		for y := 0; y < 480; y++ {
			for x := 0; x < 640; x++ {
				readerr := binary.Read(file, binary.LittleEndian, &w)
				if readerr != nil {
					log.Fatalln("Error while reading file.", readerr, filepath.Base(filename))
				}
				r = uint8(((w >> 11) & 0x1F) * 8)
				g = uint8(((w >> 5) & 0x3F) * 4)
				b = uint8((w & 0x1F) * 8)
				visImage.SetRGB255(int(r), int(g), int(b))
				visImage.SetPixel(x, y)
			}
		}
		outFile2, _ := os.Create(visfilepath)
		defer outFile2.Close()
		jpegerr := jpeg.Encode(outFile2, visImage.Image(), &jpeg.Options{Quality: 100})
		if jpegerr != nil {
			log.Fatalln("Can't encode jpeg.", jpegerr, visfilepath)
		}
	}
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	if fi.Size() > 784080 {
		var w uint16
		_, seekerr := file.Seek(784080, 0)
		if seekerr != nil {
			log.Fatalln("File corrupt. Position [offset+784080 audiodata] not found.", seekerr, filepath.Base(filename))
		}
		wavOut, _ := os.Create(filename + ".wav")
		defer wavOut.Close()
		meta := wav.File{
			Channels:        1,
			SampleRate:      8000,
			SignificantBits: 16,
		}
		writer, _ := meta.NewWriter(wavOut)
		defer writer.Close()
		for {
			err = binary.Read(file, binary.LittleEndian, &w)
			if err != nil {
				if err != io.EOF {
					fmt.Println("read error:", err)
				}
				break
			}
			b := make([]byte, 2)
			binary.LittleEndian.PutUint16(b, uint16(w))
			_, writeerr := writer.Write(b)
			if writeerr != nil {
				log.Fatalln("Error while writing.", writeerr, filename+".wav")
			}
		}
	}
	return nil
}

// Decode new fileformat
func decodeNewIS2(filename string, irfilepath string, visfilepath string, bgtemp float64, emission float64, mintemp float64, maxtemp float64) error {
	tempdir := strings.Replace(filepath.Base(filename), ".IS2", "", -1)
	tempdir = strings.Replace(tempdir, ".is2", "", -1)
	tempdir = "./temp_" + tempdir
	_, err := Unzip(filename, tempdir)
	if err != nil {
		return err
	}
	// 028001E0.jpg
	// 028001E1.jpg
	// IR.data
	file, err := os.Open(tempdir + "/Images/Main/" + "IR.data")
	if err != nil {
		log.Println("Error while opening IR.data", err)
		return err
	}
	defer file.Close()
	fi, err := os.Stat(irfilepath)
	if err != nil {
		log.Println("Erstelle:", irfilepath)
	} else {
		switch mode := fi.Mode(); {
		case mode.IsDir():
			log.Println(irfilepath, "ist ein Verzeichnis.")
			friendlyname := strings.Replace(filepath.Base(filename), ".IS2", ".jpg", -1)
			friendlyname = strings.Replace(friendlyname, ".is2", ".jpg", -1)
			irfilepath = irfilepath + "/" + friendlyname
			log.Println(irfilepath)
		case mode.IsRegular():
			log.Println("Overwrite:", irfilepath)
		}
	}
	decbinerr := decodeNEWIS2IRBinaryData2JPG(640, file, irfilepath, bgtemp, emission, mintemp, maxtemp)
	if decbinerr != nil {
		file.Close()
		err = os.RemoveAll(tempdir)
		if err != nil {
			log.Println("Error while deleting the temporary directory.", err)
		}
		log.Fatalln("Error while decoding ir data.", decbinerr, irfilepath)
	}
	file.Close()

	if visfilepath != "" {
		fi, err := os.Stat(visfilepath)
		if err != nil {
			log.Println("Create:", visfilepath)
		} else {
			switch mode := fi.Mode(); {
			case mode.IsDir():
				log.Println(visfilepath, "is a directory.")
				visfilepath = visfilepath + "/" + strings.Replace(filepath.Base(filename), ".IS2", ".jpg", -1)
			case mode.IsRegular():
				log.Println("Overwrite:", visfilepath)
			}
		}
		_, err = Copyfile(tempdir+"/Images/Main/"+"028001E0.jpg", visfilepath)
		if err != nil {
			log.Println("File:", tempdir+"/Images/Main/"+"028001E0.jpg", ".Can't copy file.")
		}
	}

	err = os.RemoveAll(tempdir)
	if err != nil {
		log.Println("Error while deleting the temporary directory.", err)
	}
	return nil
}

// Abs returns the absolute value of x.
func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// AbsFloat64 Return the absolut value of x
func AbsFloat64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// HTMLColorToRGB convert #RRGGBB to an (R, G, B)
func HTMLColorToRGB(colorstring string) (uint8, uint8, uint8) {
	colorstring = strings.Replace(colorstring, "#", "", -1)
	r, _ := strconv.ParseUint("0x"+colorstring[:2], 0, 8)
	g, _ := strconv.ParseUint("0x"+colorstring[2:4], 0, 8)
	b, _ := strconv.ParseUint("0x"+colorstring[4:], 0, 8)
	return uint8(r), uint8(g), uint8(b)
}

// raypower2degrees converts an raypower value to degree celcius
// using the Stefan Boltzmann Law
func raypower2degrees(P uint16, backgroundtemp float64, emissionfactor float64) float64 {
	// backgroundtemp = 20
	// emissionfactor = 0.95
	// Stefan-Boltzmann Law
	sigma := 5.67e-8 // stefan-boltzmann constant
	//t := (P/(sigma*2.4))**(1.0/4) // Stefan-Boltzmann Law
	t := math.Pow(float64(P)/(sigma*2.4), (1.0 / 4))
	t = t - 273.15 // kelvin in degree
	t = ((t - backgroundtemp) / emissionfactor) + backgroundtemp
	if t < -30 {
		t = -30
	}
	return t
}

// with calibrate you can find the calibration values for a measuring instrument
/*func calibrate(mintemp, maxtemp, minv float64, maxv float64, bgtemp float64, emission float64) {
	for x := 0.01; x < 1; x = x + 0.001 {
		for y := 1.00; y < 200; y = y + 0.005 {
			min := raypower2degrees(uint16(float64(minv)*x+y), bgtemp, emission)
			max := raypower2degrees(uint16(float64(maxv)*x+y), bgtemp, emission)
			if min >= mintemp-0.05 && min <= mintemp+0.05 {
				if max <= maxtemp+0.05 && max >= maxtemp-0.05 {
					fmt.Println("DEBUG Calibration Data:", x, y, min, max)
					return
				}
			}
			//fmt.Println(min, max)
		}
		fmt.Println(x)
	}
}*/

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {
	var filenames []string
	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)
		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: bad directory.", fpath)
		}
		filenames = append(filenames, fpath)
		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}
		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		_, err = io.Copy(outFile, rc)
		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()
		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

// copyfile, copy a file from src to dst
func Copyfile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is a bad file.", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
