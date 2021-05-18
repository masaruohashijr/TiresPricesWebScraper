package models

type Taxonomy struct {
	widths []Width
}

type Width struct {
	data     string
	profiles []Profile
}

type Profile struct {
	data      string
	diameters []Diameter
}

type Diameter struct {
	data  string
	tires []Tire
}

type Tire struct {
	Name        string
	Width       string
	Profile     string
	Diameter    string
	Model       string
	Price       string
	Rebate      string
	Xl          string
	SpeedRating string
}

type Config struct {
	TiresPricesURL string
}
