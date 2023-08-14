//go:build go1.18
// +build go1.18

package mapstructurex

import (
	"net/netip"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// StringToNetIPAddrHookFunc returns a DecodeHookFunc that converts
// strings to netip.Addr.
//
// Will be removed once https://github.com/mitchellh/mapstructure/pull/315 is merged.
func StringToNetIPAddrHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(netip.Addr{}) {
			return data, nil
		}

		// Convert it by parsing
		return netip.ParseAddr(data.(string))
	}
}

// StringToNetIPAddrPortHookFunc returns a DecodeHookFunc that converts
// strings to netip.AddrPort.
//
// Will be removed once https://github.com/mitchellh/mapstructure/pull/315 is merged.
func StringToNetIPAddrPortHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(netip.AddrPort{}) {
			return data, nil
		}

		// Convert it by parsing
		return netip.ParseAddrPort(data.(string))
	}
}
