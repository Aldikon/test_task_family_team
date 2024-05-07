package unzip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func UnzipFile(path, destToUnzip string) error {
	r, err := zip.OpenReader(path)
	if err != nil {
		return fmt.Errorf("open file '%s', err: %w", path, err)
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("open file in zip '%s', err: %w", f.Name, err)
		}
		defer rc.Close()

		filePath := filepath.Join(destToUnzip, f.Name)

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(filePath, f.Mode())
			if err != nil {
				return fmt.Errorf("can't create dir '%s', err: %w", filePath, err)
			}
		} else {
			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return fmt.Errorf("open file to be copied '%v', err: %w", filePath, err)
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, rc)
			if err != nil {
				return fmt.Errorf("copy file, err: %w", err)
			}
		}
	}

	return nil
}
