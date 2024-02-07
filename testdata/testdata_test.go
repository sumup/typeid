package testdata_test

import (
	_ "embed"
	"strings"
	"sync"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/sumup/typeid"
)

type testPrefix struct{}

func (testPrefix) Prefix() string {
	return "prefix"
}

type nilPrefix struct{}

func (nilPrefix) Prefix() string {
	return ""
}

type (
	RandomTestID   = typeid.Random[testPrefix]
	SortableTestID = typeid.Sortable[testPrefix]

	RandomNilID   = typeid.Random[nilPrefix]
	SortableNilID = typeid.Sortable[nilPrefix]
)

type idImpl interface {
	Type() string
	String() string
	UUID() uuid.UUID
}

func TestValidIDs(t *testing.T) {
	t.Parallel()

	type testcase struct {
		Name   string
		TypeID string
		Prefix string
		UUID   string
	}

	assertions := func(t *testing.T, tc testcase, tid idImpl, err error) {
		t.Helper()
		if err != nil {
			t.Errorf("unexpected error: cannot parse valid typeid: %v", err.Error())
		}
		if tc.Prefix != tid.Type() {
			t.Errorf("type id prefix does not match expectected value:\nExpected:%s\nGot: %s)", tc.Prefix, tid.Type())
		}
		if tc.UUID != tid.UUID().String() {
			t.Errorf("type id UUID does not match expectected value:\nExpected:%s\nGot: %s)", tc.UUID, tid.UUID().String())
		}
	}

	t.Run("common", func(t *testing.T) {
		t.Parallel()

		for _, tt := range []testcase{
			{
				Name:   "nil id (no prefix, all bytes zero)",
				TypeID: "00000000000000000000000000",
				UUID:   "00000000-0000-0000-0000-000000000000",
			},
			{
				Name:   "one",
				TypeID: "00000000000000000000000001",
				UUID:   "00000000-0000-0000-0000-000000000001",
			},
			{
				Name:   "thirtytwo (typeid string uses base32)",
				TypeID: "00000000000000000000000010",
				UUID:   "00000000-0000-0000-0000-000000000020",
			},
		} {
			tc := tt
			t.Run(tc.Name, func(t *testing.T) {
				t.Parallel()

				sid, err := typeid.FromString[SortableNilID](tc.TypeID)
				assertions(t, tc, sid, err)

				rid, err := typeid.FromString[RandomNilID](tc.TypeID)
				assertions(t, tc, rid, err)
			})
		}

		t.Run("Sortable", func(t *testing.T) {
			t.Parallel()
			for _, tt := range []testcase{
				{
					Name:   "ten",
					TypeID: "0000000000000000000000000a",
					UUID:   "00000000-0000-0000-0000-00000000000a",
				},
				{
					Name:   "sixteen",
					TypeID: "0000000000000000000000000g",
					UUID:   "00000000-0000-0000-0000-000000000010",
				},

				{
					Name:   "maximum",
					TypeID: "7zzzzzzzzzzzzzzzzzzzzzzzzz",
					UUID:   "ffffffff-ffff-ffff-ffff-ffffffffffff",
				},
				{
					Name:   "full alphabet",
					TypeID: "prefix_0123456789abcdefghjkmnpqrs",
					Prefix: "prefix",
					UUID:   "0110c853-1d09-52d8-d73e-1194e95b5f19",
				},
				{
					Name:   "UUIDv7",
					TypeID: "prefix_01hp1aybq6f6athhfcvp1j8fpt",
					Prefix: "prefix",
					UUID:   "018d82af-2ee6-7995-a8c5-ecdd83243eda",
				},
			} {
				tc := tt
				t.Run(tc.Name, func(t *testing.T) {
					t.Parallel()
					var (
						tid idImpl
						err error
					)

					switch tc.Prefix {
					case "prefix":
						tid, err = typeid.FromString[SortableTestID](tc.TypeID)
					case "":
						tid, err = typeid.FromString[SortableNilID](tc.TypeID)
					default:
						t.Fatalf("unknown prefix %s in testdata", tc.Prefix)
					}
					assertions(t, tc, tid, err)
				})
			}
		})

		t.Run("Random", func(t *testing.T) {
			t.Parallel()
			for _, tt := range []testcase{
				{
					Name:   "ten",
					TypeID: "0000000000000000000000000A",
					UUID:   "00000000-0000-0000-0000-00000000000a",
				},
				{
					Name:   "sixteen",
					TypeID: "0000000000000000000000000G",
					UUID:   "00000000-0000-0000-0000-000000000010",
				},

				{
					Name:   "maximum",
					TypeID: "7ZZZZZZZZZZZZZZZZZZZZZZZZZ",
					UUID:   "ffffffff-ffff-ffff-ffff-ffffffffffff",
				},
				{
					Name:   "full alphabet",
					TypeID: "prefix_0123456789ABCDEFGHJKMNPQRS",
					Prefix: "prefix",
					UUID:   "0110c853-1d09-52d8-d73e-1194e95b5f19",
				},
				{
					Name:   "UUIDv7",
					TypeID: "prefix_01HP1AYBQ6F6ATHHFCVP1J8FPT",
					Prefix: "prefix",
					UUID:   "018d82af-2ee6-7995-a8c5-ecdd83243eda",
				},
			} {
				tc := tt
				t.Run(tc.Name, func(t *testing.T) {
					t.Parallel()
					var (
						tid idImpl
						err error
					)

					switch tc.Prefix {
					case "prefix":
						tid, err = typeid.FromString[RandomTestID](tc.TypeID)
					case "":
						tid, err = typeid.FromString[RandomNilID](tc.TypeID)
					default:
						t.Fatalf("unknown prefix %s in testdata", tc.Prefix)
					}
					assertions(t, tc, tid, err)
				})
			}
		})

	})

	for _, tt := range []testcase{

		{
			Name:   "ten",
			TypeID: "0000000000000000000000000a",
			UUID:   "00000000-0000-0000-0000-00000000000a",
		},
		{
			Name:   "sixteen",
			TypeID: "0000000000000000000000000g",
			UUID:   "00000000-0000-0000-0000-000000000010",
		},

		{
			Name:   "maximum",
			TypeID: "7zzzzzzzzzzzzzzzzzzzzzzzzz",
			UUID:   "ffffffff-ffff-ffff-ffff-ffffffffffff",
		},
		{
			Name:   "full alphabet",
			TypeID: "prefix_0123456789abcdefghjkmnpqrs",
			Prefix: "prefix",
			UUID:   "0110c853-1d09-52d8-d73e-1194e95b5f19",
		},
		{
			Name:   "UUIDv7",
			TypeID: "prefix_01hp1aybq6f6athhfcvp1j8fpt",
			Prefix: "prefix",
			UUID:   "018d82af-2ee6-7995-a8c5-ecdd83243eda",
		},
	} {
		tc := tt
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			var (
				tid idImpl
				err error
			)

			switch tc.Prefix {
			case "prefix":
				tid, err = typeid.FromString[SortableTestID](tc.TypeID)
			case "":
				tid, err = typeid.FromString[SortableNilID](tc.TypeID)
			default:
				t.Fatalf("unknown prefix %s in testdata", tc.Prefix)
			}

			if err != nil {
				t.Errorf("unexpected error: cannot parse valid typeid: %v", err.Error())
			}
			if tc.Prefix != tid.Type() {
				t.Errorf("type id prefix does not match expectected value:\nExpected:%s\nGot: %s)", tc.Prefix, tid.Type())
			}
			if tc.UUID != tid.UUID().String() {
				t.Errorf("type id UUID does not match expectected value:\nExpected:%s\nGot: %s)", tc.UUID, tid.UUID().String())
			}
		})
	}
}

