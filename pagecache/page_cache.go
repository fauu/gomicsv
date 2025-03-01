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

package pagecache

import (
	"sync"
	"time"

	"github.com/fauu/gomicsv/util"
	"github.com/gotk3/gotk3/gdk"
)

type KeepReason uint8

const (
	KeepReasonPreload KeepReason = 1 << iota
	KeepReasonJumpmark
)

type PageCache struct {
	Pages      map[int]CachedPage
	pagesMutex sync.Mutex
	keep       map[int]KeepReason
	keepMutex  sync.Mutex
}

func NewPageCache() PageCache {
	return PageCache{
		make(map[int]CachedPage),
		sync.Mutex{},
		make(map[int]KeepReason),
		sync.Mutex{},
	}
}

type CachedPage struct {
	Pixbuf *gdk.Pixbuf
	Time   time.Time
}

func NewCachedPage(pixbuf *gdk.Pixbuf) CachedPage {
	return CachedPage{pixbuf, time.Now()}
}

func (cache *PageCache) Get(i int) (*CachedPage, bool) {
	cache.pagesMutex.Lock()
	defer cache.pagesMutex.Unlock()
	if page, ok := cache.Pages[i]; ok {
		return &page, true
	} else {
		return nil, false
	}
}

func (cache *PageCache) Insert(i int, pixbuf *gdk.Pixbuf, keepReason KeepReason) {
	cache.set(i, pixbuf)
	cache.Keep(i, keepReason)
}

func (cache *PageCache) Keep(i int, reason KeepReason) {
	cache.keepMutex.Lock()
	defer cache.keepMutex.Unlock()
	cache.keep[i] |= reason
}

func (cache *PageCache) DontKeep(i int, reason KeepReason) {
	cache.keepMutex.Lock()
	defer cache.keepMutex.Unlock()
	cache.doDontKeep(i, reason)
}

func (cache *PageCache) DontKeepSlice(indices []int, reason KeepReason) {
	cache.keepMutex.Lock()
	defer cache.keepMutex.Unlock()
	for _, i := range indices {
		cache.doDontKeep(i, reason)
	}
}

func (cache *PageCache) doDontKeep(i int, reason KeepReason) {
	keep, ok := cache.keep[i]
	if !ok {
		return
	}
	keep &^= reason
	if keep == 0 {
		delete(cache.keep, i)
	} else {
		cache.keep[i] = keep
	}
}

func (cache *PageCache) Trim() {
	var toRemove []int
	for i := range cache.Pages {
		if cache.keep[i] == 0 {
			toRemove = append(toRemove, i)
		}
	}
	if len(toRemove) > 0 {
		cache.remove(toRemove)
	}
}

func (cache *PageCache) set(i int, pixbuf *gdk.Pixbuf) {
	cache.pagesMutex.Lock()
	defer cache.pagesMutex.Unlock()
	if pixbuf == nil {
		delete(cache.Pages, i)
	} else {
		cache.Pages[i] = NewCachedPage(pixbuf)
	}
}

func (cache *PageCache) remove(indices []int) {
	cache.pagesMutex.Lock()
	defer cache.pagesMutex.Unlock()
	for _, i := range indices {
		delete(cache.Pages, i)
	}
	util.GC()
}
