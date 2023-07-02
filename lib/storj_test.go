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
	"bytes"
	"context"
	"testing"

	"storj.io/uplink"
)

func TestStorj(t *testing.T) {
	Init("/tmp", true)
	bucketName := "shift"
	accessToken := "1GW7L5Hab3vR4twJARK4mMuatA2D319NyYboQXnRQU9JcLDj2BEwwtiZ5whRtwDV4KRPvsfV4HcSjq9DutvF2NLr6yMgij6N6debnCzeLEfPZJds2uLtj4PcQHPXUyzqStdxwTAZrMDJX4RQcvdpqAtbRUVxtbrkg7hRCrjgwTFNCAoATvfeeoXacMkUBMSxpNXLfp3NYWk9KjGgbRC9SkFHDurkrHg8aVs1mMs2vRqW2Y1mcHbpzYthWJxfJB1sQP1shfRyCUZxTY4okb5gnZH3tSSyCPSsSkbLh6KSYnVrb2bqRAr1AgvfQVaB"
	ctx := context.Background()

	access, err := uplink.ParseAccess(accessToken)
	if err != nil {
		t.Errorf("parse access failed %s", err.Error())
		return
	}
	uploadBuffer := []byte("one fish two fish red fish blue fish")
	err = put("foo/bar/baz", uploadBuffer, bucketName, ctx, access)
	if err != nil {
		t.Errorf("put failed: " + err.Error())
	}

	buffer, _, err := get("foo/bar/baz", bucketName, ctx, access)
	if err != nil {
		t.Errorf("get failed: " + err.Error())
	}

	if !bytes.Equal(uploadBuffer, buffer) {
		t.Error("Storj buffers are not identical")
	}

	res, err := exists("foo/bar/baz", bucketName, ctx, access)
	if err != nil {
		t.Errorf("exists failed: " + err.Error())
	}
	if res != true {
		t.Error("exists returned false")
	}

	err = delete("foo/bar/baz", bucketName, ctx, access)
	if err != nil {
		t.Errorf("delete failed: " + err.Error())
	}

	res, err = exists("foo/bar/baz2", bucketName, ctx, access)
	if err == nil {
		t.Errorf("exists failed: " + err.Error())
	}
	if res != false {
		t.Error("exists returned true")
	}
}
