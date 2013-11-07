package scrapegoat

import (
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func readResponse(resp *http.Response) string {
	reader := getReader(resp)
	r, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return string(r)
}

func getCharset(header http.Header) (string, error) {
	char := header["Content-Type"][0]
	idx := strings.Index(char, "charset")
	if idx >= 0 {
		return char[idx+8 : len(char)], nil
	}
	return "", errors.New("Charset header not found in " + char)
}

func getReader(resp *http.Response) io.Reader {
	cs, err := getCharset(resp.Header)
	if err != nil {
		return resp.Body
	}

	r, err := charset.NewReader(cs, resp.Body)
	if err != nil {
		panic("Charset error " + cs)
	}
	return r
}
