package netutils

import (
	"net"
	"reflect"
	"testing"
)

type testcase struct {
	Input    []net.IPNet
	Expected []net.IPNet
}

func mustParseCIDR(s string) net.IPNet {
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return *ipnet
}

func NewTestCase(input, expected []string) (tc testcase) {
	tc.Input = make([]net.IPNet, len(input))
	for i, v := range input {
		tc.Input[i] = mustParseCIDR(v)
	}
	tc.Expected = make([]net.IPNet, len(expected))
	for i, v := range expected {
		tc.Expected[i] = mustParseCIDR(v)
	}
	return
}

var tcs = []testcase{
	NewTestCase(nil, nil),
	NewTestCase([]string{}, []string{}),
	NewTestCase(
		[]string{"192.0.2.0/24"},
		[]string{"192.0.2.0/24"},
	),
	NewTestCase(
		[]string{"2001:db8:a0b:12f0::/32"},
		[]string{"2001:db8:a0b:12f0::/32"},
	),
	NewTestCase(
		[]string{"192.0.2.0/26", "192.0.2.64/26"},
		[]string{"192.0.2.0/25"},
	),
	NewTestCase(
		[]string{"192.0.2.0/26", "192.0.2.128/26"},
		[]string{"192.0.2.0/26", "192.0.2.128/26"},
	),
	NewTestCase(
		[]string{"192.0.2.0/26", "192.0.2.64/26", "192.0.2.128/26", "192.0.2.192/26"},
		[]string{"192.0.2.0/24"},
	),
	NewTestCase(
		[]string{"2001:db8::/64", "2001:db8:0:1::/64"},
		[]string{"2001:db8::/63"},
	),
	NewTestCase(
		[]string{"2001:db8::/64", "2001:db8:0:1::/64", "2001:db8:0:2::/64"},
		[]string{"2001:db8::/63", "2001:db8:0:2::/64"},
	),
}

func TestAggregate(t *testing.T) {
	for i, v := range tcs {
		got := Aggregate(v.Input)
		if !reflect.DeepEqual(got, v.Expected) {
			t.Errorf("[%d] error: Got: %v, Expected: %v", i, got, v.Expected)
		}
	}
}

func TestMerge(t *testing.T) {
	type tc struct {
		A      string
		B      string
		C      string
		Merged bool
	}
	tcs := []tc{
		{"192.168.0.0/24", "192.168.1.0/24", "192.168.0.0/23", true},
		{"192.168.0.0/24", "192.168.0.0/24", "192.168.0.0/24", true},
		{"192.168.0.0/24", "192.168.2.0/24", "192.168.0.0/24", false},
		{"192.168.0.0/24", "192.168.1.0/25", "192.168.0.0/24", false},
		{"192.168.2.0/26", "192.168.2.128/26", "192.0.0.0/24", false},
	}
	for i, v := range tcs {
		a := mustParseCIDR(v.A)
		b := mustParseCIDR(v.B)
		c := mustParseCIDR(v.C)
		got, merged := merge(a, b)
		if merged != v.Merged {
			t.Errorf("[%d]unexpected merged. Got: %#v, Expected: %#v", i, merged, v.Merged)
		}
		if merged == false {
			continue
		}
		if !reflect.DeepEqual(got, c) {
			t.Errorf("[%d]ipnet unexpected. Got: %#v, Expected: %#v", i, got, c)
		}
	}
}

func TestMaskup(t *testing.T) {
	tcs := [][2]net.IPMask{
		{{0x00, 0x00, 0x00, 0x00}, {0x00, 0x00, 0x00, 0x00}},
		{{0x80, 0x00, 0x00, 0x00}, {0x00, 0x00, 0x00, 0x00}},
		{{0xc0, 0x00, 0x00, 0x00}, {0x80, 0x00, 0x00, 0x00}},
		{{0xe0, 0x00, 0x00, 0x00}, {0xc0, 0x00, 0x00, 0x00}},
		{{0xf0, 0x00, 0x00, 0x00}, {0xe0, 0x00, 0x00, 0x00}},
		{{0xff, 0xf0, 0x00, 0x00}, {0xff, 0xe0, 0x00, 0x00}},
		{{0xff, 0x00, 0x00, 0x00}, {0xfe, 0x00, 0x00, 0x00}},
		{{0xff, 0xff, 0xff, 0xff}, {0xff, 0xff, 0xff, 0xfe}},
	}
	for i, v := range tcs {
		got := maskup(v[0])
		if !reflect.DeepEqual(got, v[1]) {
			t.Errorf("[%d]error: got %v, expected: %v", i, got, v[1])
		}
	}
}
