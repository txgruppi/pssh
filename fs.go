package main

import (
	"io"
	"os"
)

func isFile(filepath string) bool {
	info, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

func makeTempExecFile(content []byte) (string, error) {
	tmpfile, err := os.CreateTemp("", "pssh-temp-exec-*")
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()
	if _, err := tmpfile.Write([]byte("#!/usr/bin/env bash\nset -euo pipefail\n")); err != nil {
		return "", err
	}
	if _, err := tmpfile.Write(content); err != nil {
		return "", err
	}
	if err := tmpfile.Chmod(0755); err != nil {
		return "", err
	}
	return tmpfile.Name(), nil
}

func copyTempFile(src string) (string, error) {
	dst, err := os.CreateTemp("", "pssh-temp-copy-*")
	if err != nil {
		return "", err
	}
	defer dst.Close()

	srcFile, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	if _, err := io.Copy(dst, srcFile); err != nil {
		return "", err
	}

	info, err := srcFile.Stat()
	if err != nil {
		return "", err
	}

	mode := info.Mode()
	mode |= 0700
	if err := dst.Chmod(mode); err != nil {
		return "", err
	}

	return dst.Name(), nil
}
