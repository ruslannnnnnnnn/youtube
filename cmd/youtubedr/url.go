package main

import (
	"fmt"

	"github.com/ruslannnnnnnnn/youtube/v2/pkg"
	"github.com/spf13/cobra"
)

// urlCmd represents the url command
var urlCmd = &cobra.Command{
	Use:   "url",
	Short: "Only output the stream-url to desired video",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		video, format, err := pkg.getVideoWithFormat(args[0])
		exitOnError(err)

		url, err := pkg.downloader.GetStreamURL(video, format)
		exitOnError(err)

		fmt.Println(url)
	},
}

func init() {
	pkg.addVideoSelectionFlags(urlCmd.Flags())
	rootCmd.AddCommand(urlCmd)
}
