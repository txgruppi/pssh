package main

import (
	"fmt"
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
		return "", fmt.Errorf("error creating temp exec file: %w", err)
	}
	defer tmpfile.Close()
	if _, err := tmpfile.Write([]byte("#!/usr/bin/env bash\nset -euo pipefail\n")); err != nil {
		return "", fmt.Errorf("error writing to temp exec file: %w", err)
	}
	if _, err := tmpfile.Write(content); err != nil {
		return "", fmt.Errorf("error writing to temp exec file: %w", err)
	}
	if err := tmpfile.Chmod(0o755); err != nil {
		return "", fmt.Errorf("error chmodding temp exec file: %w", err)
	}
	return tmpfile.Name(), nil
}

func copyTempFile(src string) (string, error) {
	dst, err := os.CreateTemp("", "pssh-temp-copy-*")
	if err != nil {
		return "", fmt.Errorf("error opening source file: %w", err)
	}
	defer dst.Close()

	srcFile, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("error copying temp file: %w", err)
	}
	defer srcFile.Close()

	if _, err := io.Copy(dst, srcFile); err != nil {
		return "", fmt.Errorf("error copying temp file: %w", err)
	}

	info, err := srcFile.Stat()
	if err != nil {
		return "", fmt.Errorf("error statting source file: %w", err)
	}

	mode := info.Mode()
	mode |= 0o700
	if err := dst.Chmod(mode); err != nil {
		return "", fmt.Errorf("error chmodding temp file: %w", err)
	}

	return dst.Name(), nil
}
