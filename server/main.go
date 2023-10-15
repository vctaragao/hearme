package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Logger.SetHeader(`{"level":"${level}}"`)

	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	e.GET("/", handleTrack)

	e.Logger.Fatal(e.Start(":8008"))
}

func handleTrack(c echo.Context) error {
	trackInfo, err := getTrackInfo(c, "Daydream - Soobin Hoang SonThaoboy (Hiderway Remix)")
	if err != nil {
		c.Logger().Fatal(err)
	}

	trackFilePath := "tracks/Daydream - Soobin Hoang SonThaoboy (Hiderway Remix).mp3"
	trackInfo.FileInfo, err = os.Stat(trackFilePath)
	if err != nil {
		c.Logger().Fatal(err)
	}

	c.Response().Header().Set(echo.HeaderContentLength, trackInfo.StrSize())

	file, err := os.Open(trackFilePath)
	if err != nil {
		c.Logger().Fatal(err)
	}
	defer file.Close()

	streamTrack(c, file, trackInfo)

	return nil
}

func getTrackInfo(c echo.Context, songName string) (TrackInfo, error) {
	conf, err := os.ReadFile("./tracks/tracks.json")
	if err != nil {
		c.Logger().Error(err)
	}

	var tracks Tracks
	if err := json.Unmarshal(conf, &tracks); err != nil {
		c.Logger().Error(err)
	}

	track, exists := tracks["Daydream - Soobin Hoang SonThaoboy (Hiderway Remix)"]
	if !exists {
		c.Logger().Error("Track not found")
		return TrackInfo{}, errors.New("Track not found")
	}

	return track, nil
}

func streamTrack(c echo.Context, file *os.File, trackInfo TrackInfo) {
	songRead := int64(0)
	buffer := make([]byte, trackInfo.BytesPerSecond()*Duration)
	for {
		if _, err := file.Seek(songRead, 0); err != nil {
			c.Logger().Fatal(err)
		}

		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			c.Logger().Fatal(err)
		}
		songRead += int64(n)

		if _, err := c.Response().Write(buffer[:n]); err != nil {
			c.Logger().Fatal(err)
		}

		c.Response().Flush()

		if err == io.EOF {
			break
		}

		time.Sleep(time.Second)
	}
}
