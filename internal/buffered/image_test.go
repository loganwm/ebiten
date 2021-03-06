// Copyright 2019 The Ebiten Authors
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

package buffered_test

import (
	"image/color"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten"
)

var mainCh = make(chan func())

func runOnMainThread(f func()) {
	ch := make(chan struct{})
	mainCh <- func() {
		f()
		close(ch)
	}
	<-ch
}

func TestMain(m *testing.M) {
	go func() {
		os.Exit(m.Run())
	}()

	for {
		if err := ebiten.Run(func(*ebiten.Image) error {
			f := <-mainCh
			f()
			return nil
		}, 320, 240, 1, ""); err != nil {
			panic(err)
		}
	}
}

type testResult struct {
	want color.RGBA
	got  <-chan color.RGBA
}

var testSetBeforeMainResult = func() testResult {
	clr := color.RGBA{1, 2, 3, 4}
	img, _ := ebiten.NewImage(16, 16, ebiten.FilterDefault)
	img.Set(0, 0, clr)

	ch := make(chan color.RGBA, 1)
	go func() {
		runOnMainThread(func() {
			ch <- img.At(0, 0).(color.RGBA)
		})
	}()

	return testResult{
		want: clr,
		got:  ch,
	}
}()

func TestSetBeforeMain(t *testing.T) {
	got := <-testSetBeforeMainResult.got
	want := testSetBeforeMainResult.want

	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

var testDrawImageBeforeMainResult = func() testResult {
	const w, h = 16, 16
	src, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	dst, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	src.Set(0, 0, color.White)
	dst.DrawImage(src, nil)

	ch := make(chan color.RGBA, 1)
	go func() {
		runOnMainThread(func() {
			ch <- dst.At(0, 0).(color.RGBA)
		})
	}()

	return testResult{
		want: color.RGBA{0xff, 0xff, 0xff, 0xff},
		got:  ch,
	}
}()

func TestDrawImageBeforeMain(t *testing.T) {
	got := <-testDrawImageBeforeMainResult.got
	want := testDrawImageBeforeMainResult.want

	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

var testDrawTrianglesBeforeMainResult = func() testResult {
	const w, h = 16, 16
	src, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	dst, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	src.Set(0, 0, color.White)
	vs := []ebiten.Vertex{
		{
			DstX:   0,
			DstY:   0,
			SrcX:   0,
			SrcY:   0,
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
		{
			DstX:   w,
			DstY:   0,
			SrcX:   w,
			SrcY:   0,
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
		{
			DstX:   0,
			DstY:   h,
			SrcX:   0,
			SrcY:   h,
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
	}
	dst.DrawTriangles(vs, []uint16{0, 1, 2}, src, nil)

	ch := make(chan color.RGBA, 1)
	go func() {
		runOnMainThread(func() {
			ch <- dst.At(0, 0).(color.RGBA)
		})
	}()

	return testResult{
		want: color.RGBA{0xff, 0xff, 0xff, 0xff},
		got:  ch,
	}
}()

func TestDrawTrianglesBeforeMain(t *testing.T) {
	got := <-testDrawTrianglesBeforeMainResult.got
	want := testDrawTrianglesBeforeMainResult.want

	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}
