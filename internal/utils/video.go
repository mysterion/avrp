package utils

import (
	"strings"
)

var video = []string{
	"3g2",
	"3gp",
	"aaf",
	"asf",
	"avchd",
	"avi",
	"drc",
	"flv",
	"m2v",
	"m3u8",
	"m4p",
	"m4v",
	"mkv",
	"mng",
	"mov",
	"mp2",
	"mp4",
	"mpe",
	"mpeg",
	"mpg",
	"mpv",
	"mxf",
	"nsv",
	"ogg",
	"ogv",
	"qt",
	"rm",
	"rmvb",
	"roq",
	"svi",
	"vob",
	"webm",
	"wmv",
	"yuv",
}

var subs = []string{
	"srt",
}

func CanServe(file string) bool {
	f := strings.ToLower(file)
	for _, e := range video {
		if strings.HasSuffix(f, e) {
			return true
		}
	}
	for _, e := range subs {
		if strings.HasSuffix(f, e) {
			return true
		}
	}
	return false
}

func IsVideo(file string) bool {
	f := strings.ToLower(file)
	for _, e := range video {
		if strings.HasSuffix(f, e) {
			return true
		}
	}
	return false
}
