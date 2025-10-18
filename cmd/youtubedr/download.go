package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/ruslannnnnnnnn/youtube/v2/pkg"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "Downloads a video from youtube",
	Example: `youtubedr -o "Campaign Diary".mp4 https://www.youtube.com/watch\?v\=XbNghLqsVwU`,
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		exitOnError(download(args[0]))
	},
}

var (
	ffmpegCheck error
	outputFile  string
	OutputDir   string
)

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&outputFile, "filename", "o", "", "The output file, the default is genated by the video title.")
	downloadCmd.Flags().StringVarP(&OutputDir, "directory", "d", ".", "The output directory.")
	pkg.addVideoSelectionFlags(downloadCmd.Flags())
}

func download(id string) error {
	video, format, err := pkg.getVideoWithFormat(id)
	if err != nil {
		return err
	}

	log.Println("download to directory", OutputDir)

	if strings.HasPrefix(pkg.outputQuality, "hd") {
		if err := checkFFMPEG(); err != nil {
			return err
		}
		return pkg.downloader.DownloadComposite(context.Background(), outputFile, video, pkg.outputQuality, pkg.mimetype, pkg.language)
	}

	return pkg.downloader.Download(context.Background(), video, format, outputFile)
}

func checkFFMPEG() error {
	fmt.Println("check ffmpeg is installed....")
	if err := exec.Command("ffmpeg", "-version").Run(); err != nil {
		ffmpegCheck = fmt.Errorf("please check ffmpegCheck is installed correctly")
	}

	return ffmpegCheck
}
