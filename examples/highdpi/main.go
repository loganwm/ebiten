// Copyright 2018 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build example

package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"net/http"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	count        int
	highDPIImage *ebiten.Image
)

func init() {
	// licensed under Public Domain
	// https://commons.wikimedia.org/wiki/File:As08-16-2593.jpg
	const url = "https://upload.wikimedia.org/wikipedia/commons/1/1f/As08-16-2593.jpg"

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	highDPIImage, err = ebiten.NewImageFromImage(img, ebiten.FilterLinear)
	if err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {
	if ebiten.IsRunningSlowly() {
		return nil
	}

	scale := ebiten.DeviceScaleFactor()
	sw, sh := screen.Size()

	w, h := highDPIImage.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(-w)/2, float64(-h)/2)
	// The image is just too big. Adjust the scale.
	op.GeoM.Scale(0.25, 0.25)
	// Scale the image by the device ratio so that the rendering result can be same
	// on various (diffrent-DPI) environments.
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(sw)/2, float64(sh)/2)
	screen.DrawImage(highDPIImage, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Device Scale Ratio: %0.2f", scale))
	return nil
}

func main() {
	const (
		screenWidth  = 640
		screenHeight = 480
	)

	// Pass the invert of scale so that Ebiten's auto scaling by device scale is disabled.
	s := ebiten.DeviceScaleFactor()
	if err := ebiten.Run(update, int(screenWidth*s), int(screenHeight*s), 1/s, "High DPI (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
}
