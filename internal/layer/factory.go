package layer

import (
	"archive/tar"
	"github.com/buildpacks/pack/internal/archive"
	"io"

	"github.com/buildpacks/imgutil"
)

type TarWriterFactory interface {
	NewTarWriter(io.Writer) TarWriter
}

type tarWriterFactory struct {
	os string
}

type TarWriter interface {
	WriteHeader(hdr *tar.Header) error
	Write(b []byte) (int, error)
	Close() error
}

var DefaultTarWriterFactory = tarWriterFactory{os: "linux"}

func NewTarWriterFactory(image imgutil.Image) (archive.TarWriterFactory, error) {
	os, err := image.OS()
	if err != nil {
		return nil, err
	}

	return tarWriterFactory{os: os}, nil
}

func (f tarWriterFactory) NewTarWriter(fileWriter io.Writer) archive.TarWriter {
	if f.os == "windows" {
		return NewWindowsWriter(fileWriter)
	}

	// Linux images use tar.Writer
	return tar.NewWriter(fileWriter)
}

// TODO: Move to method on `imgutil.Image`
func NewWriterForImage(image imgutil.Image, fileWriter io.Writer) (TarWriter, error) {
	os, err := image.OS()
	if err != nil {
		return nil, err
	}
	if os == "windows" {
		return NewWindowsWriter(fileWriter), nil
	}

	// Linux images use tar.Writer
	return tar.NewWriter(fileWriter), nil
}

/*
imgutil      lifecycle
    ^          ^
     \        /
        pack


*/
