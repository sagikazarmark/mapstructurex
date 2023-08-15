package mapstructurex

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// MapDecoder is the interface implemented by an object that can decode a raw map representation of itself.
type MapDecoder interface {
	DecodeMap(map[string]any) error
}

// MapDecoderHookFunc returns a {mapstructure.DecodeHookFunc} that applies maps to the DecodeMap function,
// when the target type implements the {MapDecoder} interface.
func MapDecoderHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.Map || f.Key().Kind() != reflect.String || f.Elem().Kind() != reflect.Interface {
			return data, nil
		}

		result := reflect.New(t).Interface()
		mapDecoder, ok := result.(MapDecoder)
		if !ok {
			return data, nil
		}

		if err := mapDecoder.DecodeMap(data.(map[string]any)); err != nil {
			return nil, err
		}

		return result, nil
	}
}
