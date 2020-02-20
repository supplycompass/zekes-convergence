package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Image struct {
	Filename 	string
	ContentType string
	Data		[]byte
	Size		int
}

func (i *Image) DataURI() string {
	return fmt.Sprintf("data:%s;base64,%s", i.ContentType, base64.StdEncoding.EncodeToString(i.Data))
}

func (i *Image) ThumbnailJPEG(width int, height int, quality int) (*Image, error) {
	return ThumbnailJPEG(i, width, height, quality)
}

func okContentType(contentType string) bool {
	return contentType == "image/png" || contentType == "image/jpeg" || contentType == "image/gif"
}

func (i *Image) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", i.ContentType)
	w.Header().Set("Content-Length", strconv.Itoa(i.Size))
	w.Write(i.Data)
}

func Convert(r *http.Request, field string) (*Image, error) {
	file, info, err := r.FormFile(field)
	if err != nil {
		return nil, err
	}

	contentType := info.Header.Get("Content-Type")

	if !okContentType(contentType) {
		return nil, errors.New(fmt.Sprintf("Wrong content type: %s", contentType))
	}

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	_, _, err = image.Decode(bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}

	i := &Image {
		Filename:		info.Filename,
		ContentType:	contentType,
		Data:			bs,
		Size:			len(bs),
	}

	return i, nil
}

func ThumbnailJPEG(i *Image, width int, height int, quality int) (*Image, error) {
	img, _, err := image.Decode(bytes.NewReader(i.Data))

	thumbnail := resize.Thumbnail(uint(width), uint(height), img, resize.Lanczos3)

	data := new(bytes.Buffer)
	err = jpeg.Encode(data, thumbnail, &jpeg.Options{
		Quality: quality,
	})

	if err != nil {
		return nil, err
	}

	bs := data.Bytes()

	t := &Image{
		Filename:    "thumbnail.jpg",
		ContentType: "image/jpeg",
		Data:        bs,
		Size:        len(bs),
	}

	return t, nil
}

func ThumbnailPNG(i *Image, width int, height int) (*Image, error) {
	img, _, err := image.Decode(bytes.NewReader(i.Data))

	thumbnail := resize.Thumbnail(uint(width), uint(height), img, resize.Lanczos3)

	data := new(bytes.Buffer)
	err = png.Encode(data, thumbnail)

	if err != nil {
		return nil, err
	}

	bs := data.Bytes()

	t := &Image{
		Filename:    "thumbnail.png",
		ContentType: "image/png",
		Data:        bs,
		Size:        len(bs),
	}

	return t, nil
}
