package config_test

import (
	"testing"

	"example.com/layered-orders/internal/config"
)

func TestHTTPAddr(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value string
		want  string
	}{
		{name: "default", value: "", want: ":8080"},
		{name: "override", value: "127.0.0.1:9090", want: "127.0.0.1:9090"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("HTTP_ADDR", tc.value)
			if got := config.HTTPAddr(); got != tc.want {
				t.Fatalf("HTTPAddr() = %q, ожидалось %q", got, tc.want)
			}
		})
	}
}
