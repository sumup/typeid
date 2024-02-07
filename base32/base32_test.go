package base32

import (
	"crypto/rand"
	"encoding/base32"
	"testing"
	"testing/quick"
)

func TestEncodeDecode(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name          string
		stdLibEncoder *base32.Encoding
		EncodeFunc    func(src [16]byte) string
		DecodeFunc    func(src string) ([]byte, error)
	}{
		{
			name:          "upppercase encoding",
			stdLibEncoder: base32.NewEncoding("0123456789ABCDEFGHJKMNPQRSTVWXYZ"),
			EncodeFunc:    EncodeUpper,
			DecodeFunc:    DecodeUpper,
		},
		{
			name:          "lowercase encoding",
			stdLibEncoder: base32.NewEncoding("0123456789abcdefghjkmnpqrstvwxyz"),
			EncodeFunc:    EncodeLower,
			DecodeFunc:    DecodeLower,
		},
	} {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for i := 0; i < 1000; i++ {
				// Generate 16 random bytes
				data := make([]byte, 16)
				_, err := rand.Read(data)
				if err != nil {
					t.Fatalf("unexpected error:\n%+v", err)
				}

				// Encode them using our library, and encode them using go's standard library:
				actual := tc.EncodeFunc([16]byte(data))

				// The standard base32 library decodes in groups of 5 bytes, otherwise it needs
				// to pad, by default it pads at the end of the byte array, but to match our
				// encoding we need to pad in the front.
				// Pad manually, and then remove the extra 000000 from the resulting string.
				padded := append([]byte{0x00, 0x00, 0x00, 0x00}, data...)
				expected := tc.stdLibEncoder.EncodeToString(padded)[6:]

				// They should be equal
				if expected != actual {
					t.Errorf("padded value is not equal to exted:\nExpected: %s\nActual: %s", expected, actual)
				}

				// Decoding again should yield the original result:
				decoded, err := tc.DecodeFunc(actual)
				if err != nil {
					t.Errorf("unexpected error during decode:\n%+v", err)
				}
				for i := 0; i < 16; i++ {
					if data[i] != decoded[i] {
						t.Errorf("decoded value is not equal to original:\nExpected: %v\nActual: %v", data[i], decoded[i])
					}
				}
			}
		})
	}

	encoder := base32.NewEncoding("0123456789ABCDEFGHJKMNPQRSTVWXYZ")

	for i := 0; i < 1000; i++ {
		// Generate 16 random bytes
		data := make([]byte, 16)
		_, err := rand.Read(data)
		if err != nil {
			t.Fatalf("unexpected error:\n%+v", err)
		}

		// Encode them using our library, and encode them using go's standard library:
		actual := EncodeUpper([16]byte(data))

		// The standard base32 library decodes in groups of 5 bytes, otherwise it needs
		// to pad, by default it pads at the end of the byte array, but to match our
		// encoding we need to pad in the front.
		// Pad manually, and then remove the extra 000000 from the resulting string.
		padded := append([]byte{0x00, 0x00, 0x00, 0x00}, data...)
		expected := encoder.EncodeToString(padded)[6:]

		// They should be equal
		if expected != actual {
			t.Errorf("padded value is not equal to exted:\nExpected: %s\nActual: %s", expected, actual)
		}

		// Decoding again should yield the original result:
		decoded, err := DecodeUpper(actual)
		if err != nil {
			t.Errorf("unexpected error during decode:\n%+v", err)
		}
		for i := 0; i < 16; i++ {
			if data[i] != decoded[i] {
				t.Errorf("decoded value is not equal to original:\nExpected: %v\nActual: %v", data[i], decoded[i])
			}
		}
	}
}

func TestEncodeDecodeProp(t *testing.T) {
	t.Parallel()

	f := func(input [16]byte) bool {
		enc := EncodeUpper(input)
		dec, err := DecodeUpper(enc)
		if err != nil {
			t.Errorf("decode: received unexpected error:\n%+v", err)
		}
		return input == [16]byte(dec)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestAlphabetValidity(t *testing.T) {
	for i := range alphUp {
		if decUpper[alphUp[i]] == 0xFF {
			t.Errorf("char from alphabet not in base64 lookup table (%c)", alphUp[i])
		}
	}
}
