package afterarch

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

type Writer struct {
	w               *zip.Writer
	f               *os.File
	closeUnderlying bool
	buf             bytes.Buffer
}

// NewWriterAfter returns a new Writer that appends data to the given file.
// The underlying File will not be closed when this Writer is closed.
func NewWriterAfter(f *os.File) *Writer {
	ret := &Writer{
		w:               nil,
		f:               f,
		closeUnderlying: false,
	}
	ret.w = zip.NewWriter(&ret.buf)
	return ret
}

// NewWriterAfterThis copies the current binary (as found in os.Args[0]) to the
// given path, and then creates a Writer after the new file.
func NewWriterAfterThis(output string) (*Writer, error) {
	// Open ourselves
	ourselves, err := os.Open(os.Args[0])
	if err != nil {
		return nil, err
	}
	defer ourselves.Close()

	// Create a new executable file.
	outFile, err := os.OpenFile(
		output,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		0755,
	)
	if err != nil {
		return nil, err
	}

	// Copy ourselves to output.
	_, err = io.Copy(outFile, ourselves)
	if err != nil {
		outFile.Close()
		return nil, err
	}

	// Create a new writer
	w := NewWriterAfter(outFile)

	// Tell the writer to close the underlying file.
	w.closeUnderlying = true
	return w, nil
}

// Close finishes writing the archive.
func (w *Writer) Close() error {
	// Close the ZIP file to force it to flush to our buffer.
	if err := w.w.Close(); err != nil {
		return err
	}

	// Get the size of the buffer.
	bufSize := w.buf.Len()

	// Write the data at the end of the file.
	if _, err := w.f.Seek(0, 2); err != nil {
		return err
	}
	if _, err := io.Copy(w.f, &w.buf); err != nil {
		return err
	}

	// Write a trailer indicating the size of the buffer and a magic value.
	t := trailer{
		Magic:       trailerMagic,
		ArchiveSize: int64(bufSize),
	}
	if err := binary.Write(w.f, binary.LittleEndian, &t); err != nil {
		return err
	}

	// Possibly close file
	if w.closeUnderlying {
		if err := w.f.Close(); err != nil {
			return err
		}
	}

	// Done!
	return nil
}

// Create simply proxies to the underlying zip.Writer implementation.
func (w *Writer) Create(name string) (io.Writer, error) {
	return w.w.Create(name)
}

// CreateHeader simply proxies to the underlying zip.Writer implementation.
func (w *Writer) CreateHeader(fh *zip.FileHeader) (io.Writer, error) {
	return w.w.CreateHeader(fh)
}

// Flush simply proxies to the underlying zip.Writer implementation.
func (w *Writer) Flush() error {
	return w.w.Flush()
}
