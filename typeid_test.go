package typeid

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	userIDPrefix    = "user"
	accountIDPrefix = "account"

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
	require.NoError(t, err, "can create userid")
	assert.Equal(t, uuid.V4, userID.UUID().Version())

	accountID, err := New[AccountID]()
	require.NoError(t, err, "can create accountID")
	assert.Equal(t, uuid.V7, accountID.UUID().Version())
}

func TestTypeID_Nil(t *testing.T) {
	t.Parallel()

	nilUserID := Nil[UserID]()
	assert.Equal(t, "user_"+emptyID, nilUserID.String())
	assert.Equal(t, nilUserID, Nil[UserID](), "two nil id's are equal")
	nilAccountID := Nil[AccountID]()
	assert.Equal(t, "account_"+emptyID, nilAccountID.String())

	nilID := Nil[NilID]()
	assert.Equal(t, emptyID, nilID.String())
}

func TestTypeID_ToFrom(t *testing.T) {
	t.Parallel()

	type RandomID = UserID
	t.Run("typeid.Random", runToFromQuickTests[RandomID])

	type SortableID = AccountID
	t.Run("typeid.Sortable", runToFromQuickTests[SortableID])
}

func runToFromQuickTests[T idImplementation[P], P Prefix](t *testing.T) {
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
	return func(wid wrappedID[T, P]) bool {
		parsedID, err := FromString[T](wid.ID().String())
		require.NoError(t, err, "can parse typeid from string")
		return wid.ID() == parsedID
	}
}

func fromUUIDTester[T idImplementation[P], P Prefix](t *testing.T) func(wid wrappedID[T, P]) bool {
	return func(wid wrappedID[T, P]) bool {
		parsedID, err := FromUUID[T](wid.ID().UUID())
		require.NoError(t, err, "can parse typeid from string")
		return wid.ID() == parsedID
	}
}

func fromUUIDStringTester[T idImplementation[P], P Prefix](t *testing.T) func(wid wrappedID[T, P]) bool {
	return func(wid wrappedID[T, P]) bool {
		parsedID, err := FromUUIDStr[T](wid.ID().UUID().String())
		require.NoError(t, err, "can parse typeid from uuid string")
		return wid.ID() == parsedID
	}
}

func fromUUIDBytesTester[T idImplementation[P], P Prefix](t *testing.T) func(wid wrappedID[T, P]) bool {
	return func(wid wrappedID[T, P]) bool {
		parsedID, err := FromUUIDBytes[T](wid.id.UUID().Bytes())
		require.NoError(t, err, "can parse typeid from uuid bytes")
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

func (w wrappedID[T, P]) Generate(rand *rand.Rand, _ int) reflect.Value {
	// gen the processor to determine the UUID version to use
	procGenUUID, err := (T{}).processor().generateUUID()
	if err != nil {
		panic(err)
	}
	version := procGenUUID.Version()

	uuidGen := uuid.NewGenWithOptions(uuid.WithRandomReader(rand))
	var uid uuid.UUID

	if version == uuid.V4 {
		if uid, err = uuidGen.NewV4(); err != nil {
			panic("failed to generate uuid v4")
		}
	} else if version == uuid.V7 {
		if uid, err = uuidGen.NewV7(); err != nil {
			panic("failed to generate uuid v7")
		}
	} else {
		panic("unknown uuid version")
	}

	tid, err := FromUUID[T](uid)
	if err != nil {
		panic(err)
	}

	return reflect.ValueOf(wrappedID[T, P]{tid})
}
