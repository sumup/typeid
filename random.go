package typeid

import (
	"database/sql/driver"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sumup/x/typeid/base32"
)

// Random represents an unique identifier that is entirely random.
// Internally, it's based on UUIDv4.
type Random[P Prefix] struct{ typedID[P] }

var randomIDProc = &processor{
	b32Encode: func(u uuid.UUID) string {
		return base32.EncodeUpper([16]byte(u))
	},
	b32Decode: func(s string) (uuid.UUID, error) {
		decoded, err := base32.DecodeUpper(s)
		if err != nil {
			return uuid.Nil, err
		}
		return uuid.FromBytes(decoded)
	},
	generateUUID: uuid.NewV4,
}

func (Random[P]) processor() *processor {
	return randomIDProc
}

func (Random[P]) Type() string {
	return getPrefix[P]()
}

func (r Random[P]) String() string {
	return toString[P](r.uuid, r.processor())
}

func (r Random[P]) UUID() uuid.UUID {
	return r.uuid
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
// It parses a TypeID string using [FromString]
func (r *Random[P]) UnmarshalText(text []byte) error {
	return unmarshalText(r, text)
}

// MarshalText implements the [encoding.TextMarshaler] interface.
// Internally it use [Random.String]
func (r Random[P]) MarshalText() ([]byte, error) {
	return marshalText(r)
}

func (r Random[P]) Value() (driver.Value, error) {
	return value(r)
}

func (r *Random[P]) Scan(src any) error {
	return scan(r, src)
}

func (r Random[P]) TextValue() (pgtype.Text, error) {
	return textValue(r)
}

func (r *Random[P]) ScanText(v pgtype.Text) error {
	return scanText(r, v)
}

func (r Random[P]) UUIDValue() (pgtype.UUID, error) {
	return uuidValue(r)
}

func (r *Random[P]) ScanUUID(v pgtype.UUID) error {
	return scanUUID(r, v)
}
