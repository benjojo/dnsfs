package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func handleUpload(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Query().Get("name") == "" {
		http.Error(rw, "Please supply a file name as ?name=", http.StatusInternalServerError)
		return
	}
	filename := req.URL.Query().Get("name")

	fullfile, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "Unable to read data that was submitted", http.StatusInternalServerError)
		return
	}

	chunkcount := len(fullfile) / 180
	fmt.Printf("%d chunks need to be uploaded...\n", chunkcount)

	var submissionSlice []byte
	for bytePos := 0; bytePos < len(fullfile); bytePos = bytePos + 180 {
		if bytePos+180 > len(fullfile) {
			submissionSlice = fullfile[bytePos:]
		} else {
			submissionSlice = fullfile[bytePos : bytePos+179]
		}

		b64string := base64.StdEncoding.EncodeToString(submissionSlice)

		go uploadChunk(filename, bytePos/180, b64string)
		time.Sleep(time.Millisecond * 100)
	}
}
