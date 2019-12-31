// Copyright 2019 ShouDong Zheng. All rights reserved.
// Use of this source code is governed by a Apache-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"./slowlogsparser"
)

var (
	_slowlogPath       = flag.String("slowlog", "", "Slow log file or directory path")
	_redisCommand      = flag.String("command", "", "Redis command")
	_durationThreshold = flag.Float64("duration", 0.0, "Slow query duration threshold")
	slowlogResult      []slowlogsparser.Slowlog
)

func main() {
	flag.Parse()
	if len(*_slowlogPath) == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	if isDir(*_slowlogPath) {
		slowlogResult = slowlogsparser.ParserLogs(getLogFiles(*_slowlogPath), *_durationThreshold, *_redisCommand)
	} else {
		slowlogResult = slowlogsparser.ParserLogs([]string{*_slowlogPath}, *_durationThreshold, *_redisCommand)
	}

	for _, log := range slowlogResult {
		fmt.Println(slowlogsparser.ToString(log))
	}
}

func getLogFiles(filesDir string) []string {
	filePaths := []string{}
	files, _ := ioutil.ReadDir(filesDir)
	for _, file := range files {
		fileName := file.Name()
		if strings.HasPrefix(fileName, ".") {
			continue
		}

		filePath := path.Join(filesDir, fileName)
		if file.IsDir() {
			for _, p := range getLogFiles(filePath) {
				filePaths = append(filePaths, p)
			}
		} else {
			if strings.HasSuffix(fileName, ".log") {
				filePaths = append(filePaths, filePath)
			}
		}
	}

	return filePaths
}

func isDir(p string) bool {
	f, err := os.Stat(p)
	if err != nil {
		panic(err)
	}

	return f.IsDir()
}
