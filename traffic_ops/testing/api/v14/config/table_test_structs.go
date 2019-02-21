package config

type NegativeTest struct {
	Description string
	Config      string
	Expected    int
}

type PositiveTest struct {
	Description string
	Config      string
}
