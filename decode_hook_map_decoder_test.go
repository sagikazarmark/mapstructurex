package mapstructurex

import (
	"reflect"
	"testing"

	"github.com/mitchellh/mapstructure"
)

type mapDecoder struct {
	Key1 string
	Sub  mapDecoderSub
}

func (d *mapDecoder) DecodeMap(v map[string]any) error {
	return mapstructure.Decode(v, d)
}

type mapDecoderSub struct {
	Key11 string
}

func TestMapDecoderHookFunc(t *testing.T) {
	cases := []struct {
		f, t   reflect.Value
		result interface{}
		err    bool
	}{
		{
			reflect.ValueOf(map[string]any{
				"key1": "value",
				"sub": map[string]any{
					"key11": "value",
				},
			}),
			reflect.ValueOf(mapDecoder{}),
			&mapDecoder{
				Key1: "value",
				Sub: mapDecoderSub{
					Key11: "value",
				},
			},
			false,
		},
	}

	for i, tc := range cases {
		f := MapDecoderHookFunc()

		actual, err := mapstructure.DecodeHookExec(f, tc.f, tc.t)
		if tc.err != (err != nil) {
			t.Fatalf("case %d: expected err %#v", i, tc.err)
		}

		if !reflect.DeepEqual(actual, tc.result) {
			t.Fatalf(
				"case %d: expected %#v, got %#v",
				i, tc.result, actual,
			)
		}
	}
}
