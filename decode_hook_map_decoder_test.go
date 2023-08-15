package mapstructurex

import (
	"errors"
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

type mapDecoderConfig struct {
	mapDecoderConfigType
}

type mapDecoderConfigRaw struct {
	Type   string
	Config map[string]any
}

func (c *mapDecoderConfig) DecodeMap(v map[string]any) error {
	var rawConfig mapDecoderConfigRaw

	err := mapstructure.Decode(v, &rawConfig)
	if err != nil {
		return err
	}

	switch rawConfig.Type {
	case "bar":
		var bar mapDecoderConfigTypeBar
		err := mapstructure.Decode(rawConfig.Config, &bar)
		if err != nil {
			return err
		}
		c.mapDecoderConfigType = bar
	case "baz":
		var baz mapDecoderConfigTypeBaz
		err := mapstructure.Decode(rawConfig.Config, &baz)
		if err != nil {
			return err
		}
		c.mapDecoderConfigType = baz
	default:
		return errors.New("unknown type")

	}

	return nil
}

type mapDecoderConfigType interface {
	Foo() string
}

type mapDecoderConfigTypeBar struct {
	Bar string
}

func (b mapDecoderConfigTypeBar) Foo() string {
	return b.Bar
}

type mapDecoderConfigTypeBaz struct {
	Baz string
}

func (b mapDecoderConfigTypeBaz) Foo() string {
	return b.Baz
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
		{
			reflect.ValueOf(map[string]any{
				"type": "baz",
				"config": map[string]any{
					"baz": "baz",
				},
			}),
			reflect.ValueOf(mapDecoderConfig{}),
			&mapDecoderConfig{
				mapDecoderConfigType: mapDecoderConfigTypeBaz{
					Baz: "baz",
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

type mapDecoderInfiniteLoop struct {
	Key string
}

func (d *mapDecoderInfiniteLoop) DecodeMap(v map[string]any) error {
	type alias mapDecoderInfiniteLoop

	var result alias

	decoder, err := createDecoderWithMapDecoderHook(&result)
	if err != nil {
		return err
	}

	err = decoder.Decode(v)
	if err != nil {
		return err
	}

	*d = mapDecoderInfiniteLoop(result)

	return nil
}

func createDecoderWithMapDecoderHook(result any) (*mapstructure.Decoder, error) {
	return mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: MapDecoderHookFunc(),
		Metadata:   nil,
		Result:     result,
	})
}

func TestMapDecoderHookFunc_DecodeInfiniteLoop(t *testing.T) {
	var result mapDecoderInfiniteLoop

	decoder, err := createDecoderWithMapDecoderHook(&result)
	if err != nil {
		t.Fatal(err)
	}

	input := map[string]any{
		"key": "value",
	}

	err = decoder.Decode(input)
	if err != nil {
		t.Fatal(err)
	}

	expected := mapDecoderInfiniteLoop{
		Key: "value",
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf(
			"case map decoder infinite loop: expected %#v, got %#v",
			expected, result,
		)
	}
}
