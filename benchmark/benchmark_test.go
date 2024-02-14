package benchmark_test

import (
	"fmt"
	"testing"

	"github.com/sumup/typeid"
	jpTypeId "go.jetpack.io/typeid"
)

func BenchmarkNew(b *testing.B) {
	b.Run("sumup/typeid", func(b *testing.B) {
		b.Run("Random", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				//nolint:errcheck // Benchmark.
				typeid.New[RandomTestID]()
			}
		})
		b.Run("Sortable", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				//nolint:errcheck // Benchmark.
				typeid.New[SortableTestID]()
			}
		})
	})

	b.Run("jetpack-io/typeid", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			//nolint:errcheck // Benchmark.
			jpTypeId.New[JetpackID]()
		}
	})
}

func BenchmarkString(b *testing.B) {
	b.Run("sumup/typeid", func(b *testing.B) {
		b.Run("Random", func(b *testing.B) {
			b.Run(benchStringRandom(1))
			b.Run(benchStringRandom(8))
			b.Run(benchStringRandom(64))
			b.Run(benchStringRandom(4096))
		})
		b.Run("Sortable", func(b *testing.B) {
			b.Run(benchStringSortable(1))
			b.Run(benchStringSortable(8))
			b.Run(benchStringSortable(64))
			b.Run(benchStringSortable(4096))
		})
	})

	b.Run("jetpack-io/typeid", func(b *testing.B) {
		b.Run(benchStringJp(1))
		b.Run(benchStringJp(8))
		b.Run(benchStringJp(64))
		b.Run(benchStringJp(4096))
	})
}

func BenchmarkFromString(b *testing.B) {
	b.Run("sumup/typeid", func(b *testing.B) {
		b.Run("Random", func(b *testing.B) {
			b.Run(benchFromStringRandom(1))
			b.Run(benchFromStringRandom(8))
			b.Run(benchFromStringRandom(64))
			b.Run(benchFromStringRandom(4096))
		})
		b.Run("Sortable", func(b *testing.B) {
			b.Run(benchFromStringSortable(1))
			b.Run(benchFromStringSortable(8))
			b.Run(benchFromStringSortable(64))
			b.Run(benchFromStringSortable(4096))
		})
	})

	b.Run("jetpack-io/typeid", func(b *testing.B) {
		b.Run(benchFromStringJetpack(1))
		b.Run(benchFromStringJetpack(8))
		b.Run(benchFromStringJetpack(64))
		b.Run(benchFromStringJetpack(4096))
	})
}

func benchStringRandom(n int) (string, func(*testing.B)) {
	ids := makeSortableIDs(n)
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for idx := range ids {
				_ = ids[idx].String()
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func benchStringSortable(n int) (string, func(*testing.B)) {
	ids := makeSortableIDs(n)
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for idx := range ids {
				_ = ids[idx].String()
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func benchStringJp(n int) (string, func(*testing.B)) {
	ids := makeJpTypeIDs(n)
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for idx := range ids {
				_ = ids[idx].String()
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func benchFromStringRandom(n int) (string, func(*testing.B)) {
	idStrings := toString(makeRandomIDs(n))
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for idx := range idStrings {
				//nolint:errcheck // Benchmark.
				typeid.FromString[RandomTestID](idStrings[idx])
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func benchFromStringSortable(n int) (string, func(*testing.B)) {
	idStrings := toString(makeSortableIDs(n))
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			//nolint:errcheck // Benchmark.
			for idx := range idStrings {
				//nolint:errcheck // Benchmark.
				typeid.FromString[SortableTestID](idStrings[idx])
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func benchFromStringJetpack(n int) (string, func(*testing.B)) {
	idStrings := toString(makeJpTypeIDs(n))
	return fmt.Sprintf("n=%d", n), func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			for idx := range idStrings {
				//nolint:errcheck // Benchmark.
				jpTypeId.Parse[JetpackID](idStrings[idx])
			}
		}
		b.ReportMetric(float64(n*b.N)/b.Elapsed().Seconds(), "id/s")
	}
}

func makeRandomIDs(cnt int) []RandomTestID {
	ids := make([]RandomTestID, 0, cnt)
	for i := 0; i < cnt; i++ {
		ids = append(ids, typeid.MustNew[RandomTestID]())
	}
	return ids
}

func makeSortableIDs(cnt int) []SortableTestID {
	ids := make([]SortableTestID, 0, cnt)
	for i := 0; i < cnt; i++ {
		ids = append(ids, typeid.MustNew[SortableTestID]())
	}
	return ids
}

func makeJpTypeIDs(cnt int) []JetpackID {
	ids := make([]JetpackID, 0, cnt)
	for i := 0; i < cnt; i++ {
		ids = append(ids, jpTypeId.Must(jpTypeId.New[JetpackID]()))
	}
	return ids
}

type stringable interface {
	String() string
}

func toString[T stringable](in []T) []string {
	out := make([]string, 0, len(in))
	for _, id := range in {
		out = append(out, id.String())
	}
	return out
}

func BenchmarkEncodeDecode(b *testing.B) {
	b.Run("sumup/typeid", func(b *testing.B) {
		b.Run("Random", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tid := typeid.MustNew[RandomTestID]()
				_ = typeid.Must(typeid.FromString[RandomTestID](tid.String()))
			}
		})
		b.Run("Sortable", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tid := typeid.MustNew[SortableTestID]()
				_ = typeid.Must(typeid.FromString[SortableTestID](tid.String()))
			}
		})
	})

	b.Run("jetpack-io/typeid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tid := jpTypeId.Must(jpTypeId.New[JetpackID]())
			_ = jpTypeId.Must(jpTypeId.Parse[JetpackID](tid.String()))
		}
	})
}

type TestPrefix struct{}

func (TestPrefix) Prefix() string {
	return "test"
}

type (
	RandomTestID   = typeid.Random[TestPrefix]
	SortableTestID = typeid.Sortable[TestPrefix]

	JetpackID struct {
		jpTypeId.TypeID[TestPrefix]
	}
)
