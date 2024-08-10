/*
 * Copyright (c) 2013-2021 Utkan Güngördü <utkan@freeconsole.org>
 * Copyright (c) 2021-2024 Piotr Grabowski
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

package util

import (
	"runtime"
)

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Wrap(val, low, mod int) int {
	val %= mod
	if val < low {
		val = mod + val
	}
	return val
}

func Clamp(val, low, high float64) float64 {
	if val < low {
		val = low
	} else if val > high {
		val = high
	}

	return val
}

func Fit(sw, sh, fw, fh int) (int, int) {
	r := float64(sw) / float64(sh)

	var nw, nh float64
	if float64(fw) >= float64(fh)*r {
		nw, nh = float64(fh)*r, float64(fh)
	} else {
		nw, nh = float64(fw), float64(fw)/r
	}
	return int(nw), int(nh)
}

func IntPtrToUintPtr(l *int) *uint {
	if l != nil {
		ul := uint(*l)
		return &ul
	}
	return nil
}

func GC() {
	runtime.GC()
	runtime.GC()
}
