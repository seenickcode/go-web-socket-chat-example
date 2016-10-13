package api

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/jmcvetta/randutil"
	"github.com/juju/errgo"

	log "github.com/Sirupsen/logrus"
)

// WrappedAPIResponse .
type WrappedAPIResponse struct {
	Data interface{} `json:"data"`
}

// NewWrappedAPIResponse .
func NewWrappedAPIResponse(obj interface{}) *WrappedAPIResponse {
	return &WrappedAPIResponse{
		Data: obj,
	}
}

func bytesFromMultiPartFile(file multipart.File) (data []byte, err error) {
	reader := bufio.NewReader(file) // to bufio.Reader
	buffer := bytes.NewBuffer(make([]byte, 0))
	_, err = buffer.ReadFrom(reader) // io.Reader -> bytes.Buffer
	if err != nil {
		err = errgo.Mask(err)
		return
	}
	data = buffer.Bytes()
	return
}

func httpResponseToString(r *http.Response) string {
	data, _ := ioutil.ReadAll(r.Body)
	return string(data)
}

func randomString() string {
	rand, _ := randutil.String(12, randutil.Alphabet)
	ts := time.Now().Unix()
	return fmt.Sprintf("%v%v", rand, ts)
}

func randomInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func randomIntAsString(min, max int) string {
	return strconv.Itoa(randomInt(min, max))
}

func randInt64() int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(100000000000))
}

func renderIfError(w http.ResponseWriter, err error, status int) bool {
	if err != nil {
		renderError(w, err, status)
		return false
	}
	return true
}

func trunc(s string, i int, filler string) string {
	if len(s) < i {
		return s
	}
	if utf8.ValidString(s[:i]) {
		to := i - len(filler)
		return s[:to] + filler
	}
	return s[:i+1] + filler
}

func isErr(err error) bool {
	if err != nil {
		logErr(err)
	}
	return (err != nil)
}

func logErr(err error) {
	log.Error(err)
	debug.PrintStack()
}
