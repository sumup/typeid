package typeid

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrParse = errors.New("parse typeid")
)

const (
	// MaxPrefixLen is the maximum string length of a [Prefix]. Any generation or parsing of an ID type with a longer prefix will fail.
	MaxPrefixLen = 63
	suffixStrLen = 26 // base32 of UUID
)

type Prefix interface {
	Prefix() string
}

type typedID[P Prefix] struct {
	uuid uuid.UUID
}

func getPrefix[P Prefix]() string {
	var prefix P
	return prefix.Prefix()
}

type instance[P Prefix] interface {
	~struct{ typedID[P] }
	processor() *processor
}

func (tid typedID[P]) Prefix() string {
	return getPrefix[P]()
}

func (tid typedID[P]) UUID() string {
	return tid.uuid.String()
}

// New returns a new TypeID of the specified type with a randomly generated suffix.
//
// Use the 'T' generic argument to indicate your ID type.
//
// Example:
//
//	type UserID = typeid.Sortable[UserPrefix]
//	id, err := typeid.New[UserID]()
func New[T instance[P], P Prefix]() (T, error) {
	tid, err := generate[P](T{}.processor())
	return T{tid}, err
}

// Nil returns the nil identifier for the specified ID type. The nil identifier is a type identifier (typeid) with all corresponding UUID bytes set to zero.
// Functions in this package return the nil identifier in case of errors.
func Nil[T instance[P], P Prefix]() T {
	return T{nilID[P]()}
}

// Must returns a TypeID if the error is nil, otherwise panics. This is a helper function to ease intialization in tests, etc.
// For generating a new id, use [MustNew]
//
// Example:
//
//	testID := typeid.Must(typeid.FromString[UserID]("user_01hf98sp99fs2b4qf2jm11hse4"))
func Must[T any](tid T, err error) T {
	if err != nil {
		panic(err)
	}
	return tid
}

// MustNew returns a generated TypeID if the error is null, otherwise panics.
// Equivalent to:
//
//	typeid.Must(typeid.New[IDType]())
func MustNew[T instance[P], P Prefix]() T {
	return Must(New[T]())
}

func FromString[T instance[P], P Prefix](s string) (T, error) {
	prefix := getPrefix[P]()
	if prefix != "" && !strings.HasPrefix(s, prefix+"_") {
		return Nil[T](), fmt.Errorf("%w: invalid prefix for %T, expected %q", ErrParse, T{}, prefix)
	}

	suffix := strings.TrimPrefix(s, prefix+"_")

	tid, err := from[P](suffix, T{}.processor())
	if err != nil {
		return Nil[T](), fmt.Errorf("%w: invalid suffix %q: %s", ErrParse, suffix, err.Error())
	}
	return T{tid}, nil
}

func FromUUID[T instance[P], P Prefix](u uuid.UUID) (T, error) {
	if err := validatePrefix(getPrefix[P]()); err != nil {
		return Nil[T](), err
	}
	// TODO: Add UUID validation for specific type
	return T{typedID[P]{u}}, nil
}

func FromUUIDStr[T instance[P], P Prefix](uuidStr string) (T, error) {
	u, err := uuid.FromString(uuidStr)
	if err != nil {
		return Nil[T](), fmt.Errorf("%w: uuid from string: %s", ErrParse, err.Error())
	}
	return FromUUID[T](u)
}

func FromUUIDBytes[T instance[P], P Prefix](bytes []byte) (T, error) {
	u, err := uuid.FromBytes(bytes)
	if err != nil {
		return Nil[T](), fmt.Errorf("%w: uuid from bytes: %s", ErrParse, err.Error())
	}
	return FromUUID[T](u)
}
