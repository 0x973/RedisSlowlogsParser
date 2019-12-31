// Copyright 2019 ShouDong Zheng. All rights reserved.
// Use of this source code is governed by a Apache-style
// license that can be found in the LICENSE file.

package slowlogsparser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"../concurrentmap"
)

const (
	_carriageReturn = "\r"
	_linefeed       = "\n"
	_lienfeedByte   = '\n'
	_crlf           = "\r\n"
)

// Slowlog - records slowlog message
type Slowlog struct {
	slowID     int64
	timestamp  int64
	duration   float64
	command    string
	key        string
	parameters []string
}

// ParserLogs - Start to parser file(s)
func ParserLogs(filePaths []string, durationThreshold float64, redisCommand string) []Slowlog {
	// remove duplication logs
	slowlogMap := concurrentmap.New()
	for _, filePath := range filePaths {
		parserLog(filePath, func(slowlogObject Slowlog) {
			slowlogMap.Put(slowlogObject.slowID, slowlogObject)
		})
	}

	var result []Slowlog
	for _, slowlogObject := range slowlogMap.Values() {
		log := slowlogObject.(Slowlog)
		if log.duration > durationThreshold*1000 {
			if len(redisCommand) > 0 && redisCommand != log.command {
				continue
			}

			result = append(result, log)
		}
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].timestamp < result[j].timestamp
	})

	return result
}

// ToString - Slowlog convert to string
func ToString(log Slowlog) string {
	date := time.Unix(log.timestamp, 0).Format("2006-01-02 15:04:05")
	var stringBuilder strings.Builder
	stringBuilder.WriteString(fmt.Sprintf("date: %s, ", date))
	stringBuilder.WriteString(fmt.Sprintf("slowID: %d, ", log.slowID))
	stringBuilder.WriteString(fmt.Sprintf("duration: %.2fms, ", log.duration/1000.0))
	stringBuilder.WriteString(fmt.Sprintf("command: %s, ", log.command))

	if log.key != "" {
		stringBuilder.WriteString(fmt.Sprintf("key: %s, ", log.key))
	}

	// for index, parm := range log.parameters {
	// 	stringBuilder.WriteString(fmt.Sprintf("parameter%d: %s, ", index, parm))
	// }

	str := stringBuilder.String()
	return str[0 : len(str)-2]
}

func parserLog(p string, callback func(Slowlog)) {
	file, err := os.Open(p)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_reader := bufio.NewReader(file)
	count := 0
	slowlogObject := Slowlog{}

	for {
		lineString, err := _reader.ReadString(_lienfeedByte)
		if err != nil && err == io.EOF {
			break
		}

		lineString = strings.TrimSuffix(lineString, _carriageReturn)
		lineString = strings.TrimSuffix(lineString, _linefeed)
		count++

		// normally, the fourth is command, the fifth start is the argument
		if count > 4 {
			prevSlowID := slowlogObject.slowID - 1
			afterSlowID := slowlogObject.slowID + 1
			tempIDValue := tryParseToInt(lineString)
			if tempIDValue != -1 {
				if tempIDValue == afterSlowID || tempIDValue == prevSlowID {
					callback(slowlogObject)
					// reset
					slowlogObject = Slowlog{}
					count = 1
				}
			}
		}

		switch count {
		case 1:
			slowlogObject.slowID = parseToInt(lineString)
			break
		case 2:
			slowlogObject.timestamp = parseToInt(lineString)
			break
		case 3:
			slowlogObject.duration = parseToFloat(lineString)
			break
		case 4:
			slowlogObject.command = lineString
			break
		default:
			if len(slowlogObject.key) == 0 {
				slowlogObject.key = lineString
			} else {
				slowlogObject.parameters = append(slowlogObject.parameters, lineString)
			}
			break
		}
	}
}

func parseToInt(str string) int64 {
	value, err := strconv.ParseInt(str, 0, 64)
	if err != nil {
		panic(err)
	}

	return value
}

func parseToFloat(str string) float64 {
	value, err := strconv.ParseFloat(str, 0)
	if err != nil {
		panic(err)
	}

	return value
}

func tryParseToInt(str string) int64 {
	value, err := strconv.ParseInt(str, 0, 64)
	if err != nil {
		return -1
	}

	return value
}
