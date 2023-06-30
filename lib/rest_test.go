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
	"testing"
)

func TestCreateAccount(t *testing.T) {
	res := registerAccount("Hans", "a4f388e4-d9a3-4b78-9f9e-3387b6ae87cc", "1cdae0a0-7896-4c9b-b103-6015a7ae1b22", "Timbuktu", "de", true)
	if res != 0 {
		t.Errorf("Expected 0 as return but got %d", res)
	}
}

func TestSetScooping(t *testing.T) {
	Init("/tmp")
	res := setScooping(true)
	if res != 0 {
		t.Errorf("Expected to get Success but got %d", res)
	}
	readAccount()
	if account.Level_1_count != 9 {
		t.Errorf("Expected to get Level1 as 9 but got %d", account.Level_1_count)
	}

	if account.Level_2_count != 99 {
		t.Errorf("Expected to get Level1 as 99 but got %d", account.Level_2_count)
	}

	if account.Level_3_count != 999 {
		t.Errorf("Expected to get Level1 as 999 but got %d", account.Level_3_count)
	}
	os.Remove("/tmp/shift.db")
}

func TestGetMatelist(t *testing.T) {
	list := getMatelist(true)
	if len(list) != 3 {
		t.Errorf("Expected to get 3 mates but got %d", len(list))
	}
	count := 0
	for _, m := range list {
		count += m.FriendsCount
	}
	if count != 12 {
		t.Errorf("Expected to get 12 friends but got %d", count)
	}
}
