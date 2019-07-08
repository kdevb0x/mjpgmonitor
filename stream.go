// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package mjpgmonitor

import (
	"image"
	_ "image/jpeg"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

type streamDecoder struct {
	r *multipart.Reader
	i chan image.Image
}

func NewStreamDecoder(url string, user, password string) (*streamDecoder, error) {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	_, param, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	s := &streamDecoder{
		r: multipart.NewReader(resp.Body, strings.Trim(param["boundary"], "-")),
		i: make(chan image.Image),
	}
	return s, nil
}

func (s *streamDecoder) decodeLoop() {
	for {
		np, err := s.r.NextPart()
		if err != nil {
			break
		}

		img, _, err := image.Decode(np)
		if err != nil {
			break
		}
		go func() {
			// goroutine exits when img == nil
			for img != nil {
				s.i <- img
			}
		}()
	}
	return
}
