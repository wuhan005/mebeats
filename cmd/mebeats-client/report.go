// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

type reportOptions struct {
	HeartRate int
}

func reportToServer(serverAddr string, opts reportOptions) error {
	type report struct {
		HeartRate int `json:"heart_rate"`
	}

	body, err := jsoniter.Marshal(
		report{
			HeartRate: opts.HeartRate,
		},
	)
	if err != nil {
		return errors.Wrap(err, "json marshal")
	}

	url := strings.TrimSuffix(serverAddr, "/") + "/report"
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	defer func() {
		_ = resp.Body.Close()
	}()
	if err != nil {
		return errors.Wrap(err, "post")
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	return nil
}
