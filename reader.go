package afterarch

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

var (
	ErrInvalidMagic = errors.New("afterarch: invalid magic number at end of file")
)

// NewReader will open a ZIP archive at the end of the given file, if any
// exists.  The offset into the file (file pointer) will be modified and will
// not be reset.
func NewReader(f *os.File) (*zip.Reader, error) {
	// Get file size.
	fileSize, err := f.Seek(0, 2)
	if err != nil {
		return nil, err
	}
	if fileSize < trailerSize {
		return nil, ErrInvalidMagic
	}

	// Read the trailer.
	if _, err := f.Seek(-trailerSize, 2); err != nil {
		return nil, err
	}

	var t trailer
	if err := binary.Read(f, binary.LittleEndian, &t); err != nil {
		return nil, err
	}

	// Validate the magic number.
	if !bytes.Equal(t.Magic[:], trailerMagic[:]) {
		return nil, ErrInvalidMagic
	}

	// Seek to beginning of archive.
	if _, err := f.Seek(t.ArchiveSize+trailerSize, 2); err != nil {
		return nil, err
	}

	// Create a reader that can read only the archive
	sr := io.NewSectionReader(f, fileSize-t.ArchiveSize-trailerSize, t.ArchiveSize)

	// Return our ZIP reader.
	return zip.NewReader(sr, t.ArchiveSize)
}
