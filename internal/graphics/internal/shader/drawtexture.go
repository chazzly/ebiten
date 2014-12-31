// Copyright 2014 Hajime Hoshi
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

package shader

import (
	"github.com/go-gl/gl"
	"github.com/hajimehoshi/ebiten/internal/opengl"
)

func glMatrix(m [4][4]float64) [16]float32 {
	result := [16]float32{}
	for j := 0; j < 4; j++ {
		for i := 0; i < 4; i++ {
			result[i+j*4] = float32(m[i][j])
		}
	}
	return result
}

type Matrix interface {
	Element(i, j int) float64
}

type TextureQuads interface {
	Len() int
	Vertex(i int) (x0, y0, x1, y1 float32)
	Texture(i int) (u0, v0, u1, v1 float32)
}

var initialized = false

// TODO: Use unsafe.SizeOf?
const uint16Size = 2
const float32Size = 4

func DrawTexture(c *opengl.Context, texture opengl.Texture, projectionMatrix [4][4]float64, quads TextureQuads, geo Matrix, color Matrix) error {
	// TODO: Check len(quads) and gl.MAX_ELEMENTS_INDICES?
	const stride = 4 * 4
	if !initialized {
		if err := initialize(c); err != nil {
			return err
		}
		initialized = true
	}

	if quads.Len() == 0 {
		return nil
	}

	// TODO: Check performance
	program := useProgramColorMatrix(glMatrix(projectionMatrix), geo, color)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.Texture(texture).Bind(gl.TEXTURE_2D)

	vertexAttrLocation := program.GetAttributeLocation("vertex")
	texCoordAttrLocation := program.GetAttributeLocation("tex_coord")

	vertexAttrLocation.EnableArray()
	texCoordAttrLocation.EnableArray()
	defer func() {
		texCoordAttrLocation.DisableArray()
		vertexAttrLocation.DisableArray()
	}()

	vertexAttrLocation.AttribPointer(stride, uintptr(float32Size*0))
	texCoordAttrLocation.AttribPointer(stride, uintptr(float32Size*2))

	vertices := []float32{}
	for i := 0; i < quads.Len(); i++ {
		x0, y0, x1, y1 := quads.Vertex(i)
		u0, v0, u1, v1 := quads.Texture(i)
		vertices = append(vertices,
			x0, y0, u0, v0,
			x1, y0, u1, v0,
			x0, y1, u0, v1,
			x1, y1, u1, v1,
		)
	}
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, float32Size*len(vertices), vertices)
	gl.DrawElements(gl.TRIANGLES, 6*quads.Len(), gl.UNSIGNED_SHORT, uintptr(0))

	gl.Flush()
	return nil
}