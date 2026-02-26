//go:build !windows && !darwin

package main

import (
	"os"
	"time"
)

func fileCreationTime(info os.FileInfo) (time.Time, bool) {
	return time.Time{}, false
}
