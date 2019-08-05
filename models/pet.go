package models

type Pet struct {
	ID      int     `db:"user_id"`
	Name    string  `db:"name"`
	PetType string  `db:"animal_type"`
	Breed   string  `db:"breed"`
	Age     int     `db:"age"`
	Weight  float32 `db: "weight"`
	Gender  string  `db:"gender"`
}
