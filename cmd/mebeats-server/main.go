// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"embed"
	"flag"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/flamego/flamego"
	"github.com/golang/freetype"
	jsoniter "github.com/json-iterator/go"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/mebeats/report"
)

//go:embed assets/*
var assets embed.FS

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	key := flag.String("key", "", "Server report key")
	flag.Parse()

	heartRate := 0

	// Load fonts.
	fontFile, err := assets.Open("assets/PressStart2P-Regular.ttf")
	if err != nil {
		log.Fatal("Failed to open font: %v", err)
	}
	fontBytes, err := io.ReadAll(fontFile)
	if err != nil {
		log.Fatal("Failed to read font: %v", err)
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal("Failed to parse font: %v", err)
	}

	// Read heart image.
	heartImage, err := assets.Open("assets/heart.png")
	if err != nil {
		log.Fatal("Failed to load asset heart.png: %v", err)
	}
	heart, err := imaging.Decode(heartImage)
	heart = imaging.Resize(heart, 50, 50, imaging.Lanczos)
	if err != nil {
		log.Fatal("Failed to decode image: %v", err)
	}

	f := flamego.New()
	f.Get("/rate.png", func(ctx flamego.Context) {
		background := imaging.New(150, 50, color.NRGBA{})

		textImage := image.NewNRGBA(image.Rect(0, 0, 200, 200))
		text := freetype.NewContext()
		text.SetDPI(72)
		text.SetFont(font)
		text.SetFontSize(26)
		text.SetClip(textImage.Bounds())
		text.SetDst(textImage)
		text.SetSrc(image.NewUniform(color.RGBA{R: 0, G: 0, B: 0, A: 255}))
		pt := freetype.Pt(60, 35+int(text.PointToFixed(26))>>8)
		_, err := text.DrawString(strconv.Itoa(heartRate), pt)
		if err != nil {
			log.Error("Failed to draw string: %v", err)
			ctx.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			return
		}
		background = imaging.Paste(background, textImage, image.Pt(0, 0))
		background = imaging.Paste(background, heart, image.Pt(0, 0))

		encoder := png.Encoder{}
		err = encoder.Encode(ctx.ResponseWriter(), background)
		if err != nil {
			log.Error("Failed to encode image: %v", err)
			ctx.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			return
		}
		ctx.ResponseWriter().WriteHeader(http.StatusOK)
	})

	f.Post("/report", func(ctx flamego.Context) {
		var body report.Body
		err := jsoniter.NewDecoder(ctx.Request().Request.Body).Decode(&body)
		if err != nil {
			log.Error("Failed to parse request body: %v", err)
			ctx.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			return
		}

		if body.Key != *key {
			ctx.ResponseWriter().WriteHeader(http.StatusForbidden)
			return
		}

		heartRate = body.HeartRate
		ctx.ResponseWriter().WriteHeader(http.StatusCreated)
	})

	f.Run()
}
