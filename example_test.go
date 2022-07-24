package gojson

type CommonStruct struct {
	ExampleInt     int     `json:"example_int"`
	ExampleFloat64 float64 `json:"example_float_64"`
	ExampleString  string  `json:"example_string"`
}

type NestedStruct struct {
	Nested         *CommonStruct `json:"nested"`
	ExampleInt     int           `json:"example_int"`
	ExampleFloat64 float64       `json:"example_float_64"`
	ExampleString  string        `json:"example_string"`
}
