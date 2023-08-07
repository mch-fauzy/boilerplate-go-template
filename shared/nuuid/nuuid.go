package nuuid

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gofrs/uuid"
)

// NUUID is a nullable UUID, Use it if you want to make UUID nullable in DB/Json model/etc
type NUUID struct {
	UUID  uuid.UUID
	Valid bool
}

// From takes a UUID and returns a NUUID
func From(id uuid.UUID) NUUID {
	return New(id, true)
}

// FromString takes a string and returns a NUUID
func FromString(str string) NUUID {
	id, err := uuid.FromString(str)
	return New(id, err == nil)
}

// New creates a new NUUID
func New(id uuid.UUID, valid bool) NUUID {
	return NUUID{
		UUID:  id,
		Valid: valid,
	}
}

// Scan implements the Scanner interface.
func (nid *NUUID) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case []byte:
		return nid.UnmarshalText(x)
	case string:
		return nid.UnmarshalText([]byte(x))
	case nil:
		nid.Valid = false
		return nil
	default:
		err = fmt.Errorf("null: cannot scan type %T into nuuid.NUUID: %v", value, value)
		nid.Valid = false
		return err
	}
}

// Value implements the driver Valuer interface.
func (nid NUUID) Value() (driver.Value, error) {
	if !nid.Valid {
		return nil, nil
	}
	return nid.UUID.String(), nil
}

// UnmarshalJSON implements the UnmarshalJSON method
func (nid *NUUID) UnmarshalJSON(data []byte) error {
	var err error
	var id uuid.UUID
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		id, err = uuid.FromString(x)
		nid.UUID = id
	case map[string]interface{}:
		str, strOK := x["UUID"].(string)
		valid, validOK := x["Valid"].(bool)
		if !strOK || !validOK {
			return fmt.Errorf(`json: unmarshalling object into Go value of type nuuid.NUUID requires key "UUID" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively`, x["UUID"], x["Valid"])
		}
		id, err = uuid.FromString(str)
		nid.UUID = id
		nid.Valid = valid
		return err
	case nil:
		nid.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type nuuid.UUID", reflect.TypeOf(v).Name())
	}
	nid.Valid = err == nil
	return err
}

// MarshalJSON implements the MarshalJSON method
func (nid NUUID) MarshalJSON() ([]byte, error) {
	if !nid.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nid.UUID.String())
}

// MarshalText implements the MarshalText method
func (nid NUUID) MarshalText() ([]byte, error) {
	if !nid.Valid {
		return []byte{}, nil
	}
	return nid.UUID.MarshalText()
}

// UnmarshalText implements the UnmarshalText method
func (nid *NUUID) UnmarshalText(text []byte) error {
	str := string(text)
	id, err := uuid.FromString(str)
	nid.UUID = id
	if err != nil {
		return err
	}
	nid.Valid = true
	return nil
}

// SetValid sets the value of a NUUID
func (nid *NUUID) SetValid(v uuid.UUID) {
	nid.UUID = v
	nid.Valid = true
}

// Ptr returns the pointer to a NUUID
func (nid NUUID) Ptr() *uuid.UUID {
	if !nid.Valid {
		return nil
	}
	return &nid.UUID
}

// IsZero checks whether the NUUID is null or not
func (nid NUUID) IsZero() bool {
	return !nid.Valid
}
