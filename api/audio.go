package api

import (
	"io"
	"sync"
)

type MP3Chunks struct {
	rcs []io.ReadCloser
	mu  sync.Mutex
}

var _ io.ReadCloser = (*MP3Chunks)(nil)

func NewMP3Chunks(content []byte) *MP3Chunks {
	return &MP3Chunks{
		rcs: []io.ReadCloser{},
	}
}

func (m *MP3Chunks) Add(rc io.ReadCloser) {
	m.mu.Lock()
	m.rcs = append(m.rcs, rc)
	m.mu.Unlock()
}

func (m *MP3Chunks) Read(p []byte) (n int, err error) {
	if len(m.rcs) == 0 {
		return 0, io.EOF
	}

	for len(m.rcs) > 0 && n < len(p) {
		m.mu.Lock()
		n, err = m.rcs[0].Read(p)
		m.mu.Unlock()

		if err != nil {
			if err == io.EOF {
				m.mu.Lock()
				m.rcs[0].Close()
				m.rcs = m.rcs[1:]
				m.mu.Unlock()

				continue
			}
			return n, err
		}

	}

	return n, nil
}

func (m *MP3Chunks) Close() error {
	for _, rc := range m.rcs {
		if err := rc.Close(); err != nil {
			return err
		}
	}
	return nil
}
