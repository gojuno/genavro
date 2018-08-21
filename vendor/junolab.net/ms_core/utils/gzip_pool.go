package utils

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
)

type (
	GzipPool struct {
		gzipWriterPool *sync.Pool
		gzipReaderPool *sync.Pool
	}
)

func NewGzipPool() *GzipPool {
	return &GzipPool{
		gzipWriterPool: &sync.Pool{
			New: func() interface{} {
				return gzip.NewWriter(nil)
			},
		},
		gzipReaderPool: &sync.Pool{},
	}
}

func (p *GzipPool) AcquireReader(r io.Reader) (*gzip.Reader, error) {
	gz, ok := p.gzipReaderPool.Get().(*gzip.Reader)
	var err error
	if !ok || gz == nil {
		gz, err = gzip.NewReader(r)
	} else {
		err = gz.Reset(r)
	}
	return gz, err
}

func (p *GzipPool) ReleaseReader(r io.ReadCloser) error {
	err := r.Close()
	p.gzipReaderPool.Put(r)
	return err
}

func (p *GzipPool) AcquireWriter(w io.Writer) *gzip.Writer {
	gz := p.gzipWriterPool.Get().(*gzip.Writer)
	gz.Reset(w)
	return gz
}

func (p *GzipPool) ReleaseWriter(w io.WriteCloser) error {
	err := w.Close()
	p.gzipWriterPool.Put(w)
	return err
}

func (p *GzipPool) Compress(msg []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzipWriter := p.AcquireWriter(&buf)
	defer p.ReleaseWriter(gzipWriter)

	if _, err := gzipWriter.Write(msg); err != nil {
		return nil, fmt.Errorf("failed to gzip msg with err %s", err)
	}
	// need for flushing EOF
	gzipWriter.Close()
	return buf.Bytes(), nil
}

func (p *GzipPool) Decompress(b []byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write(b)
	gzipReader, err := p.AcquireReader(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire gzip reader: %q", err)
	}
	bytes, err := ioutil.ReadAll(gzipReader)
	p.ReleaseReader(gzipReader)
	if err != nil {
		return nil, fmt.Errorf("failed to unzip blob: %q", err)
	}
	return bytes, nil
}
