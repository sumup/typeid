package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
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
			if err != nil {
				t.Fatalf("can scan id string into typeid.TypeID: unexpected error:\n%+v", err)
			}
			if userIDStr != target.String() {
				t.Errorf("scanned typeid.TypeID has the correct uuid string: expected %s, got %s", userIDStr, target.String())
			}
		})

		t.Run("error on invalid char count", func(t *testing.T) {
			var target UserID
			err := codec.PlanScan(pgtypeMap, pgtype.UUIDOID, pgtype.BinaryFormatCode, &target).
				Scan([]byte{23, 12, 125, 54}, &target)
			if err == nil {
				t.Error("must error if byte count is not 16")
			}
		})
	})

	t.Run("pgx text scan", func(t *testing.T) {
		var target UserID
		err := codec.PlanScan(pgtypeMap, pgtype.TextOID, pgtype.TextFormatCode, &target).
			Scan([]byte(userIDStr), &target)
		if err != nil {
			t.Fatalf("can scan id string into typeid.TypeID: unexpected error:\n%+v", err)
		}
		if userIDStr != target.String() {
			t.Errorf("scanned typeid.TypeID has the correct uuid string: expected %s, got %s", userIDStr, target.String())
		}
	})
	t.Run("error on invalid byte count", func(t *testing.T) {
		var target UserID
		err := codec.PlanScan(pgtypeMap, pgtype.UUIDOID, pgtype.BinaryFormatCode, &target).
			Scan([]byte("2345-2222-111-222"), &target)
		if err == nil {
			t.Error("must error if byte count is not 16")
		}
	})

	t.Run("error on nil scan", func(t *testing.T) {
		var target UserID
		err := codec.PlanScan(pgtypeMap, pgtype.UUIDOID, pgtype.BinaryFormatCode, &target).
			Scan(nil, &target)
		if err == nil {
			t.Error("must error on a nil scan")
		}
		if !strings.Contains(err.Error(), "cannot scan NULL into *typeid.TypeID") {
			t.Error("error must be cannot scan NULL into *typeid.TypeID")
		}
	})
}

func TestTypeID_Pgx_Value(t *testing.T) {
	t.Parallel()

	original, err := New[UserID]()
	if err != nil {
		t.Fatalf("create UserID: unexpected error:\n%+v", err)
	}

	pgtypeMap := pgtype.NewMap()
	codec := pgtype.TextCodec{}

	t.Run("binary encoding", func(t *testing.T) {
		var buf []byte
		newBuf, err := codec.PlanEncode(pgtypeMap, pgtype.TextOID, pgtype.BinaryFormatCode, original).
			Encode(original, buf)
		if err != nil {
			t.Fatalf("binary encoding: unexpected error:\n%+v", err)
		}
		if !bytes.Equal([]byte(original.String()), newBuf) {
			t.Errorf("binary encoding should return the uuid bytes: expected %v, got %v", []byte(original.String()), newBuf)
		}
	})

	t.Run("text encoding", func(t *testing.T) {
		var buf []byte
		newBuf, err := codec.PlanEncode(pgtypeMap, pgtype.TextOID, pgtype.TextFormatCode, original).
			Encode(original, buf)
		if err != nil {
			t.Fatalf("text encoding: unexpected error:\n%+v", err)
		}
		if !bytes.Equal([]byte(original.String()), newBuf) {
			t.Errorf("text encoding should return the uuid string representation: expected %s, got %s", original.String(), newBuf)
		}
	})
}

func TestTypeID_SQL_Scan(t *testing.T) {
	t.Parallel()

	original, err := New[UserID]()
	if err != nil {
		t.Fatalf("create UserID: unexpected error:\n%+v", err)
	}

	scannerType := reflect.TypeOf((*sql.Scanner)(nil)).Elem()
	if !reflect.TypeOf(&UserID{}).Implements(scannerType) {
		t.Fatalf("typeid.TypeID instantiation implements the `sql.Scanner` interface")
	}

	str := original.String()

	otherPrefixID, err := New[AccountID]()
	if err != nil {
		t.Fatalf("create AccountID: unexpected error:\n%+v", err)
	}

	for _, tt := range []struct {
		name       string
		input      any
		shouldFail bool
	}{
		{
			name:  "scan id string type",
			input: str,
		},
		{
			name:       "fail on invalid type prefix",
			input:      otherPrefixID.String(),
			shouldFail: true,
		},
		{
			name:       "fail on non invalid type",
			input:      12345,
			shouldFail: true,
		},
	} {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var scannedID UserID
			err := (&scannedID).Scan(tc.input)
			if tc.shouldFail {
				if err == nil {
					t.Fatalf("scan should fail (expected error %+v)", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("scan should succeed (unexpected error %+v)", err)
			}
			if original != scannedID {
				t.Errorf("scanned TypeId is equal to the original one: expected %v, got %v", original, scannedID)
			}
		})
	}
}

func TestTypeID_SQL_Value(t *testing.T) {
	t.Parallel()

	id, err := New[UserID]()
	if err != nil {
		t.Fatalf("create UserID: unexpected error:\n%+v", err)
	}

	val, err := id.Value()
	if err != nil {
		t.Fatalf("value should succeed (unexpected error %+v)", err)
	}
	if id.String() != val {
		t.Errorf("value should return the uuid string: expected %s, got %s", id.String(), val)
	}
}

func TestJSON(t *testing.T) {
	str := "account_00041061050r3gg28a1c60t3gf"
	tid := Must(FromString[AccountID](str))

	encoded, err := json.Marshal(tid)
	if err != nil {
		t.Fatalf("unexpected error:\n%+v", err)
	}
	if `"`+str+`"` != string(encoded) {
		t.Fatalf("json encoding should return the uuid string: expected %s, got %s", `"`+str+`"`, string(encoded))
	}

	var decoded AccountID
	err = json.Unmarshal(encoded, &decoded)
	if err != nil {
		t.Fatalf("unexpected error:\n%+v", err)
	}

	if tid != decoded {
		t.Errorf("json decoding should return the original uuid: expected %v, got %v", tid, decoded)
	}
	if str != decoded.String() {
		t.Errorf("json decoding should return the original uuid string: expected %s, got %s", str, decoded.String())
	}
}
