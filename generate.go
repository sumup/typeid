package typeid

import (
	"fmt"
	"unsafe"

	"github.com/gofrs/uuid/v5"
)

// processor is an internal structure to handle different types of typedIDs, as they
// may differ in the exact encoding and uuid generator function used.
type processor struct {
	// b32Encode applies a base32 encoding to a UUID.
	b32Encode func(uuid.UUID) string
	// b32EncodeTo applies a base32 encoding to a UUID and copies the result into a provided 26-byte buffer.
	b32EncodeTo func([]byte, uuid.UUID)
	// b32Decode decode a UUID using the resp. base32 decoding.
	b32Decode func(string) (uuid.UUID, error)
	// Generates a new universal unique identifier.
	generateUUID func() (uuid.UUID, error)
}

func from[P Prefix](suffix string, p *processor) (typedID[P], error) {
	var err error

	if err = validatePrefix(getPrefix[P]()); err != nil {
		return nilID[P](), err
	}

	tid := typedID[P]{}
	tid.uuid, err = decodeSuffix(suffix, p)
	if err != nil {
		return nilID[P](), err
	}
	return tid, nil
}

// generate generates
func generate[P Prefix](p *processor) (typedID[P], error) {
	var err error

	if err = validatePrefix(getPrefix[P]()); err != nil {
		return nilID[P](), err
	}

	tid := typedID[P]{}
	tid.uuid, err = p.generateUUID()
	if err != nil {
		return nilID[P](), err
	}
	return tid, nil
}

func nilID[P Prefix]() typedID[P] {
	return typedID[P]{
		uuid: uuid.Nil,
	}
}

func validatePrefix(prefix string) error {
	if prefix == "" {
		return nil
	}

	if len(prefix) > MaxPrefixLen {
		return fmt.Errorf("invalid prefix: %s. Prefix length is %d, expected <= %d", prefix, len(prefix), MaxPrefixLen)
	}

	// Ensure that the prefix has only lowercase ASCII characters
	for _, c := range prefix {
		if c < 'a' || c > 'z' {
			return fmt.Errorf("invalid prefix: '%s'. Prefix should match [a-z]{0,%d}", prefix, MaxPrefixLen)
		}
	}
	return nil
}

func decodeSuffix(suffix string, p *processor) (uuid.UUID, error) {
	if len(suffix) != suffixStrLen {
		return uuid.Nil, fmt.Errorf("invalid suffix: %s. Suffix length is %d, expected %d", suffix, len(suffix), suffixStrLen)
	}

	if suffix[0] > '7' {
		return uuid.Nil, fmt.Errorf("invalid suffix: '%s'. Suffix must start with a 0-7 digit to avoid overflows", suffix)
	}

	return p.b32Decode(suffix)
}

func toString[P Prefix](suffix uuid.UUID, p *processor) string {
	prefix := getPrefix[P]()
	if prefix == "" {
		return p.b32Encode(suffix)
	}

	buf := make([]byte, len(prefix)+1+suffixStrLen)
	copy(buf, prefix)
	copy(buf[len(prefix):], "_")
	p.b32EncodeTo(buf[len(prefix)+1:], suffix)

	return unsafe.String(unsafe.SliceData(buf), len(buf))
}
