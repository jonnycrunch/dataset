package dataset

import (
	"fmt"
)

func CompareSchemas(a, b *Schema) error {
	// if a.PrimaryKey != b.PrimaryKey {
	// 	return fmt.Errorf("primary key mismatch: %s != %s", a.PrimaryKey, b.PrimaryKey)
	// }
	if a == nil && b == nil {
		return nil
	}
	if a != nil && b == nil || a == nil && b != nil {
		return fmt.Errorf("nil mismatch: %s != %s", a, b)
	}
	if a.Fields == nil && b.Fields != nil || a.Fields != nil && b.Fields == nil {
		return fmt.Errorf("fields slice mismatch: %s != %s", a.Fields, b.Fields)
	}
	if a.Fields == nil && b.Fields == nil {
		return nil
	}

	if len(a.Fields) != len(b.Fields) {
		return fmt.Errorf("field length mismatch: %d != %d", len(a.Fields), len(b.Fields))
	}

	for i, af := range a.Fields {
		bf := b.Fields[i]
		if err := CompareFields(af, bf); err != nil {
			return fmt.Errorf("field %d mismatch: %s", i, err.Error())
		}
	}

	return nil
}

func CompareFields(a, b *Field) error {
	if a == nil && b == nil {
		return nil
	} else if a == nil && b != nil || a != nil && b == nil {
		return fmt.Errorf("nil mismatch: %s != %s", a, b)
	}

	if a.Name != b.Name {
		return fmt.Errorf("name mismatch: %s != %s", a.Name, b.Name)
	}
	if a.Type != b.Type {
		return fmt.Errorf("field type mismatch: %s != %s", a.Type.String(), b.Type.String())
	}
	if a.Title != b.Title {
		return fmt.Errorf("title mismatch: %s != %s", a.Title, b.Title)
	}
	if a.Description != b.Description {
		return fmt.Errorf("description mismatch: %s != %s", a.Description, b.Description)
	}

	// TODO - finish comparison of field constraints, primary keys, format, etc.

	return nil
}
