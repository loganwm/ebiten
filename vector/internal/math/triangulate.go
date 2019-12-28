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

package math

import (
	"fmt"
)

func cross(v0, v1 Vector) float32 {
	return v0.X*v1.Y - v0.Y*v1.X
}

func triangleCross(pt0, pt1, pt2 Point) float32 {
	return cross(Vector{pt1.X - pt0.X, pt1.Y - pt0.Y}, Vector{pt2.X - pt1.X, pt2.Y - pt1.Y})
}

func adjacentIndices(indices []uint16, idx int) (uint16, uint16, uint16) {
	return indices[(idx+len(indices)-1)%len(indices)], indices[idx], indices[(idx+1)%len(indices)]
}

func InTriangle(pt, pt0, pt1, pt2 Point) bool {
	c0 := cross(Vector{pt.X - pt0.X, pt.Y - pt0.Y}, Vector{pt1.X - pt0.X, pt1.Y - pt0.Y})
	c1 := cross(Vector{pt.X - pt1.X, pt.Y - pt1.Y}, Vector{pt2.X - pt1.X, pt2.Y - pt1.Y})
	c2 := cross(Vector{pt.X - pt2.X, pt.Y - pt2.Y}, Vector{pt0.X - pt2.X, pt0.Y - pt2.Y})
	return (c0 <= 0 && c1 <= 0 && c2 <= 0) || (c0 >= 0 && c1 >= 0 && c2 >= 0)
}

func Triangulate(pts []Point) []uint16 {
	if len(pts) < 3 {
		return nil
	}

	var currentIndices []uint16

	// Remove duplicated points
dup:
	for i := range pts {
		for j := 0; j < i; j++ {
			if pts[i] == pts[j] {
				continue dup
			}
		}
		currentIndices = append(currentIndices, uint16(i))
	}
	if len(currentIndices) < 3 {
		return nil
	}

	// Determine the direction of the polygon from the upper-left point.
	var upperLeft int
	for ci, i := range currentIndices {
		if pts[upperLeft].X < pts[i].X {
			upperLeft = int(ci)
		}
		if pts[upperLeft].X == pts[i].X && pts[upperLeft].Y < pts[i].Y {
			upperLeft = int(ci)
		}
	}
	i0, i1, i2 := adjacentIndices(currentIndices, upperLeft)
	pt0 := pts[i0]
	pt1 := pts[i1]
	pt2 := pts[i2]
	clockwise := triangleCross(pt0, pt1, pt2) < 0

	var indices []uint16

	// Triangulation by Ear Clipping.
	// https://www.geometrictools.com/Documentation/TriangulationByEarClipping.pdf
	for len(currentIndices) >= 3 {
		// Calculate cross-products and remove unneeded vertices.
		cs := make([]float32, len(currentIndices))
		idx := -1
		for i := range currentIndices {
			i0, i1, i2 := adjacentIndices(currentIndices, i)
			pt0 := pts[i0]
			pt1 := pts[i1]
			pt2 := pts[i2]
			c := triangleCross(pt0, pt1, pt2)
			if c == 0 {
				idx = i
				break
			}
			cs[i] = c
		}
		if idx != -1 {
			currentIndices = append(currentIndices[:idx], currentIndices[idx+1:]...)
			continue
		}

		idx = -1
	index:
		for i := range currentIndices {
			i0, i1, i2 := adjacentIndices(currentIndices, i)
			pt0 := pts[i0]
			pt1 := pts[i1]
			pt2 := pts[i2]

			c := cs[i]
			if c == 0 {
				panic("math: cross value must not be 0")
			}
			if c < 0 && !clockwise || c > 0 && clockwise {
				// The angle is more than 180 degrees. This is not an ear.
				continue
			}
			for j := range currentIndices {
				if l := len(currentIndices); j == (i+l-1)%l || j == i || j == (i+1)%l {
					continue
				}
				if InTriangle(pts[currentIndices[j]], pt0, pt1, pt2) {
					// If the triangle includes another point, the triangle is not an ear.
					continue index
				}
			}
			// The angle is less than 180 degrees. This is an ear.
			idx = i
			break
		}
		if idx < 0 {
			// TODO: This happens when there is self-crossing.
			panic(fmt.Sprintf("math: there is no ear in the polygon: %v", pts))
		}
		i0, i1, i2 := adjacentIndices(currentIndices, idx)
		indices = append(indices, i0, i1, i2)
		currentIndices = append(currentIndices[:idx], currentIndices[idx+1:]...)
	}
	return indices
}