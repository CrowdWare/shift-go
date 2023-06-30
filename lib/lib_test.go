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
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestToHex(t *testing.T) {
	guuid := uuid.NewString()
	hex := decodeUuid(guuid)
	huuid := encodeUuid(hex)

	if guuid != huuid {
		t.Errorf("Uuids are different %s and %s", guuid, huuid)
	}
}

func TestFormatTime(t *testing.T) {
	now := time.Now().Local()
	date := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.Local)

	comp := fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute())
	if formatTime(date) != comp {
		t.Errorf("Expected result to be %s but got, %s", comp, formatTime(date))
	}

	date = date.AddDate(0, 0, -2)
	if len(formatTime(date)) != 2 {
		t.Errorf("Expected result to be 2 chars long but got, %s", formatTime(date))
	}

	date = date.AddDate(0, 0, -7)
	if len(formatTime(date)) < 3 {
		t.Errorf("Expected result to be longer than 2 chars long but got, %s", formatTime(date))
	}
}
