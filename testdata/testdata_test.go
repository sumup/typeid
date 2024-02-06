package testdata_test

import (
	_ "embed"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sumup/x/typeid"
	"gopkg.in/yaml.v3"
)

type testPrefix struct{}

func (testPrefix) Prefix() string {
	return "prefix"
}

type nilPrefix struct{}

func (nilPrefix) Prefix() string {
	return ""
}

type TestID = typeid.Sortable[testPrefix]
type NilID = typeid.Sortable[nilPrefix]

var (
	//go:embed valid.yaml
	validYAML []byte
	//go:embed invalid.yaml
	invalidYAML []byte
)

type ValidExample struct {
	Name   string `yaml:"name"`
	Tid    string `yaml:"typeid"`
	Prefix string `yaml:"prefix"`
	UUID   string `yaml:"uuid"`
}

type idImpl interface {
	Type() string
	String() string
	UUID() uuid.UUID
}

func TestValidTestdate(t *testing.T) {
	require.Greater(t, len(validYAML), 0, "testdata/valid.yaml is embedded and not empty")
	var data []ValidExample
	err := yaml.Unmarshal(validYAML, &data)
	require.NoError(t, err, "testdata/valid.yaml is valid yaml")
	require.Greater(t, len(data), 0, "valid.yaml includes test data")

	for i := range data {
		td := data[i]
		t.Run(td.Name, func(t *testing.T) {
			t.Parallel()

			var (
				tid idImpl
				err error
			)

			switch td.Prefix {
			case "prefix":
				tid, err = typeid.FromString[TestID](td.Tid)
			case "":
				tid, err = typeid.FromString[NilID](td.Tid)
			default:
				t.Fatalf("unknown prefix %s in testdata", td.Prefix)
			}

			assert.NoError(t, err, "can parse valid typeid")
			assert.Equal(t, td.Prefix, tid.Type(), "expected prefix matches")
			assert.Equal(t, td.UUID, tid.UUID().String(), "expected UUID string matches")
		})
	}
}

type InvalidExample struct {
	Name        string `yaml:"name"`
	Tid         string `yaml:"typeid"`
	Description string `yaml:"description"`
}

func TestInvalidTestdata(t *testing.T) {
	require.Greater(t, len(invalidYAML), 0, "testdata/invalid.yaml is embedded and not empty")
	var data []InvalidExample
	err := yaml.Unmarshal(invalidYAML, &data)
	require.NoError(t, err, "testdata/invalid.yaml is valid yaml")
	require.Greater(t, len(data), 0, "invalid.yaml includes test data")

	for i := range data {
		td := data[i]
		t.Run(td.Name, func(t *testing.T) {
			t.Parallel()
			_, err := typeid.FromString[TestID](td.Tid)
			assert.Errorf(t, err, "Expected error: %s", td.Description)
		})
	}
}
