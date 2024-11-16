package api

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/nharu-0630/twitter-api-client/model"
	"go.uber.org/zap"
)

func (c *Client) DownloadAllMedia(dir string, tweet model.Tweet) error {
	for i, media := range tweet.Legacy.Entities.Media {
		nameWithoutExt := media.MediaURLHTTPS[strings.LastIndex(media.MediaURLHTTPS, "/")+1:]
		nameWithoutExt = nameWithoutExt[:strings.LastIndex(nameWithoutExt, ".")]
		nameExt := media.MediaURLHTTPS[strings.LastIndex(media.MediaURLHTTPS, "/")+1:]
		nameExt = nameExt[strings.LastIndex(nameExt, ".")+1:]
		if media.Type == "photo" {
			originalURL := "https://pbs.twimg.com/media/" + nameWithoutExt + "?format=" + nameExt + "&name=orig"
			err := c.downloadMediaByOriginalURL(dir, originalURL, tweet.Core.UserResults.Result.Legacy.ScreenName, tweet.RestID, tweet.Legacy.CreatedAt, "img", "img", i+1, nameExt)
			if err != nil {
				zap.L().Error(err.Error())
			}
		} else if media.Type == "video" || media.Type == "animated_gif" {
			err := c.downloadMediaByOriginalURL(dir, media.MediaURLHTTPS, tweet.Core.UserResults.Result.Legacy.ScreenName, tweet.RestID, tweet.Legacy.CreatedAt, "thumb", media.Type, i+1, nameExt)
			if err != nil {
				zap.L().Error(err.Error())
			}
			bitrate := 0
			videoURL := ""
			for _, variant := range media.VideoInfo.Variants {
				if variant.ContentType == "application/x-mpegURL" {
					continue
				}
				if variant.Bitrate > bitrate {
					bitrate = variant.Bitrate
					videoURL = variant.URL
				}
			}
			if videoURL != "" {
				videoNameExt := videoURL[strings.LastIndex(videoURL, "/")+1:]
				videoNameExt = videoNameExt[strings.LastIndex(videoNameExt, ".")+1:]
				videoNameExt = videoNameExt[:strings.LastIndex(videoNameExt, "?")]
				err = c.downloadMediaByOriginalURL(dir, videoURL, tweet.Core.UserResults.Result.Legacy.ScreenName, tweet.RestID, tweet.Legacy.CreatedAt, media.Type, media.Type, i+1, videoNameExt)
				if err != nil {
					zap.L().Error(err.Error())
				}
			}
		}
	}
	return nil
}

func (c *Client) downloadMediaByOriginalURL(dirName string, originalURL string, screenName string, tweetID string, createdAt string, suffix string, dirSuffix string, index int, ext string) error {
	outputDir := os.Getenv("OUTPUT_DIR")
	if outputDir == "" {
		return errors.New("output directory is not set")
	}
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0755)
	}
	outputDir = path.Join(outputDir, dirName)
	parsedCreatedAt, err := time.Parse(time.RubyDate, createdAt)
	if err != nil {
		return err
	}
	createdAt = parsedCreatedAt.Format("20060102_150405")
	outputDir = path.Join(outputDir, screenName+"-"+tweetID+"-"+createdAt+"-"+dirSuffix)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0755)
	}
	fileName := screenName + "-" + tweetID + "-" + createdAt + "-" + suffix + strconv.Itoa(index)
	mediaPath := path.Join(outputDir, fileName+"."+ext)
	if _, err := os.Stat(mediaPath); !os.IsNotExist(err) {
		return errors.New("File already exists: " + mediaPath)
	}
	resp, err := c.client.Get(originalURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Failed to download media: " + originalURL)
	}
	out, err := os.Create(mediaPath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
