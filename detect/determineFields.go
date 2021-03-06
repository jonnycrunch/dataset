package detect

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/qri-io/dataset"
	"github.com/qri-io/dataset/datatypes"
)

var (
	startsWithNumberRegex = regexp.MustCompile(`^[0-9]`)
)

// Fields determines the fields of a given reader, for a given structure
func Fields(r *dataset.Structure, data io.Reader) (fields []*dataset.Field, err error) {
	if r.Format == dataset.UnknownDataFormat {
		return nil, errors.New("dataset format must be specified to determine fields")
	}

	switch r.Format {
	case dataset.CSVDataFormat:
		return CSVFields(r, data)
	case dataset.JSONDataFormat:
		return JSONFields(r, data)
	case dataset.XMLDataFormat:
		return XMLFields(r, data)
	}

	return nil, fmt.Errorf("'%s' is not supported for field detection", r.Format.String())
}

// CSVFields determines the field names and types of an io.Reader of CSV-formatted data
func CSVFields(resource *dataset.Structure, data io.Reader) (fields []*dataset.Field, err error) {
	r := csv.NewReader(data)
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = true
	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	fields = make([]*dataset.Field, len(header))
	types := make([]map[datatypes.Type]int, len(header))

	for i := range fields {
		fields[i] = &dataset.Field{
			Name: fmt.Sprintf("field_%d", i+1),
			Type: datatypes.Any,
		}
		types[i] = map[datatypes.Type]int{}
	}

	if possibleCsvHeaderRow(header) {
		for i, f := range fields {
			f.Name = Camelize(header[i])
			f.Type = datatypes.Any
		}
		resource.FormatConfig = &dataset.CSVOptions{
			HeaderRow: true,
		}
		// ds.HeaderRow = true
	} else {
		for i, cell := range header {
			types[i][datatypes.ParseDatatype([]byte(cell))]++
		}
	}

	count := 0
	for {
		rec, err := r.Read()
		// fmt.Println(rec)
		if count > 2000 {
			break
		}
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fields, err
		}

		for i, cell := range rec {
			types[i][datatypes.ParseDatatype([]byte(cell))]++
		}

		count++
	}

	for i, tally := range types {
		for typ, count := range tally {
			if count > tally[fields[i].Type] {
				fields[i].Type = typ
			}
		}
	}

	return fields, nil
}

// JSONFields determines the field names and types of a given io.Reader of JSON-formatted data
func JSONFields(ds *dataset.Structure, data io.Reader) (fields []*dataset.Field, err error) {
	// TODO
	return nil, errors.New("json field detection not yet implemented")
}

// XMLFields determines the field names and types of a given io.Reader of XML-formatted data
func XMLFields(ds *dataset.Structure, data io.Reader) (fields []*dataset.Field, err error) {
	// TODO
	return nil, errors.New("xml field detection not yet implemented")
}

// PossibleHeaderRow makes an educated guess about weather or not this csv file has a header row.
// If this returns true, a determination about weather this data contains a header row should be
// made by comparing with the destination schema.
// This is because it's not totally possible to determine if csv data has a header row based on the
// data alone.
// For example, if all columns are a string data type, and all fields in the first row
// are provided, it isn't possible to distinguish between a header row and an entry
func possibleCsvHeaderRow(header []string) bool {
	for _, rawCol := range header {
		col := strings.TrimSpace(rawCol)
		if _, err := datatypes.ParseInteger([]byte(col)); err == nil {
			// if the row contains valid numeric data, we out.
			return false
		} else if _, err := datatypes.ParseFloat([]byte(col)); err == nil {
			return false
		} else if col == "" {
			// empty columns can't be headers
			return false
		} else if col == "true" || col == "false" {
			// true & false are keywords, and cannot be headers
			return false
		}
	}
	return true
}
