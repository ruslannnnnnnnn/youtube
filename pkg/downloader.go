package pkg

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/spf13/pflag"
	"golang.org/x/net/http/httpproxy"

	"github.com/kkdai/youtube/v2"
	ytdl "github.com/kkdai/youtube/v2/downloader"
)

var (
	InsecureSkipVerify bool   // skip TLS server validation
	OutputQuality      string // itag number or quality string
	Mimetype           string
	Language           string
	Downloader         *ytdl.Downloader
)

func AddVideoSelectionFlags(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(&OutputQuality, "quality", "q", "medium", "The itag number or quality label (hd720, medium)")
	flagSet.StringVarP(&Mimetype, "Mimetype", "m", "", "Mime-Type to filter (mp4, webm, av01, avc1) - applicable if --quality used is quality label")
	flagSet.StringVarP(&Language, "Language", "l", "", "Language to filter")
}

func GetDownloader() *ytdl.Downloader {
	if Downloader != nil {
		return Downloader
	}

	proxyFunc := httpproxy.FromEnvironment().ProxyFunc()
	httpTransport := &http.Transport{
		// Proxy: http.ProxyFromEnvironment() does not work. Why?
		Proxy: func(r *http.Request) (uri *url.URL, err error) {
			return proxyFunc(r.URL)
		},
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	youtube.SetLogLevel("info")

	if InsecureSkipVerify {
		youtube.Logger.Info("Skip server certificate verification")
		httpTransport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	Downloader = &ytdl.Downloader{
		OutputDir: "/tmp",
	}
	Downloader.HTTPClient = &http.Client{Transport: httpTransport}

	return Downloader
}

func GetVideoWithFormat(videoID string) (*youtube.Video, *youtube.Format, error) {
	dl := GetDownloader()
	video, err := dl.GetVideo(videoID)
	if err != nil {
		return nil, nil, err
	}

	itag, _ := strconv.Atoi(OutputQuality)
	formats := video.Formats

	if Language != "" {
		formats = formats.Language(Language)
	}
	if Mimetype != "" {
		formats = formats.Type(Mimetype)
	}
	if OutputQuality != "" {
		formats = formats.Quality(OutputQuality)
	}
	if itag > 0 {
		formats = formats.Itag(itag)
	}
	if formats == nil {
		return nil, nil, fmt.Errorf("unable to find the specified format")
	}

	formats.Sort()

	// select the first format
	return video, &formats[0], nil
}
