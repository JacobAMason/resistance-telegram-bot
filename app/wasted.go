package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

var wastedCommand = BotCommand{
	name:        "wasted",
	description: "Generated a 'wasted' version of a gif",
	matcher:     messageContainsCommandMatcher("wasted"),
	execute: func(bot TeleBot, update Update, respChan chan BotResponse) {

		var fileID string

		if update.Message.ReplyToMessage.Document != nil {
			fileID = update.Message.ReplyToMessage.Document.FileID
		} else if update.Message.ReplyToMessage.Video != nil {
			fileID = update.Message.ReplyToMessage.Video.FileID
		}

		if fileID == "" {
			respChan <- *NewTextBotResponse("Please respond to a gif ", update.Message.Chat.ID)
			return
		}

		wholeCommand := getContentFromCommand(update.Message.Text, "wasted")
		desiredFlash, err := strconv.ParseFloat(wholeCommand, 64)
		if err != nil {
			desiredFlash = .5
		}

		path, err := bot.GetFile(fileID, 2097152)
		if err != nil {
			respChan <- *NewTextBotResponse("Error downloading file: "+err.Error(), update.Message.Chat.ID)
			return
		}

		if len(path) <= 4 || !(path[len(path)-4:] == ".mp4" || path[len(path)-5:] == ".webm" ) {
			respChan <- *NewTextBotResponse(fmt.Sprintf("Unsupported filetype: %s", path), update.Message.Chat.ID)
		}

		frameRate, err := VideoFramerate(path)
		if err != nil {
			respChan <- *NewTextBotResponse(fmt.Sprintf("Error retrieving framerate: %s", err.Error()), update.Message.Chat.ID)
			return
		}

		os.MkdirAll("wastedTemp", os.ModePerm)
		RemoveContents("wastedTemp")
		err = UnstitchImages(path, "wastedTemp")

		if err != nil {
			respChan <- *NewTextBotResponse("Error unstitching file: "+err.Error(), update.Message.Chat.ID)
			return
		}

		fontByteData, err := ioutil.ReadFile("./pricedown bl.ttf")
		font, err := truetype.Parse(fontByteData)
		if err != nil {
			bot.errorReport.Log(err.Error())
			respChan <- *NewTextBotResponse("Error Loading font: "+err.Error(), update.Message.Chat.ID)
			return
		}

		frames, _ := ioutil.ReadDir("wastedTemp")
		savedFrameIndex := 1
		flashed := false
		for i := 0; i < len(frames); i++ {
			filePath := fmt.Sprintf("wastedTemp/%04d.png", i+1)
			im, err := gg.LoadPNG(filePath)
			if err != nil {
				bot.errorReport.Log("unable to load image: " + err.Error())
				return
			}

			dc := gg.NewContextForImage(im)
			Greyscale(dc, math.Min(float64(i)/(float64(len(frames))*desiredFlash), 1))
			if float64(i)/float64(len(frames)) > desiredFlash {
				if flashed == false {
					Shift(dc, 1, 1, 1, 1)
					flashed = true
				} else {
					size := math.Min(float64(dc.Width()), float64(dc.Height()))

					dc.SetRGBA(0, 0, 0, .5)
					dc.DrawRectangle(0, (float64(dc.Height())/2.0)-(size*.2), float64(dc.Width()), size*.25)
					dc.Fill()

					face := truetype.NewFace(font, &truetype.Options{
						Size: size * .25,
					})
					dc.SetFontFace(face)
					dc.SetRGB(1, 0, 0)
					dc.DrawStringAnchored("WASTED", float64(dc.Width()/2)-(size*.4), float64(dc.Height()/2), 0.0, 0.0)
					dc.SavePNG(fmt.Sprintf("wastedTemp/F_%04d.png", savedFrameIndex+1))
					savedFrameIndex++
				}
			}
			dc.SavePNG(fmt.Sprintf("wastedTemp/F_%04d.png", savedFrameIndex+1))
			savedFrameIndex++
		}

		err = StichPicturesTogether("wastedTemp/F_%04d.png", "wastedout.mp4", frameRate)
		if err != nil {
			respChan <- *NewTextBotResponse("Error stitching file: "+err.Error(), update.Message.Chat.ID)
			return
		}

		respChan <- *NewFileUploadBotResponse("wastedout.mp4", update.Message.Chat.ID)

		err = os.Remove(path)
		if err != nil {
			respChan <- *NewTextBotResponse("Error deleting file: "+err.Error(), update.Message.Chat.ID)
			return
		}

	},
}
