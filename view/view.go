package view

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Xjs/gopher-rating/model"

	"github.com/Xjs/gopher-rating/storage"

	"github.com/gorilla/mux"
)

// A Handler displays a Gopher webpage
type Handler struct {
	storer   storage.Interface
	template *template.Template
}

// NewHandler initialises a new Handler
func NewHandler(storer storage.Interface, template *template.Template) *Handler {
	return &Handler{storer, template}
}

// TemplateData is data used in the templates
type TemplateData struct {
	Gophers []*model.Gopher
	Ratings []int
	Errors  []error
}

const count = 200

// Gopher shows a single gopher and rating links
func (s *Handler) Gopher(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	var hash model.Hash
	if n, err := hex.Decode(hash[:], []byte(v["hash"])); err != nil || n != sha256.Size {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request")
		return
	}

	gopher, err := s.storer.Load(r.Context(), hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error loading list: %v", err)
		return
	}

	if gopher == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No such gopher")
		return
	}

	rating, err := s.storer.Rating(r.Context(), hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting rating: %v", err)
		return
	}

	if err := s.template.ExecuteTemplate(w, "gopher.html", &TemplateData{
		Gophers: []*model.Gopher{gopher},
		Ratings: []int{rating},
		Errors:  nil,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
	}
}

// Gophers shows a lot of gophers and a upload form
func (s *Handler) Gophers(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	var start int

	if s, err := strconv.ParseInt(v["start"], 10, 64); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request")
		return
	} else {
		start = int(s)
	}

	hashes, err := s.storer.List(r.Context(), start, count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error loading list: %v", err)
		return
	}

	gophers := make([]*model.Gopher, len(hashes))
	ratings := make([]int, len(hashes))
	var errs []error
	for i, hash := range hashes {
		var err error
		gophers[i], err = s.storer.Load(r.Context(), hash)
		if err != nil {
			errs = append(errs, err)
		}
		ratings[i], err = s.storer.Rating(r.Context(), hash)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(gophers) == 0 && len(errs) != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		for _, err := range errs {
			fmt.Fprintf(w, "%v\n", err)
		}
		return
	}

	if err := s.template.ExecuteTemplate(w, "gophers.html", &TemplateData{
		Gophers: gophers,
		Ratings: ratings,
		Errors:  errs,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
	}
}

// Rate rates a single gopher
func (s *Handler) Rate(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	var rating int
	if s, err := strconv.ParseInt(v["start"], 10, 64); err != nil || s < 1 || s > 5 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request")
		return
	} else {
		rating = int(s)
	}

	var hash model.Hash
	if n, err := hex.Decode(hash[:], []byte(v["hash"])); err != nil || n != sha256.Size {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request")
		return
	}

	if err := s.storer.Rate(r.Context(), hash, rating); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error rating: %v", err)
		return
	}

	if err := s.template.ExecuteTemplate(w, "rated.html", nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
	}
}

const maxMem = 1024 * 1024 * 20 // 20 MiB

// Upload uploads a single gopher
func (s *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxMem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad gopher data")
	}

	file, header, err := r.FormFile("gopher")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad gopher data")
		return
	}

	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No image")
		return
	}

	raw, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Couldn't read gopher")
		return
	}

	if err := s.storer.Save(r.Context(), model.NewGopher(raw)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Couldn't save gopher")
		return
	}

	if err := s.template.ExecuteTemplate(w, "saved.html", nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
	}
}
