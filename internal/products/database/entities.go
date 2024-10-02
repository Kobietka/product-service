package database

type productEntity struct {
	Ean  string
	Name string
}

type packagingEntity struct {
	Ean   string
	Value float32
	Unit  string
}

type nutritionEntity struct {
	Ean  string
	Kcal int32
}

type nutritionQuantityEntity struct {
	Ean   string
	Value float32
	Unit  string
}

type nutrientEntity struct {
	Ean   string
	Type  string
	Value float32
	Unit  string
}

type vitaminEntity struct {
	Ean   string
	Type  string
	Value float32
	Unit  string
}

type mineralEntity struct {
	Ean   string
	Type  string
	Value float32
	Unit  string
}
