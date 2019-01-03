package kazaam_test

import (
	"testing"

	"github.com/mbordner/kazaam"
)

const (
	benchmarkJSON = `{"topKeyA": {"arrayKey": [{"foo": 0}, {"foo": 1}, {"foo": 1}, {"foo": 2}], "notArrayKey": "Sun Jul 23 08:15:27 +0000 2017", "deepArrayKey": [{"key0":["val0", "val1"]}]}, "topKeyB":{"nextKeyB": "valueB"}}`
)

var (
	benchmarkSlice = []byte(benchmarkJSON)
)

// Just for emulating field access, so it will not throw "evaluated but not
// used." Borrowed from:
// https://github.com/qntfy/jsonparser/blob/master/benchmark/benchmark_small_payload_test.go
func nothing(_ ...interface{}) {}

func BenchmarkShift(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "shift", "spec": {"outputKey": "topKeyA.arrayKey"}}]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transform, _ := kazaam.NewKazaam(spec)
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}
}

func BenchmarkShiftTranformOnly(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "shift", "spec": {"outputKey": "topKeyA.arrayKey"}}]`
	transform, _ := kazaam.NewKazaam(spec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}
}

func BenchmarkShiftWithWildcard(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "shift", "spec": {"outputKey": "topKeyA.arrayKey[*].foo"}}]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transform, _ := kazaam.NewKazaam(spec)
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}
}

func BenchmarkShiftWithWildcardTransformOnly(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "shift", "spec": {"outputKey": "topKeyA.arrayKey[*].foo"}}]`
	transform, _ := kazaam.NewKazaam(spec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}
}

func BenchmarkDeepShiftWithWildcard(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "shift", "spec": {"outputKey": "topKeyA.deepArrayKey[0].key0[*]"}}]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transform, _ := kazaam.NewKazaam(spec)
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}
}

func BenchmarkDeepShiftWithWildcardTrnsformOnly(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "shift", "spec": {"outputKey": "topKeyA.deepArrayKey[0].key0[*]"}}]`
	transform, _ := kazaam.NewKazaam(spec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}
}

func BenchmarkShiftEncapsulateTransform(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "shift", "spec": {"data": ["$"]}}]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transform, _ := kazaam.NewKazaam(spec)
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}
}

func BenchmarkShiftEncapsulateTransformTransformOnly(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "shift", "spec": {"data": ["$"]}}]`
	transform, _ := kazaam.NewKazaam(spec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}
}

func BenchmarkConcat(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "concat", "spec": {"sources": [{"value": "TEST"}, {"path": "topKeyA.notArrayKey"}], "targetPath": "a.output", "delim": "," }}]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transform, _ := kazaam.NewKazaam(spec)
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}
}

func BenchmarkConcatTransformOnly(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "concat", "spec": {"sources": [{"value": "TEST"}, {"path": "topKeyA.notArrayKey"}], "targetPath": "a.output", "delim": "," }}]`
	transform, _ := kazaam.NewKazaam(spec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}
}

func BenchmarkConcatWithWildcardNested(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "concat", "spec": {"sources": [{"value": "TEST"}, {"path": "topKeyA.arrayKey[*].foo"}], "targetPath": "a.output", "delim": "," }}]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transform, _ := kazaam.NewKazaam(spec)
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}
}

func BenchmarkConcatWithWildcardNestedTransformOnly(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "concat", "spec": {"sources": [{"value": "TEST"}, {"path": "topKeyA.arrayKey[*].foo"}], "targetPath": "a.output", "delim": "," }}]`
	transform, _ := kazaam.NewKazaam(spec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}
}

func BenchmarkTransformMultiOpWithOver(b *testing.B) {
	b.ReportAllocs()

	spec := `[{
        "operation": "concat",
        "over": "topKeyA.arrayKey",
        "spec": {"sources": [{"path": "foo"}, {"value": "KEY"}], "targetPath": "url", "delim": ":" }
    }, {
        "operation": "shift",
        "spec": {"urls": "topKeyA.arrayKey[*].url" }
    }]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transform, _ := kazaam.NewKazaam(spec)
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}
}

func BenchmarkTransformMultiOpWithOverTransformOnly(b *testing.B) {
	b.ReportAllocs()

	spec := `[{
        "operation": "concat",
        "over": "topKeyA.arrayKey",
        "spec": {"sources": [{"path": "foo"}, {"value": "KEY"}], "targetPath": "url", "delim": ":" }
    }, {
        "operation": "shift",
        "spec": {"urls": "topKeyA.arrayKey[*].url" }
    }]`
	transform, _ := kazaam.NewKazaam(spec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}
}
func BenchmarkCoalesce(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "coalesce", "spec": {"foo": ["key.foo", "topKeyB.nextKeyB"]}}]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transform, _ := kazaam.NewKazaam(spec)
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}

}

func BenchmarkCoalesceTransformOnly(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "coalesce", "spec": {"foo": ["key.foo", "topKeyB.nextKeyB"]}}]`
	transform, _ := kazaam.NewKazaam(spec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}

}

func BenchmarkExtract(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "extract", "spec": {"path": "topKeyB"}}]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transform, _ := kazaam.NewKazaam(spec)
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}

}

func BenchmarkExtractTransformOnly(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "extract", "spec": {"path": "topKeyB"}}]`
	transform, _ := kazaam.NewKazaam(spec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}

}

func BenchmarkTimestamp(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "timestamp", "spec": {"notArrayKey":{"inputFormat":"Mon Jan _2 15:04:05 -0700 2006","outputFormat":"2006-01-02T15:04:05-0700}}}]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transform, _ := kazaam.NewKazaam(spec)
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}

}

func BenchmarkTimestampTransformOnly(b *testing.B) {
	b.ReportAllocs()

	spec := `[{"operation": "timestamp", "spec": {"notArrayKey":{"inputFormat":"Mon Jan _2 15:04:05 -0700 2006","outputFormat":"2006-01-02T15:04:05-0700}}}]`
	transform, _ := kazaam.NewKazaam(spec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kazaamOut, _ := transform.TransformJSONStringToString(benchmarkJSON)
		nothing(kazaamOut)
	}

}

func BenchmarkIsJsonFast(b *testing.B) {
	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		val := kazaam.IsJsonFast(benchmarkSlice)
		nothing(val)
	}

}
