/****************************************************************************
 * Copyright (C) 2023 CrowdWare
 *
 * This file is part of SHIFT.
 *
 *  SHIFT is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  SHIFT is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with SHIFT.  If not, see <http://www.gnu.org/licenses/>.
 *
 ****************************************************************************/
package lib

import (
	"os"
	"runtime"
	"strings"
	"time"
)

var dbFile string
var peerFile string
var messageFile string
var account _account
var lastTransaction _transaction
var initialAmount = int64(initial_amount)
var growLevel0 = int64(10000)
var growLevel1 = int64(1800)
var growLevel2 = int64(360)
var growLevel3 = int64(75)
var useWebService bool

const (
	algorithm = "AES/GCM/NoPadding"
)

type BalanceError struct {
	message string
}

func (e *BalanceError) Error() string {
	return e.message
}

func Encrypt(value string) string {
	return encryptStringGCM(value, false)
}

func Decrypt(enc string) (string, error) {
	return decryptStringGCM(enc, false)
}

func contains(list []Friend, uuid string) int {
	for index, item := range list {
		if item.Uuid == uuid {
			return index
		}
	}
	return -1
}

func decodeUuid(input string) string {
	decodeMap := strings.NewReplacer(
		"0", "A",
		"1", "B",
		"2", "C",
		"3", "D",
		"4", "E",
		"5", "F",
		"6", "G",
		"7", "H",
		"8", "I",
		"9", "J",
		"-", "",
	)
	return decodeMap.Replace(input)
}

func encodeUuid(input string) string {
	encodeMap := strings.NewReplacer(
		"A", "0",
		"B", "1",
		"C", "2",
		"D", "3",
		"E", "4",
		"F", "5",
		"G", "6",
		"H", "7",
		"I", "8",
		"J", "9",
	)
	if len(input) == 36 && input[8] == '-' && input[13] == '-' && input[18] == '-' && input[23] == '-' {
		// Input is already in the correct UUID format
		return input
	}

	input = encodeMap.Replace(input)

	// Insert hyphens at specific positions
	encodedUUID := strings.Builder{}
	for i, c := range input {
		if i == 8 || i == 12 || i == 16 || i == 20 {
			encodedUUID.WriteRune('-')
		}
		encodedUUID.WriteRune(c)
	}
	return encodedUUID.String()
}

func isDevice() bool {
	// Check the value of the runtime.GOARCH
	goarch := runtime.GOARCH
	if goarch == "x86" || goarch == "amd64" {
		return false // Emulator or virtualized environment
	}

	// Check the value of the ANDROID_EMULATOR environment variable
	emulator := os.Getenv("ANDROID_EMULATOR")
	if emulator == "1" {
		return false // Emulator
	}

	// Check if the /sys/qemu_trace file exists
	_, err := os.Stat("/sys/qemu_trace")
	if err == nil {
		return false // Emulator or virtualized environment
	}
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		return false
	}
	return true
}

func formatTime(t time.Time) string {
	now := time.Now().Local()

	if t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day() {
		return t.Format("15:04")
	}

	tISOYear, tISOWeek := t.ISOWeek()
	nowISOYear, nowISOWeek := now.ISOWeek()

	if tISOYear == nowISOYear && tISOWeek == nowISOWeek {
		switch t.Weekday() {
		case time.Monday:
			return "Mo"
		case time.Tuesday:
			return "Di"
		case time.Wednesday:
			return "Mi"
		case time.Thursday:
			return "Do"
		case time.Friday:
			return "Fr"
		case time.Saturday:
			return "Sa"
		case time.Sunday:
			return "So"
		}
	}
	return t.Format("02.01.2006")
}
