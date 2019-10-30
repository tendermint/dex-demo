package testutil

import "os"

func GetBinDir() string {
	return os.Getenv("UEX_TEST_BIN_DIR")
}

func MustGetBinDir() string {
	dir := GetBinDir()
	if dir == "" {
		panic("must set UEX_TEST_BIN_DIR environment variable")
	}
	return dir
}
