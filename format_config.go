package dataset

import (
	"fmt"
)

type FormatConfig interface {
	Format() DataFormat
	Map() map[string]interface{}
}

func ParseFormatConfigMap(f DataFormat, opts map[string]interface{}) (FormatConfig, error) {
	switch f {
	case CsvDataFormat:
		return NewCsvOptions(opts)
	case JsonDataFormat:
		return NewJsonOptions(opts)
	}

	return nil, nil
}

func NewCsvOptions(opts map[string]interface{}) (FormatConfig, error) {
	o := &CsvOptions{}
	if opts == nil {
		return o, nil
	}
	if opts["header_row"] != nil {
		if headerRow, ok := opts["header_row"].(bool); ok {
			o.HeaderRow = headerRow
		} else {
			return nil, fmt.Errorf("invalid header_row value: %s", opts["header_row"])
		}
	}

	return o, nil
}

type CsvOptions struct {
	// Weather this csv file has a header row or not
	HeaderRow bool `json:"header_row"`
}

func (*CsvOptions) Format() DataFormat {
	return CsvDataFormat
}

func (o *CsvOptions) Map() map[string]interface{} {
	if o == nil {
		return nil
	}
	return map[string]interface{}{
		"header_row": o.HeaderRow,
	}
}

func NewJsonOptions(opts map[string]interface{}) (FormatConfig, error) {
	o := &JsonOptions{}
	if opts == nil {
		return o, nil
	}
	if opts["object_entries"] != nil {
		if objEntries, ok := opts["object_entries"].(bool); ok {
			o.ObjectEntries = objEntries
		} else {
			return nil, fmt.Errorf("invalid object_entries value: %s", opts["object_entries"])
		}
	}
	return o, nil
}

type JsonOptions struct {
	ObjectEntries bool `json:"object_entries"`
}

func (*JsonOptions) Format() DataFormat {
	return JsonDataFormat
}

func (o *JsonOptions) Map() map[string]interface{} {
	if o == nil {
		return nil
	}
	return map[string]interface{}{
		"object_entries": o.ObjectEntries,
	}
}
