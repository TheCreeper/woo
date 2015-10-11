package wuu

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/TheCreeper/wuu/wuu/urigen"
	"github.com/TheCreeper/wuu/wuu/verbs"
	"github.com/syndtr/goleveldb/leveldb"
)

const html = `
<html>
<body>
<pre>
wuu(1)                                wuu                                wuu(1)

NAME
    wuu: command line pastebin.

SYNOPSIS
    command | curl -F 'wuu=<-' http://{{ .BaseURL }}

DESCRIPTION
    Just post any data to this server and you'll get a "paste" link.

EXAMPLES
    cat bin/ching | curl -F 'wuu=<-' http://{{ .BaseURL }}
    curl http://{{ .BaseURL }}

SEE ALSO
    http://github.com/TheCreeper/wuu
</pre>
</body>
</html>`

// session is a wrapper around the leveldb.DB object.
type session struct{ *leveldb.DB }

// RootHandler returns a handler that serves HTTP GET requests with the
// compiled html template. If req.RequestURI is longer than 0 characters then
// its assumed that the client is requesting a paste.
func (s session) RootHandler(w http.ResponseWriter, req *http.Request) {
	// Check if the client is requesting a paste. We don't want to serve
	// any requests with a uri longer than 4 characters.
	uri := req.RequestURI[1:]
	if len(uri) != 0 && len(uri) <= 4 {
		// Copy the contents of the paste from the database to
		// memory.
		paste, err := s.Get([]byte(uri), nil)
		if err != nil {
			http.Error(w,
				http.StatusText(http.StatusNotFound),
				http.StatusNotFound)
			return
		}

		// Write out the paste bytes. This is kinda bad since its
		// copying upto 1MB to memory on every request. It might
		// be better if the database interface returned an io.Reader.
		if _, err = w.Write(paste); err != nil {
			http.Error(w,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}
		return
	}

	tmpl, err := template.New("index.html").Parse(html)
	if err != nil {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	data := struct{ BaseURL string }{req.Host}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}
}

// UploadHandler returns a request handler that stores the contents of the
// multipart form in the request as a "paste". The request is replied to with
// the URL of the paste.
func (s session) UploadHandler(w http.ResponseWriter, req *http.Request) {
	// Parse the paste with a max size of 1MB.
	if err := req.ParseMultipartForm(1048576); err != nil {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	// Generate a random name for the paste.
	// A string of 4 characters has about 500k possible combinations.
	pname, err := urigen.Generate(4)
	if err != nil {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	// Check if the form contains a paste.
	val, ok := req.MultipartForm.Value["wuu"]
	if !ok {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	// Store the paste in the database.
	if err = s.Put(pname, []byte(val[0]), nil); err != nil {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	// Return the url of the paste.
	uri := fmt.Sprintf("http://%s/%s\n", req.Host, pname)
	if _, err := w.Write([]byte(uri)); err != nil {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}
}

// Listen opens the database at the filepath dbname and listens on the TCP
// network address addr. RootHandler and UploadHandler are registered as
// handlers for Get and Post requests respectively.
func Listen(dbname, addr string) (err error) {
	db, err := leveldb.OpenFile(dbname, nil)
	if err != nil {
		return
	}
	defer db.Close()

	s := session{db}
	mux := http.NewServeMux()
	mux.Handle("/", verbs.Verbs{
		Get:  s.RootHandler,
		Post: s.UploadHandler})
	return http.ListenAndServe(addr, mux)
}
