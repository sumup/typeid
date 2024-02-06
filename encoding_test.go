package typeid

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypeID_Pgx_Scan(t *testing.T) {
	t.Parallel()

	userIDStr := "user_01HCJF4N2RER3R6SZHBPFENHVA"
	_ = pgtype.Text{Valid: true, String: userIDStr}

	pgtypeMap := pgtype.NewMap()
	codec := pgtype.TextCodec{}

	t.Run("pgx binary scan", func(t *testing.T) {
		t.Run("valid string bytes", func(t *testing.T) {
			var target UserID
			err := codec.PlanScan(pgtypeMap, pgtype.TextOID, pgtype.BinaryFormatCode, &target).
				Scan([]byte(userIDStr), &target)
			require.NoError(t, err, "can scan id string into typeid.TypeID")

			assert.Equal(t, userIDStr, target.String())
		})

		t.Run("error on invalid char count", func(t *testing.T) {
			var target UserID
			err := codec.PlanScan(pgtypeMap, pgtype.UUIDOID, pgtype.BinaryFormatCode, &target).
				Scan([]byte{23, 12, 125, 54}, &target)
			assert.Error(t, err, "must error if byte count is not 16")
		})
	})

	t.Run("pgx text scan", func(t *testing.T) {
		var target UserID
		err := codec.PlanScan(pgtypeMap, pgtype.TextOID, pgtype.TextFormatCode, &target).
			Scan([]byte(userIDStr), &target)
		require.NoError(t, err, "can scan uuid string into typeid.TypeID")

		assert.Equal(t, userIDStr, target.String(), "scanned typeid.TypeID has the correct uuid string")
	})
	t.Run("error on invalid byte count", func(t *testing.T) {
		var target UserID
		err := codec.PlanScan(pgtypeMap, pgtype.UUIDOID, pgtype.BinaryFormatCode, &target).
			Scan([]byte("2345-2222-111-222"), &target)
		assert.Error(t, err, "must error if byte count is not 16")
	})

	t.Run("error on nil scan", func(t *testing.T) {
		var target UserID
		err := codec.PlanScan(pgtypeMap, pgtype.UUIDOID, pgtype.BinaryFormatCode, &target).
			Scan(nil, &target)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot scan NULL into *typeid.TypeID")
	})
}

func TestTypeID_Pgx_Value(t *testing.T) {
	t.Parallel()

	original, err := New[UserID]()
	require.NoError(t, err)

	pgtypeMap := pgtype.NewMap()
	codec := pgtype.TextCodec{}

	t.Run("binary encoding", func(t *testing.T) {
		var buf []byte
		newBuf, err := codec.PlanEncode(pgtypeMap, pgtype.TextOID, pgtype.BinaryFormatCode, original).
			Encode(original, buf)
		assert.NoError(t, err, "binary encoding should succeed")
		assert.Equal(t, []byte(original.String()), newBuf, "binary encoding should return the uuid bytes")
	})

	t.Run("text encoding", func(t *testing.T) {
		var buf []byte
		newBuf, err := codec.PlanEncode(pgtypeMap, pgtype.TextOID, pgtype.TextFormatCode, original).
			Encode(original, buf)
		assert.NoError(t, err, "binary encoding should succeed")
		assert.Equal(t, []byte(original.String()), newBuf, "text encoding should return the uuid string representation")
	})
}

func TestTypeID_SQL_Scan(t *testing.T) {
	t.Parallel()

	original, err := New[UserID]()
	require.NoError(t, err)
	require.Implements(t, (*sql.Scanner)(nil), &UserID{}, "typeid.TypeID instantiation implements the `sql.Scanner` interface")

	str := original.String()

	otherPrefixID, err := New[AccountID]()
	require.NoError(t, err)

	for _, tc := range []struct {
		name        string
		input       any
		expectedErr bool
	}{
		{
			name:  "scan id string type",
			input: str,
		},
		{
			name:        "fail on invalid type prefix",
			input:       otherPrefixID.String(),
			expectedErr: true,
		},
		{
			name:        "fail on non invalid type",
			input:       12345,
			expectedErr: true,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var scannedID UserID
			err = (&scannedID).Scan(tc.input)
			if tc.expectedErr {
				assert.Error(t, err, "scan should fail")
				return
			}

			assert.NoError(t, err, "scan should succeed")
			assert.Equal(t, original, scannedID, "scanned TypeId is equal to the original one")
		})
	}
}

func TestTypeID_SQL_Value(t *testing.T) {
	t.Parallel()

	id, err := New[UserID]()
	require.NoError(t, err)

	val, err := id.Value()
	assert.NoError(t, err, "value should succeed")
	assert.Equal(t, id.String(), val, "value should return the uuid string")
}

func TestJSON(t *testing.T) {
	str := "account_00041061050r3gg28a1c60t3gf"
	tid := Must(FromString[AccountID](str))

	encoded, err := json.Marshal(tid)
	assert.NoError(t, err)
	assert.Equal(t, `"`+str+`"`, string(encoded))

	var decoded AccountID
	err = json.Unmarshal(encoded, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, tid, decoded)
	assert.Equal(t, str, decoded.String())
}
