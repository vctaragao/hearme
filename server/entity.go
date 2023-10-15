package main

import (
	"os"
	"strconv"
)

type (
	Track  string
	Tracks map[Track]TrackInfo

	TrackInfo struct {
		Format string `json:"format"`
		Length int64  `json:"length"`
		os.FileInfo
	}
)

func (t *TrackInfo) StrSize() string {
	return strconv.Itoa(int(t.FileInfo.Size()))
}

func (t *TrackInfo) BytesPerSecond() int64 {
	return t.FileInfo.Size() / t.Length
}
