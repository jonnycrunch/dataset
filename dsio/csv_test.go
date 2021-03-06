package dsio

import (
	"bytes"
	"testing"

	"github.com/qri-io/dataset"
	"github.com/qri-io/dataset/datatypes"
)

const csvData = `col_a,col_b,col_c,col_d
a,b,c,d
a,b,c,d
a,b,c,d
a,b,c,d
a,b,c,d`

var csvStruct = &dataset.Structure{
	Format: dataset.CSVDataFormat,
	FormatConfig: &dataset.CSVOptions{
		HeaderRow: true,
	},
	Schema: &dataset.Schema{
		Fields: []*dataset.Field{
			{Name: "col_a", Type: datatypes.String},
			{Name: "col_b", Type: datatypes.String},
			{Name: "col_c", Type: datatypes.String},
			{Name: "col_d", Type: datatypes.String},
		},
	},
}

func TestCSVReader(t *testing.T) {
	buf := bytes.NewBuffer([]byte(csvData))
	rdr, err := NewRowReader(csvStruct, buf)
	if err != nil {
		t.Errorf("error allocating RowReader: %s", err.Error())
		return
	}
	count := 0
	for {
		row, err := rdr.ReadRow()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			t.Errorf("unexpected error: %s", err.Error())
			return
		}

		if len(row) != 4 {
			t.Errorf("invalid row length for row %d. expected %d, got %d", count, 4, len(row))
		}

		count++
	}
	if count != 5 {
		t.Errorf("expected: %d rows, got: %d", 5, count)
	}
}

func TestCSVWriter(t *testing.T) {
	rows := [][][]byte{
		// TODO - vary up test input
		{[]byte("a"), []byte("b"), []byte("c"), []byte("d")},
		{[]byte("a"), []byte("b"), []byte("c"), []byte("d")},
		{[]byte("a"), []byte("b"), []byte("c"), []byte("d")},
		{[]byte("a"), []byte("b"), []byte("c"), []byte("d")},
		{[]byte("a"), []byte("b"), []byte("c"), []byte("d")},
	}

	buf := &bytes.Buffer{}
	rw, err := NewRowWriter(csvStruct, buf)
	if err != nil {
		t.Errorf("error allocating RowWriter: %s", err.Error())
		return
	}
	st := rw.Structure()
	if err := dataset.CompareStructures(st, csvStruct); err != nil {
		t.Errorf("structure mismatch: %s", err.Error())
		return
	}

	for i, row := range rows {
		if err := rw.WriteRow(row); err != nil {
			t.Errorf("row %d write error: %s", i, err.Error())
		}
	}

	if err := rw.Close(); err != nil {
		t.Errorf("close reader error: %s", err.Error())
		return
	}
	if bytes.Equal(buf.Bytes(), []byte(csvData)) {
		t.Errorf("output mismatch. %s != %s", buf.String(), csvData)
	}
}
