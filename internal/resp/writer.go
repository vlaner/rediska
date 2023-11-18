package resp

import "io"

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w Writer) Write(v Value) error {
	b := v.Marshal()
	_, err := w.writer.Write(b)
	if err != nil {
		return err
	}
	return nil
}
