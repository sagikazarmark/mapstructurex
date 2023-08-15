package mapstructurex_test

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/sagikazarmark/mapstructurex"
)

type Person struct {
	Name string
	Age  int
}

func (p *Person) DecodeMap(v map[string]any) error {
	type alias Person

	var result alias

	decoder, err := CreateDecoderWithMapDecoderHook(&result)
	if err != nil {
		return err
	}

	err = decoder.Decode(v)
	if err != nil {
		return err
	}

	*p = Person(result)

	return nil
}

func (p Person) String() string {
	return fmt.Sprintf("%s: %d", p.Name, p.Age)
}

func CreateDecoderWithMapDecoderHook(result any) (*mapstructure.Decoder, error) {
	return mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructurex.MapDecoderHookFunc(),
		Metadata:   nil,
		Result:     result,
	})
}

func ExampleMapDecoderHookFunc() {
	var person Person

	decoder, _ := CreateDecoderWithMapDecoderHook(&person)

	input := map[string]any{
		"name": "Bob",
		"age":  42,
	}

	_ = decoder.Decode(input)

	fmt.Println(person)
	// Output: Bob: 42
}

// Config is a struct that implements a custom decoding function.
// The embedded [Fooer] interface provides polymorphism.
// Alternative implementations of [Fooer] can provide different configurations.
type Config struct {
	Fooer
}

// Fooer is an interface that provides polymorphism for [Config].
type Fooer interface {
	Foo() string
}

// RawConfig mimics the structure of configuration (for example in a config file).
// Type determines the [Fooer] implementation, Config is the input for the implementation.
type RawConfig struct {
	Type   string
	Config map[string]any
}

func (c *Config) DecodeMap(v map[string]any) error {
	var rawConfig RawConfig

	err := mapstructure.Decode(v, &rawConfig)
	if err != nil {
		return err
	}

	switch rawConfig.Type {
	case "bar":
		var bar ConfigBar
		err := mapstructure.Decode(rawConfig.Config, &bar)
		if err != nil {
			return err
		}
		c.Fooer = bar
	case "baz":
		var baz ConfigBaz
		err := mapstructure.Decode(rawConfig.Config, &baz)
		if err != nil {
			return err
		}
		c.Fooer = baz
	default:
		return errors.New("unknown type")

	}

	return nil
}

// ConfigBar is a [Fooer].
type ConfigBar struct {
	Bar string
}

func (c ConfigBar) Foo() string {
	return c.Bar
}

// ConfigBaz is a [Fooer].
type ConfigBaz struct {
	Baz string
}

func (c ConfigBaz) Foo() string {
	return c.Baz
}

func ExampleMapDecoderHookFunc_polymorphic() {
	var config Config

	decoder, _ := CreateDecoderWithMapDecoderHook(&config)

	input := map[string]any{
		"type": "baz",
		"config": map[string]any{
			"baz": "bat",
		},
	}

	err := decoder.Decode(input)
	if err != nil {
		panic(err)
	}

	fmt.Println(config.Foo())
	// Output: bat
}
