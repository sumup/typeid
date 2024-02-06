package main

import (
	"errors"
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

var (
	errNilScan = errors.New("cannot scan NULL into *typeid.TypeID")
)

func textValue[T idImplementation[P], P Prefix](id T) (pgtype.Text, error) {
	return pgtype.Text{
		String: id.String(),
		Valid:  true,
	}, nil
}

func scanText[T idImplementation[P], P Prefix](dst *T, v pgtype.Text) error {
	var err error

	if !v.Valid {
		return errNilScan
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
		return errNilScan
	}

	*dst, err = FromUUIDBytes[T](v.Bytes[:])
	if err != nil {
		return fmt.Errorf("scan UUID to typeid.TypeID: %w", err)
	}

	return nil
}
