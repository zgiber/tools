package tools

import (
	"fmt"
	"testing"
)

var (
	testData1 = []byte(`
	{
		"values":{
			"red": 123,
			"green":234,
			"blue":23
		},
		"pixels": [
			{"x": 23, "y":34},
			{"x": 25, "y":35},
			{"x": 33, "y":64},
			{"x": 73, "y":11},
			{"x": 23, "y":94},
			{"x": 123, "y":14}
		]
	}
	`)

	testData2 = map[string]interface{}{
		"values": map[string]interface{}{
			"red":   123,
			"green": 234,
			"blue":  23,
		},
		"pixels": []map[string]interface{}{
			map[string]interface{}{
				"x": 23,
				"y": 34,
			},
			map[string]interface{}{
				"x": 25,
				"y": 35,
			},
			map[string]interface{}{
				"x": 33,
				"y": 64,
			},
			map[string]interface{}{
				"x": 73,
				"y": 11,
			},
			map[string]interface{}{
				"x": 23,
				"y": 94,
			},
			map[string]interface{}{
				"x": 123,
				"y": 14,
			},
		},
	}
)

func TestFlatten(t *testing.T) {
	result, err := Flatten(testData1, "_", "/")
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range result {
		fmt.Println(k, v)
	}

	expected := map[string]interface{}{
		"_values/red":   123.0,
		"_values/green": 234.0,
		"_values/blue":  23.0,
		"_pixels/0/x":   23.0,
		"_pixels/1/x":   25.0,
		"_pixels/2/x":   33.0,
		"_pixels/3/x":   73.0,
		"_pixels/4/x":   23.0,
		"_pixels/5/x":   123.0,
		"_pixels/0/y":   34.0,
		"_pixels/1/y":   35.0,
		"_pixels/2/y":   64.0,
		"_pixels/3/y":   11.0,
		"_pixels/4/y":   94.0,
		"_pixels/5/y":   14.0,
	}

	for k, v := range result {
		if expected[k] != v {
			fmt.Println(expected[k], v)
			t.Fatalf("Bad result.")
		}
	}
}
