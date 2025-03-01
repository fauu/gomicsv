/*
 * Copyright (c) 2013-2021 Utkan Güngördü <utkan@freeconsole.org>
 * Copyright (c) 2021-2025 Piotr Grabowski
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package gomicsv

import (
	"fmt"
	"math"

	"github.com/gotk3/gotk3/gdk"
)

type Color struct {
	R byte `json:"R"`
	G byte `json:"G"`
	B byte `json:"B"`
	A byte `json:"A"`
}

func NewColorFromGdkRGBA(rgba *gdk.RGBA) Color {
	return Color{
		R: colorFloat64ToByte(rgba.GetRed()),
		G: colorFloat64ToByte(rgba.GetGreen()),
		B: colorFloat64ToByte(rgba.GetBlue()),
		A: colorFloat64ToByte(rgba.GetAlpha()),
	}
}

func (color Color) ToCSS() string {
	return fmt.Sprintf("rgba(%d, %d, %d, %d)", color.R, color.G, color.B, color.A)
}

func (color Color) ToPixelInt() uint32 {
	return uint32(color.R)<<24 + uint32(color.G)<<16 + uint32(color.B)<<8 + uint32(color.A)
}

func (color Color) ToGdkRGBA() gdk.RGBA {
	return *gdk.NewRGBA(
		colorByteToFloat64(color.R),
		colorByteToFloat64(color.G),
		colorByteToFloat64(color.B),
		colorByteToFloat64(color.A),
	)
}

func colorFloat64ToByte(f float64) byte {
	return byte(math.Round(f * 255.0))
}

func colorByteToFloat64(b byte) float64 {
	return float64(b) / 255.0
}
