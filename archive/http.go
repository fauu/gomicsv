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

package archive

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/gotk3/gotk3/gdk"

	"github.com/fauu/gomicsv/pagecache"
	"github.com/fauu/gomicsv/pixbuf"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36"
)

type HTTP struct {
	urlTemplate          string
	referer              string
	firstPageOffset      int
	pageCache            *pagecache.PageCache
	fetchInProgress      map[int]bool
	fetchInProgressMutex sync.Mutex
}

func NewHTTP(url string, pageCache *pagecache.PageCache, referer string) (*HTTP, error) {
	if !isPageURLTemplate(url) {
		var ok bool
		url, ok = tryMakeURLTemplateFromSampleURL(url)
		if !ok {
			return nil, errors.New("Couldn't determine the URL template from the sample URL")
		}
	}

	newHTTP := HTTP{
		urlTemplate:          url,
		referer:              referer,
		pageCache:            pageCache,
		fetchInProgress:      make(map[int]bool),
		fetchInProgressMutex: sync.Mutex{},
	}

	firstPageIdx := 0
	var firstPixbuf *gdk.Pixbuf
	for firstPageIdx < 2 {
		var err error
		firstPixbuf, err = newHTTP.downloadImage(firstPageIdx, false)
		if err != nil {
			log.Printf("First image not located at index %d", firstPageIdx)
		} else {
			break
		}
		firstPageIdx++
	}
	if firstPixbuf == nil {
		return nil, errors.New("Couldn't locate the first image")
	}

	newHTTP.firstPageOffset = firstPageIdx
	pageCache.Insert(0, firstPixbuf, pagecache.KeepReasonPreload)

	return &newHTTP, nil
}

func (ar *HTTP) Load(i int, autorotate bool, nPreload int) (*gdk.Pixbuf, error) {
	var pixbuf *gdk.Pixbuf
	var err error
	cached, isCached := ar.pageCache.Get(i)
	if !isCached {
		if downloading := ar.getAndSetPageFetchInProgress(i, true); downloading {
			// Wait until downloading done
			var tries = 10
			for range time.Tick(time.Millisecond * 500) {
				cached, isCached = ar.pageCache.Get(i)
				if isCached {
					break
				}
				if tries <= 0 {
					return nil, fmt.Errorf("Ran out of tries when waiting for page %d to download", i)
				}
				tries--
			}
		} else {
			pixbuf, err = ar.downloadImage(i+ar.firstPageOffset, autorotate)
			ar.setPageFetchInProgress(i, false)
			if err != nil {
				return nil, err
			}
		}
	}
	if pixbuf == nil {
		pixbuf = cached.Pixbuf
	}

	if i == 0 {
		pixbuf, err = pixbuf.ApplyEmbeddedOrientation()
		if err != nil {
			return nil, err
		}
	}

	preloadStart := i - nPreload
	preloadEnd := i + nPreload
	for j := preloadStart; j <= preloadEnd; j++ {
		if j < 0 || j == i {
			continue
		}
		if _, ok := ar.pageCache.Get(j); !ok {
			if downloading := ar.getAndSetPageFetchInProgress(j, true); !downloading {
				go func(k int) {
					pixbuf, err = ar.downloadImage(k+ar.firstPageOffset, autorotate)
					ar.setPageFetchInProgress(k, false)
					if err != nil {
						log.Printf("Couldn't preload image: %v", err)
					}
					ar.pageCache.Insert(k, pixbuf, pagecache.KeepReasonPreload)
				}(j)
			}
		}
	}

	if !isCached {
		ar.pageCache.Insert(i, pixbuf, pagecache.KeepReasonPreload)
	}

	return pixbuf, nil
}

func (ar *HTTP) Kind() Kind {
	return HTTPKind
}

func (ar *HTTP) ArchiveName() string {
	return ar.urlTemplate
}

func (ar *HTTP) Name(i int) (string, error) {
	return ar.urlTemplate, nil
}

func (ar *HTTP) Len() *int {
	// Length unknown
	return nil
}

func (ar *HTTP) Close() error {
	return nil
}

func (ar *HTTP) getAndSetPageFetchInProgress(i int, value bool) bool {
	ar.fetchInProgressMutex.Lock()
	defer ar.fetchInProgressMutex.Unlock()
	_, ok := ar.fetchInProgress[i]
	ar.doSetPageFetchInProgress(i, value)
	return ok
}

func (ar *HTTP) setPageFetchInProgress(i int, value bool) {
	ar.fetchInProgressMutex.Lock()
	defer ar.fetchInProgressMutex.Unlock()
	ar.doSetPageFetchInProgress(i, value)
}

func (ar *HTTP) doSetPageFetchInProgress(i int, value bool) {
	if value {
		ar.fetchInProgress[i] = true
	} else {
		delete(ar.fetchInProgress, i)
	}
}

func (ar *HTTP) downloadImage(i int, autorotate bool) (*gdk.Pixbuf, error) {
	url := fmt.Sprintf(ar.urlTemplate, i)
	headers := map[string]string{"User-Agent": userAgent}
	if ar.referer != "" {
		headers["Referer"] = ar.referer
	}
	res, err := httpGet(url, headers)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	pixbuf, err := pixbuf.Load(res.Body, autorotate)
	if err != nil {
		return nil, err
	}

	return pixbuf, nil
}

var pagePlaceholderRegexp = regexp.MustCompile(`%\d{0,2}d`)

func isPageURLTemplate(s string) bool {
	return pagePlaceholderRegexp.MatchString(s)
}

var pageParamRegexp = regexp.MustCompile(`(?:\/(\d{1,3})(?:\/|[^\d\w]))|(?:(\d{1,3})\.(?:jp|png|gif))`)

func tryMakeURLTemplateFromSampleURL(url string) (template string, ok bool) {
	matches := pageParamRegexp.FindAllStringSubmatchIndex(url, -1)
	if len(matches) < 1 {
		return "", false
	}
	lastMatch := matches[len(matches)-1]
	if len(lastMatch) <= 2 {
		return "", false
	}
	for i := 2; i < len(lastMatch); i += 2 {
		if lastMatch[i] != -1 {
			start, end := lastMatch[i], lastMatch[i+1]
			return url[:start] + "%d" + url[end:], true
		}
	}
	return "", false
}

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func httpGet(reqURL string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	log.Printf("GET %s", reqURL)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing request: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Got status code: %d %s", res.StatusCode, res.Status)
	}

	return res, nil
}
