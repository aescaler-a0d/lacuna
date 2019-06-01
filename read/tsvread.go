/*
 * @File: tsvread.go
 * @Date: 2019-05-31 03:03:54
 * @OA:   antonioe
 * @CA:   antonioe
 * @Time: 2019-06-01 00:13:43
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package read

import (
	"encoding/csv"
	"io"
)

type FieldsReader struct {
	*csv.Reader
}

func NewFieldsReader(r io.Reader) *FieldsReader {
	fr := &FieldsReader{
		Reader: csv.NewReader(r),
	}

	return fr
}

func (r *FieldsReader) Read() (record []string, err error) {
	rec, err := r.Reader.Read()
	if err != nil {
		// parsing variable length fields is hard.
		if err, ok := err.(*csv.ParseError); ok && err.Err == csv.ErrFieldCount {
			return nil, err
		}
		return nil, err
	}
	if len(rec) == 0 {
		return nil, nil
	}
	return rec, nil
}

func (r *FieldsReader) ReadAll() (records [][]string, err error) {
loop:
	for {
		rec, err := r.Read()
		switch err {
		case io.EOF:
			break loop
		case nil:
			records = append(records, rec)
		default:
			return nil, err
		}
	}

	return records, nil
}
