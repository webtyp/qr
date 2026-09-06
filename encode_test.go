package qr

import (
	"testing"
)

func TestErrorsAndVersionSelection(t *testing.T) {
	// Encode("", LevelM) -> ErrEmptyData
	_, err := Encode("", LevelM)
	if err != ErrEmptyData {
		t.Fatalf("Encode(\"\") err = %v; want ErrEmptyData", err)
	}

	// Encode(data, Level(99)) -> ErrUnknownLevel
	_, err = Encode("test", Level(99))
	if err != ErrUnknownLevel {
		t.Fatalf("Encode(level 99) err = %v; want ErrUnknownLevel", err)
	}

	// Guard: Level(0) == LevelM
	if Level(0) != LevelM {
		t.Fatalf("Level(0) != LevelM")
	}

	// 1-byte payload -> Version 1
	v1, err := selectVersion(1, LevelM)
	if err != nil || v1 != 1 {
		t.Fatalf("selectVersion(1, LevelM) = (%d, %v); want (1, nil)", v1, err)
	}

	v2, err := selectVersion(15, LevelM)
	if err != nil || v2 != 2 {
		t.Fatalf("selectVersion(15, LevelM) = (%d, %v); want (2, nil)", v2, err)
	}

	// Exceed version 40 at level H -> ErrTooLong
	_, err = selectVersion(2000, LevelH)
	if err != ErrTooLong {
		t.Fatalf("selectVersion(2000, LevelH) err = %v; want ErrTooLong", err)
	}
}