// For testing purposes, we use a gloabal variable to enable the dynamicPrefix type to return different prefixes.
// Multiple test cases using this type cannot run in parallel.
var (
	currentDynamicPrefixMu sync.RWMutex
	currentDynamicPrefix   string
)

type dynamicPrefix struct{}

func (dynamicPrefix) Prefix() string {
	currentDynamicPrefixMu.RLock()
	defer currentDynamicPrefixMu.RUnlock()
	return currentDynamicPrefix
}

func SetDynamicPrefix(p string) {
	currentDynamicPrefixMu.Lock()
	currentDynamicPrefix = p
	currentDynamicPrefixMu.Unlock()
}

type DynamicID = typeid.Sortable[dynamicPrefix]

func TestInvalidIDs(t *testing.T) {
	t.Parallel()

	type testcase struct {
		Name        string
		Prefix      string
		TypeID      string
		ErrorReason string
	}

	cases := []testcase{
		{
			Name:        "prefix contains uppercase characters",
			Prefix:      "Prefix",
			TypeID:      "Prefix_00000000000000000000000000",
			ErrorReason: "Only lowercase letters are allowed in the prefix. No uppercase letters.",
		},
		{
			Name:        "prefix contains digits",
			Prefix:      "pref1x",
			TypeID:      "pref1x_00000000000000000000000000",
			ErrorReason: "Only lowercase letters are allowed in the prefix. No digits.",
		},
		{
			Name:        "prefix contains non alphabetic characters",
			Prefix:      "pre,fix",
			TypeID:      "pre,fix_00000000000000000000000000",
			ErrorReason: "Only lowercase letters are allowed in the prefix. No non-alphabetic characters.",
		},
		{
			Name:        "prefix contains a separator",
			Prefix:      "pre_fix",
			TypeID:      "pre_fix_00000000000000000000000000",
			ErrorReason: "Only lowercase letters are allowed in the prefix. No underscores. They act as separator for the two ID parts",
		},
		{
			Name:        "prefix is empty",
			Prefix:      " prefix",
			TypeID:      " prefix_00000000000000000000000000",
			ErrorReason: "Only lowercase letters are allowed in the prefix. No leading or trailing spaces.",
		},
		{
			Name:        "prefix to long",
			Prefix:      "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkl",
			TypeID:      "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkl_00000000000000000000000000",
			ErrorReason: "Maximum prefix length is 63 characters",
		},
		{
			Name:        "No separator if prefix is empty",
			TypeID:      "_00000000000000000000000000",
			ErrorReason: "A typeid with an empty prefix must not contain a seperator",
		},
		{
			Name:        "Underscore",
			TypeID:      "_",
			ErrorReason: "A seperator is not a valid typeid",
		},
		{
			Name:        "ID part to short",
			Prefix:      "prefix",
			TypeID:      "prefix_1234567890123456789012345",
			ErrorReason: "The ID part must contain exactly 26 characters",
		},
		{
			Name:        "ID part to long",
			Prefix:      "prefix",
			TypeID:      "prefix_123456789012345678901234567",
			ErrorReason: "The ID part must contain exactly 26 characters",
		},
		{
			Name:        "ID part contains spaces",
			Prefix:      "prefix",
			TypeID:      "prefix_1234567890123456789012345 ",
			ErrorReason: "The ID part must not contain spaces",
		},
		{
			Name:        "Suffix with hyphens",
			Prefix:      "prefix",
			TypeID:      "prefix_123456789-123456789-123456",
			ErrorReason: "The suffix must not contain hypens",
		},
		{
			Name:        "Suffix from wrong alphabet",
			Prefix:      "prefix",
			TypeID:      "prefix_ooooooiiiiiiuuuuuuulllllll",
			ErrorReason: "The suffix must only contain characters from the crockford base32 alphabnet",
		},
		{
			Name:        "Suffic not crockford",
			Prefix:      "prefix",
			TypeID:      "prefix_i23456789ol23456789oi23456",
			ErrorReason: "The suffix must not contain characters excluded by the crockford base32 rules",
		},
		{
			Name:        "No crockford hyphenation",
			Prefix:      "prefix",
			TypeID:      "prefix_123456789-0123456789-0123456",
			ErrorReason: "The suffix must not contain hyphens, even though they would be ignored under the crockford rules",
		},
		{
			Name:        "overflow",
			Prefix:      "prefix",
			TypeID:      "prefix_8zzzzzzzzzzzzzzzzzzzzzzzzz",
			ErrorReason: "The suffix encodes at most 128 bits",
		},
	}

	t.Run("Sortable", func(t *testing.T) {
		type dynamicSortableID = typeid.Sortable[dynamicPrefix]

		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				SetDynamicPrefix(tc.Prefix)
				_, err := typeid.FromString[dynamicSortableID](tc.TypeID)
				if err == nil {
					t.Fatalf("expected an error, but got nil:Input: %s\nError reason: %s\n", tc.TypeID, tc.ErrorReason)
				}
			})
		}
	})

	t.Run("Random", func(t *testing.T) {
		type dynamicRandomID = typeid.Random[dynamicPrefix]

		for _, tc := range cases {
			suffix, ok := strings.CutPrefix(tc.TypeID, tc.Prefix+"_")
			if !ok {
				t.Fatalf("could not cut prefix %s from typeid: %s", tc.Prefix, tc.TypeID)
			}
			// Reassign typeid with uppercase letters uppercase
			tc.TypeID = tc.Prefix + "_" + strings.ToUpper(suffix)
			t.Run(tc.Name, func(t *testing.T) {
				SetDynamicPrefix(tc.Prefix)
				_, err := typeid.FromString[dynamicRandomID](tc.TypeID)
				if err == nil {
					t.Fatalf("expected an error, but got nil:Input: %s\nError reason: %s\n", tc.TypeID, tc.ErrorReason)
				}
			})
		}
	})
}
