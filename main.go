// Copyright 2023 Jens Weißkopf. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/weisskopfjens/goconvertis2/convertis2"
)

func main() {
	fmt.Println("goConvertIS2 by (c)Jens Weißkopf (github.com/weisskopfjens/goconvertis)")
	fmt.Println("(*) are required parameter.")

	iPtr := flag.String("i", "", "(*) A .is2 File.")
	oIRPtr := flag.String("oi", "ir.jpg", "A .jpg file for infrared output.")
	oVISPtr := flag.String("ov", "vis.jpg", "A .jpg file for visual output.")
	bgtempPtr := flag.Float64("b", 20.0, "Background temperature.")
	emissionPtr := flag.Float64("e", 0.95, "Emission factor.")
	mintempPtr := flag.Float64("min", 20.0, "Min. temperature.")
	maxtempPtr := flag.Float64("max", 70.0, "Max. temperature.")
	flag.Parse()

	if *iPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	convertis2.ConvertIS2(*iPtr, *oIRPtr, *oVISPtr, *bgtempPtr, *emissionPtr, *mintempPtr, *maxtempPtr)
}
