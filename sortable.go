package typeid

import (
	"database/sql/driver"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sumup/typeid/base32"
)

// Sortable represents an unique identifier that is k-sortable.
// Internally, it's based on UUIDv7.
type Sortable[P Prefix] struct{ typedID[P] }

var sortableIDProc = &processor{
	b32Encode: func(u uuid.UUID) string {
		return base32.EncodeLower([16]byte(u))
	},
	b32EncodeTo: func(dst []byte, u uuid.UUID) {
		base32.EncodeLowerTo(dst, [16]byte(u))
	},
	b32Decode: func(s string) (uuid.UUID, error) {
		decoded, err := base32.DecodeLower(s)
		if err != nil {
			return uuid.Nil, err
		}
		return uuid.FromBytes(decoded)
	},
	generateUUID: uuid.NewV7,
}

func (Sortable[P]) processor() *processor {
	return sortableIDProc
}

func (Sortable[P]) Type() string {
	return getPrefix[P]()
}

func (s Sortable[P]) String() string {
	return toString[P](s.uuid, s.processor())
}

func (r Sortable[P]) UUID() uuid.UUID {
	return r.uuid
}

// MarshalText implements the [encoding.TextMarshaler] interface.
// Internally it use [Random.String]
func (r Sortable[P]) MarshalText() ([]byte, error) {
	return marshalText(r)
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
// It parses a TypeID string using [FromString]
func (r *Sortable[P]) UnmarshalText(text []byte) error {
	return unmarshalText(r, text)
}

func (s Sortable[P]) Value() (driver.Value, error) {
	return value(s)
}

func (s *Sortable[P]) Scan(src any) error {
	return scan(s, src)
}

func (s Sortable[P]) TextValue() (pgtype.Text, error) {
	return textValue(s)
}

func (s *Sortable[P]) ScanText(v pgtype.Text) error {
	return scanText(s, v)
}

func (s Sortable[P]) UUIDValue() (pgtype.UUID, error) {
	return uuidValue(s)
}

func (s *Sortable[P]) ScanUUID(v pgtype.UUID) error {
	return scanUUID(s, v)
}
