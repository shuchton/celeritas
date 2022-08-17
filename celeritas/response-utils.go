package celeritas

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
)

func (c *Celeritas) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1024 bytea * 1024 = 1 MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body has multiple json values and may only contain one")
	}
	return nil
}

func (c *Celeritas) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (c *Celeritas) WriteXML(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_, err = w.Write(out)

	return err
}

func (c *Celeritas) DownloadFile(w http.ResponseWriter, r *http.Request, pathToFile, filename string) error {
	fp := path.Join(pathToFile, filename)
	fileToServe := filepath.Clean(fp)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; file=\"%s\"", filename))
	http.ServeFile(w, r, fileToServe)

	return nil
}

func (c *Celeritas) Error404(w http.ResponseWriter, r *http.Request) {
	c.ErrorStatus(w, http.StatusNotFound)
}

func (c *Celeritas) Error500(w http.ResponseWriter, r *http.Request) {
	c.ErrorStatus(w, http.StatusInternalServerError)
}

func (c *Celeritas) ErrorUnauthorized(w http.ResponseWriter, r *http.Request) {
	c.ErrorStatus(w, http.StatusUnauthorized)
}

func (c *Celeritas) ErrorForbidden(w http.ResponseWriter, r *http.Request) {
	c.ErrorStatus(w, http.StatusForbidden)
}

func (c *Celeritas) ErrorStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
