package dataset

import (
	"encoding/json"
	"fmt"
	"github.com/ipfs/go-datastore"
	"github.com/qri-io/dataset/compression"
)

// Structure designates a deterministic definition for working with a discrete dataset.
// Structure is a concrete handle that provides precise details about how to interpret a given
// piece of data (the reference to the data itself is provided elsewhere, specifically in the dataset struct )
// These techniques provide mechanisms for joining & traversing multiple structures.
// This example is shown in a human-readable form, for storage on the network the actual
// output would be in a condensed, non-indented form, with keys sorted by lexographic order.
type Structure struct {
	// Format specifies the format of the raw data MIME type
	Format DataFormat `json:"format"`
	// FormatConfig removes as much ambiguity as possible about how
	// to interpret the speficied format.
	FormatConfig FormatConfig `json:"formatConfig,omitempty"`
	// Encoding specifics character encoding
	// should assume utf-8 if not specified
	Encoding string `json:"encoding,omitempty"`
	// Compression specifies any compression on the source data,
	// if empty assume no compression
	Compression compression.Type `json:"compression,omitempty"`
	// Schema contains the schema definition for the underlying data
	Schema *Schema `json:"schema"`
}

// Abstract returns this structure instance in it's "Abstract" form
// stripping all nonessential values &
// renaming all schema field names to standard variable names
func (s *Structure) Abstract() *Structure {
	a := &Structure{
		Format:       s.Format,
		FormatConfig: s.FormatConfig,
		Encoding:     s.Encoding,
	}
	if s.Schema != nil {
		a.Schema = &Schema{
			PrimaryKey: s.Schema.PrimaryKey,
			Fields:     make([]*Field, len(s.Schema.Fields)),
		}
		for i, f := range s.Schema.Fields {
			a.Schema.Fields[i] = &Field{
				Name:         fmt.Sprintf("col_%d", i),
				Type:         f.Type,
				MissingValue: f.MissingValue,
				Format:       f.Format,
				Constraints:  f.Constraints,
			}
		}
	}
	return a
}

// Hash gives the hash of this structure
func (r *Structure) Hash() (string, error) {
	return JSONHash(r)
}

// separate type for marshalling into & out of
// most importantly, struct names must be sorted lexographically
type _structure struct {
	Compression  compression.Type       `json:"compression,omitempty"`
	Encoding     string                 `json:"encoding,omitempty"`
	Format       DataFormat             `json:"format"`
	FormatConfig map[string]interface{} `json:"formatConfig,omitempty"`
	Schema       *Schema                `json:"schema,omitempty"`
}

// MarshalJSON satisfies the json.Marshaler interface
func (r Structure) MarshalJSON() (data []byte, err error) {
	var opt map[string]interface{}
	if r.FormatConfig != nil {
		opt = r.FormatConfig.Map()
	}

	return json.Marshal(&_structure{
		Compression:  r.Compression,
		Encoding:     r.Encoding,
		Format:       r.Format,
		FormatConfig: opt,
		Schema:       r.Schema,
	})
}

// UnmarshalJSON satisfies the json.Unmarshaler interface
func (r *Structure) UnmarshalJSON(data []byte) error {
	_r := &_structure{}
	if err := json.Unmarshal(data, _r); err != nil {
		return err
	}

	fmtCfg, err := ParseFormatConfigMap(_r.Format, _r.FormatConfig)
	if err != nil {
		return err
	}

	*r = Structure{
		Compression:  _r.Compression,
		Encoding:     _r.Encoding,
		Format:       _r.Format,
		FormatConfig: fmtCfg,
		Schema:       _r.Schema,
	}

	// TODO - question of weather we should not accept
	// invalid structure defs at parse time. For now we'll take 'em.
	// if err := d.Valid(); err != nil {
	//   return err
	// }

	return nil
}

// Valid validates weather or not this structure
func (ds *Structure) Valid() error {
	// if count := truthCount(ds.Url != "", ds.File != "", len(ds.Data) > 0); count > 1 {
	// 	return errors.New("only one of url, file, or data can be set")
	// } else if count == 1 {
	// 	if ds.Format == UnknownDataFormat {
	// 		// if format is unspecified, we need to be able to derive the format from
	// 		// the extension of either the url or filepath
	// 		if ds.DataFormat() == "" {
	// 			return errors.New("format is required for data source")
	// 		}
	// 	}
	// }
	return nil
}

// LoadStructure loads a structure from a given path in a store
func LoadStructure(store datastore.Datastore, path datastore.Key) (*Structure, error) {
	v, err := store.Get(path)
	if err != nil {
		return nil, err
	}

	return UnmarshalStructure(v)
}

// UnmarshalStructure tries to extract a structure type from an empty
// interface. Pairs nicely with datastore.Get() from github.com/ipfs/go-datastore
func UnmarshalStructure(v interface{}) (*Structure, error) {
	switch r := v.(type) {
	case *Structure:
		return r, nil
	case Structure:
		return &r, nil
	case []byte:
		structure := &Structure{}
		err := json.Unmarshal(r, structure)
		return structure, err
	default:
		return nil, fmt.Errorf("couldn't parse structure")
	}
}
