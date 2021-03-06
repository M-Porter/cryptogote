package main

import (
	"crypto/rand"
	"html/template"
	"log"
	"net/http"

	"bytes"

	b64 "encoding/base64"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/unrolled/render"
)

// Pages struct ...
type Pages struct {
	Render *render.Render
	DB     *gorm.DB
}

var noteLink = "noteLink"

// NewMessageHandler ...
func (pages Pages) NewMessageHandler(w http.ResponseWriter, r *http.Request) {
	pages.Render.HTML(w, http.StatusOK, "messages_new", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

var encoding = b64.URLEncoding.WithPadding(-1)

// PostCryptoMessageHandler ...
func (pages Pages) PostCryptoMessageHandler(w http.ResponseWriter, r *http.Request) {
	noteContent := []byte(r.FormValue("Note.Content"))

	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal(err)
	}

	cipherText, err := Encrypt(key, noteContent)
	if err != nil {
		log.Fatal(err)
	}

	pages.DB.Create(&Notes{
		Content:       encoding.EncodeToString(cipherText),
		EncryptionKey: encoding.EncodeToString(key),
	})

	log.Println(r.URL.Scheme)

	scheme := r.URL.Scheme
	if scheme == "" {
		scheme = "http://"
	} else {
		scheme = "https://"
	}

	var buffer bytes.Buffer
	buffer.WriteString(scheme)
	buffer.WriteString(r.Host)
	buffer.WriteString("/messages/")
	buffer.WriteString(encoding.EncodeToString(key))

	pages.Render.HTML(w, http.StatusOK, "messages_new", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		noteLink:         buffer.String(),
	})
}

// ShowMessageHandler ...
func (pages Pages) ShowMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	keyStr := vars["key"]

	var note Notes
	if pages.DB.First(&note, "encryption_key = ? ", keyStr).RecordNotFound() {
		// No record in DB
		http.Redirect(w, r, "/", 302)
	}

	cipherText, err := encoding.DecodeString(note.Content)
	if err != nil {
		log.Fatal(err)
		http.Redirect(w, r, "/", 302)
		return
	}

	key, err := encoding.DecodeString(keyStr)
	if err != nil {
		log.Fatal(err)
		http.Redirect(w, r, "/", 302)
		return
	}

	decrypted, err := Decrypt(key, cipherText)
	// Decryption failed. Usually means key is invalid.
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	unsafe := blackfriday.MarkdownCommon(decrypted)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	// Update note with filler values
	pages.DB.Model(&note).Update("content", "garbage")
	pages.DB.Model(&note).Update("encryption_key", gorm.Expr("NULL"))

	// Delete note
	pages.DB.Delete(&note)

	pages.Render.HTML(w, http.StatusOK, "messages_show", map[string]interface{}{
		"content": template.HTML(string(html[:])),
	})
}

// StatisticsHandler ...
func (pages Pages) StatisticsHandler(w http.ResponseWriter, r *http.Request) {
	pages.Render.HTML(w, http.StatusOK, "statistics_index", nil)
}
