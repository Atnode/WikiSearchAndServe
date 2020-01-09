package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	zim "github.com/akhenakh/gozim"
	"github.com/blevesearch/bleve"
)

const (
	ArticlesPerPage = 16
)

func cacheLookup(url string) (*CachedResponse, bool) {
	if v, ok := cache.Get(url); ok {
		c := v.(CachedResponse)
		return &c, ok
	}
	return nil, false
}

// dealing with cached response, responding directly
func handleCachedResponse(cr *CachedResponse, w http.ResponseWriter, r *http.Request) {
	if cr.ResponseType == RedirectResponse {
		log.Printf("302 from %s to %s\n", r.URL.Path, "zim/"+string(cr.Data))
		http.Redirect(w, r, "/zim/"+string(cr.Data), http.StatusMovedPermanently)
	} else if cr.ResponseType == NoResponse {
		log.Printf("404 %s\n", r.URL.Path)
		http.NotFound(w, r)
	} else if cr.ResponseType == DataResponse {
		log.Printf("200 %s\n", r.URL.Path)
		w.Header().Set("Content-Type", cr.MimeType)
		// 15 days
		w.Header().Set("Cache-control", "public, max-age=1350000")
		w.Write(cr.Data)
	}
}

// the handler receiving http request
func zimHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[5:]
	// lookup in the cache for a cached response
	if cr, iscached := cacheLookup(url); iscached {
		handleCachedResponse(cr, w, r)
		return

	} else {
		var a *zim.Article
		a, _ = Z.GetPageNoIndex(url)

		if a == nil {
			cache.Add(url, CachedResponse{ResponseType: NoResponse})
		} else if a.EntryType == zim.RedirectEntry {
			ridx, err := a.RedirectIndex()
			if err != nil {
				cache.Add(url, CachedResponse{ResponseType: NoResponse})
			} else {
				ra, err := Z.ArticleAtURLIdx(ridx)
				if err != nil {
					cache.Add(url, CachedResponse{ResponseType: NoResponse})
				} else {
					cache.Add(url, CachedResponse{
						ResponseType: RedirectResponse,
						Data:         []byte(ra.FullURL())})
				}
			}
		} else {
			data, err := a.Data()
			if err != nil {
				cache.Add(url, CachedResponse{ResponseType: NoResponse})
			} else {
				cache.Add(url, CachedResponse{
					ResponseType: DataResponse,
					Data:         data,
					MimeType:     a.MimeType(),
				})
			}
		}

		// look again in the cache for the same entry
		if cr, iscached := cacheLookup(url); iscached {
			handleCachedResponse(cr, w, r)
		}
	}
}

func removeExt(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	pageString := r.FormValue("page")
	pageNumber, _ := strconv.Atoi(pageString)
	previousPage := pageNumber - 1
	if pageNumber == 0 {
		previousPage = 0
	}
	nextPage := pageNumber + 1

	q := r.FormValue("search_data")
	d := map[string]interface{}{
		"Query":        q,
		"Filename":     removeExt(*zimPath),
		"Page":         pageNumber,
		"PreviousPage": previousPage,
		"NextPage":     nextPage,
	}

	if q == "" {
		templates["index"].Execute(w, d)
		return
	}
	itemCount := 20
	from := itemCount * pageNumber
	query := bleve.NewQueryStringQuery(q)
	search := bleve.NewSearchRequestOptions(query, itemCount, from, false)
	search.Fields = []string{"Title"}

	sr, err := index.Search(search)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if (20 * (pageNumber + 1)) > int(sr.Total) {
		d["Next"] = 0
	} else {
		d["Next"] = 1
	}

	if pageNumber < 1 {
		d["Previous"] = 0
	} else {
		d["Previous"] = 1
	}

	if sr.Total > 0 {
		d["Nbresult"] = fmt.Sprintf("%d", sr.Total)
		d["Namequery"] = fmt.Sprintf("[%s]", q)
		d["Searchtime"] = fmt.Sprintf("%s", sr.Took)

		// Constructs a list of Hits
		var l []map[string]string

		for _, h := range sr.Hits {
			idx, err := strconv.Atoi(h.ID)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			a, err := Z.ArticleAtURLIdx(uint32(idx))
			if err != nil {
				continue
			}
			l = append(l, map[string]string{
				"Score": strconv.FormatFloat(h.Score, 'f', 1, 64),
				"Title": a.Title,
				"URL":   "/zim/" + a.FullURL()})

		}
		d["Hits"] = l

	} else {
		d["Nbresult"] = fmt.Sprintf("%d", sr.Total)
		d["Namequery"] = fmt.Sprintf("[%s]", q)
		d["Searchtime"] = fmt.Sprintf("%s", sr.Took)
		d["Hits"] = 0
	}

	templates["searchResult"].Execute(w, d)
}

func robotHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "User-agent: *\nDisallow: /\n")
}
