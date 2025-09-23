package httpadapter

import (
    "fmt"
    "io"
    "mime/multipart"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/google/uuid"
)

// saveUploadedFile saves the given multipart file under uploadDir/YYYY/MM/uuid.ext
// and returns the filesystem path and public URL (/static/...).
func saveUploadedFile(uploadDir string, fh *multipart.FileHeader) (string, string, error) {
    // Ensure date-based subdirectory
    now := time.Now().UTC()
    sub := filepath.Join(fmt.Sprintf("%04d", now.Year()), fmt.Sprintf("%02d", int(now.Month())))
    dir := filepath.Join(uploadDir, sub)
    if err := os.MkdirAll(dir, 0o755); err != nil {
        return "", "", err
    }

    // Build filename
    ext := filepath.Ext(fh.Filename)
    if len(ext) > 10 { // guard weird extremely long extensions
        ext = ext[:10]
    }
    // sanitize ext a bit
    ext = strings.ToLower(ext)
    name := uuid.New().String() + ext
    path := filepath.Join(dir, name)

    // Copy file contents
    src, err := fh.Open()
    if err != nil {
        return "", "", err
    }
    defer src.Close()

    dst, err := os.Create(path)
    if err != nil {
        return "", "", err
    }
    defer dst.Close()

    if _, err := io.Copy(dst, src); err != nil {
        return "", "", err
    }

    // Public URL under /static
    rel := filepath.ToSlash(filepath.Join(sub, name))
    publicURL := "/static/" + rel
    return path, publicURL, nil
}

