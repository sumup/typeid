package typeid

import (
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// idImplementation is a helper constraint asserting the existence of the typeID methods on the type.
type idImplementation[P Prefix] interface {
	instance[P]
	String() string
	UUID() uuid.UUID
}

func marshalText[T idImplementation[P], P Prefix](id T) ([]byte, error) {
	return []byte(id.String()), nil
}

func unmarshalText[T idImplementation[P], P Prefix](dst *T, text []byte) error {
	var err error
	*dst, err = FromString[T](string(text))
	if err != nil {
		return fmt.Errorf("unmarshal text to typeid.TypeID: %w", err)
	}
	return nil
}

func value[T idImplementation[P], P Prefix](id T) (string, error) {
	return id.String(), nil
}

func scan[T idImplementation[P], P Prefix](dst *T, src any) error {
	var err error

	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("scan typeid.Typeid: espected string, got %T", src)
	}

	*dst, err = FromString[T](s)
	if err != nil {
		return fmt.Errorf("scan typeid.TypeID: %w", err)
	}

	return nil
}

func textValue[T idImplementation[P], P Prefix](id T) (pgtype.Text, error) {
	return pgtype.Text{
		String: id.String(),
		Valid:  true,
	}, nil
}

func scanText[T idImplementation[P], P Prefix](dst *T, v pgtype.Text) error {
	var err error

	if !v.Valid {
		return fmt.Errorf("cannot scan NULL into %T", dst)
	}

	*dst, err = FromString[T](v.String)
	if err != nil {
		return fmt.Errorf("scan text to typeid.TypeID: %w", err)
	}

	return nil
}

func uuidValue[T idImplementation[P], P Prefix](id T) (pgtype.UUID, error) {
	return pgtype.UUID{
		Bytes: id.UUID(),
		Valid: true,
	}, nil
}

func scanUUID[T idImplementation[P], P Prefix](dst *T, v pgtype.UUID) error {
	var err error

	if !v.Valid {
		return fmt.Errorf("cannot scan NULL into %T", dst)
	}

	*dst, err = FromUUIDBytes[T](v.Bytes[:])
	if err != nil {
		return fmt.Errorf("scan UUID to typeid.TypeID: %w", err)
	}

	return nil
}

type prefixScanner[P Prefix] struct{}

func (prefixScanner[P]) ScanText(v pgtype.Text) error {
	if !v.Valid {
		return fmt.Errorf("cannot scan NULL prefix")
	}

	var p P
	if v.String != p.Prefix() {
		return fmt.Errorf(
			"scan composite typeid: prefix mismatch: got %q, expected %q",
			v.String,
			p.Prefix(),
		)
	}

	return nil
}

// compositeUUIDScanner sets the UUID field of a PostgreSQL typeid composite.
type compositeUUIDScanner[T idImplementation[P], P Prefix] struct {
	dst *T
}

func (s compositeUUIDScanner[T, P]) ScanUUID(v pgtype.UUID) error {
	return scanUUID(s.dst, v)
}

func compositeIsNull() bool {
	return false
}

func compositeIndex[T idImplementation[P], P Prefix](id T, i int) any {
	switch i {
	case 0:
		return getPrefix[P]()
	case 1:
		return pgtype.UUID{
			Bytes: id.UUID(),
			Valid: true,
		}
	default:
		panic(fmt.Errorf("illegal composite index %d", i))
	}
}

func compositeScanNull[T idImplementation[P], P Prefix](dst *T) error {
	return fmt.Errorf("cannot scan NULL into %T", dst)
}

func compositeScanIndex[T idImplementation[P], P Prefix](dst *T, i int) any {
	switch i {
	case 0:
		return new(prefixScanner[P])
	case 1:
		return compositeUUIDScanner[T, P]{dst: dst}
	default:
		panic(fmt.Errorf("illegal composite scan index %d", i))
	}
}
