package typeid

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/gofrs/uuid/v5"
)

const (
	userIDPrefix    = "user"
	accountIDPrefix = "system_account"

	emptyID = "00000000000000000000000000" // empty is suffix is 26 zeros
)

type userPrefix struct{}

func (userPrefix) Prefix() string {
	return userIDPrefix
}

type accountPrefix struct{}

func (accountPrefix) Prefix() string {
	return accountIDPrefix
}

type nilPrefix struct{}

func (nilPrefix) Prefix() string {
	return ""
}

type UserID = Random[userPrefix]
type AccountID = Sortable[accountPrefix]
type NilID = Random[nilPrefix]

func TestTypeID_New(t *testing.T) {
	t.Parallel()

	userID, err := New[UserID]()
	if err != nil {
		t.Fatalf("create UserID: unexpected error:\n%+v", err)
	}
	if uuid.V4 != userID.UUID().Version() {
		t.Errorf("expected UUIDv4, got version byte: %x", userID.UUID().Version())
	}

	accountID, err := New[AccountID]()
	if err != nil {
		t.Fatalf("create AserID: unexpected error:\n%+v", err)
	}
	if uuid.V7 != accountID.UUID().Version() {
		t.Errorf("expected UUIDv7, got version byte: %x", accountID.UUID().Version())
	}
}

func TestTypeID_Nil(t *testing.T) {
	t.Parallel()

	nilUserID := Nil[UserID]()
	if "user_"+emptyID != nilUserID.String() {
		t.Errorf("expected nil user id, got: %s", nilUserID.String())
	}
	if nilUserID != Nil[UserID]() {
		t.Errorf("two nil id's are equal\nGot: %v\nExpected: %v", nilUserID, Nil[UserID]())
	}

	nilAccountID := Nil[AccountID]()
	if "system_account_"+emptyID != nilAccountID.String() {
		t.Errorf("expected nil account id, got: %s", nilAccountID.String())
	}

	nilID := Nil[NilID]()
	if emptyID != nilID.String() {
		t.Errorf("two nil id's are equal\nGot: %v\nExpected: %v", emptyID, nilID.String())
	}
}

func TestTypeID_ToFrom(t *testing.T) {
	t.Parallel()

	type RandomID = UserID
	t.Run("typeid.Random", runToFromQuickTests[RandomID])

	type SortableID = AccountID
	t.Run("typeid.Sortable", runToFromQuickTests[SortableID])
}

func runToFromQuickTests[T idImplementation[P], P Prefix](t *testing.T) {
	t.Helper()
	t.Parallel()

	t.Run("from string", func(t *testing.T) {
		t.Parallel()
		if err := quick.Check(fromStringTester[T](t), nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("from UUID", func(t *testing.T) {
		t.Parallel()
		if err := quick.Check(fromUUIDTester[T](t), nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("from UUID string", func(t *testing.T) {
		t.Parallel()
		if err := quick.Check(fromUUIDStringTester[T](t), nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("from UUID bytes", func(t *testing.T) {
		t.Parallel()
		if err := quick.Check(fromUUIDBytesTester[T](t), nil); err != nil {
			t.Error(err)
		}
	})
}

func fromStringTester[T idImplementation[P], P Prefix](t *testing.T) func(wid wrappedID[T, P]) bool {
	t.Helper()
	return func(wid wrappedID[T, P]) bool {
		parsedID, err := FromString[T](wid.ID().String())
		if err != nil {
			t.Fatalf("parse type id from string: unexpected error:\n%+v", err)
		}
		return wid.ID() == parsedID
	}
}

func fromUUIDTester[T idImplementation[P], P Prefix](t *testing.T) func(wid wrappedID[T, P]) bool {
	t.Helper()
	return func(wid wrappedID[T, P]) bool {
		parsedID, err := FromUUID[T](wid.ID().UUID())
		if err != nil {
			t.Fatalf("parse type id from UUID: unexpected error:\n%+v", err)
		}
		return wid.ID() == parsedID
	}
}

func fromUUIDStringTester[T idImplementation[P], P Prefix](t *testing.T) func(wid wrappedID[T, P]) bool {
	t.Helper()
	return func(wid wrappedID[T, P]) bool {
		parsedID, err := FromUUIDStr[T](wid.ID().UUID().String())
		if err != nil {
			t.Fatalf("parse type id from UUID string: unexpected error:\n%+v", err)
		}
		return wid.ID() == parsedID
	}
}

func fromUUIDBytesTester[T idImplementation[P], P Prefix](t *testing.T) func(wid wrappedID[T, P]) bool {
	t.Helper()
	return func(wid wrappedID[T, P]) bool {
		parsedID, err := FromUUIDBytes[T](wid.id.UUID().Bytes())
		if err != nil {
			t.Fatalf("parse type id from UUID bytes: unexpected error:\n%+v", err)
		}
		return wid.ID() == parsedID
	}
}

// wrappedID is a helper struct that implements [quick.Generator] to facilitate simple property-based tests.
// We require this helper because we cannot directly implement interfaces on type aliases.
type wrappedID[T instance[P], P Prefix] struct {
	id T
}

func (w wrappedID[T, P]) ID() T {
	return w.id
}

func (w wrappedID[T, P]) Generate(rnd *rand.Rand, _ int) reflect.Value {
	// gen the processor to determine the UUID version to use
	procGenUUID, err := (T{}).processor().generateUUID()
	if err != nil {
		panic(err)
	}
	version := procGenUUID.Version()

	uuidGen := uuid.NewGenWithOptions(uuid.WithRandomReader(rnd))
	var uid uuid.UUID

	switch version {
	case uuid.V4:
		if uid, err = uuidGen.NewV4(); err != nil {
			panic("failed to generate uuid v4")
		}
	case uuid.V7:
		if uid, err = uuidGen.NewV7(); err != nil {
			panic("failed to generate uuid v7")
		}
	default:
		panic("unknown uuid version")
	}

	tid, err := FromUUID[T](uid)
	if err != nil {
		panic(err)
	}

	return reflect.ValueOf(wrappedID[T, P]{tid})
}
