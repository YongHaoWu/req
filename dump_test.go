package req

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestDumpText(t *testing.T) {
	SetFlags(LstdFlags | Lcost)
	reqBody := "request body"
	respBody := "response body"
	reqHeader := "Request-Header"
	respHeader := "Response-Header"
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(respHeader, "req")
		w.Write([]byte(respBody))
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	header := Header{
		reqHeader: "hello",
	}
	resp, err := Post(ts.URL, header, reqBody)
	if err != nil {
		t.Fatal(err)
	}
	dump := resp.Dump()
	for _, keyword := range []string{reqBody, respBody, reqHeader, respHeader} {
		if !strings.Contains(dump, keyword) {
			t.Errorf("dump missing part, want: %s", keyword)
		}
	}
}

func TestDumpUpload(t *testing.T) {
	SetFlags(LreqBody)
	file1 := ioutil.NopCloser(strings.NewReader("file1"))
	uploads := []FileUpload{
		{
			FileName:  "1.txt",
			FieldName: "media",
			File:      file1,
		},
	}
	ts := newDefaultTestServer()
	r, err := Post(ts.URL, uploads, Param{"hello": "req"})
	if err != nil {
		t.Fatal(err)
	}
	dump := r.Dump()
	contains := []string{
		`Content-Disposition: form-data; name="hello"`,
		`Content-Disposition: form-data; name="media"; filename="1.txt"`,
	}
	for _, contain := range contains {
		if !strings.Contains(dump, contain) {
			t.Errorf("multipart dump should contains: %s", contain)
		}
	}
}

func TestDumpWithPrintFunc(t *testing.T) {
	reqBody := "request body"
	p := Param{
		"name": "roc",
		"job":  "programmer",
	}
	buf := bytes.NewBufferString(reqBody)
	ts := newDefaultTestServer()
	Debug = true
	PrintFunc = log.Infoln
	r, err := Post(ts.URL, p, buf)
	if err != nil {
		t.Fatal(err)
	}
	if r.Request().URL.Query().Get("name") != "roc" {
		t.Error("param should in the url when set body manually")
	}
	if string(r.reqBody) != reqBody {
		t.Error("request body not equal")
	}
}
