package proto

//easyjson:json
type User struct {
	Id        Int       `json:"id"`
	Email     String100 `json:"email"`
	FirstName String100 `json:"first_name"`
	LastName  String100 `json:"last_name"`
	Gender    String100 `json:"gender"`
	BirthDate Int64     `json:"birth_date"`
}

//easyjson:json
type Location struct {
	Id       Int       `json:"id"`
	Place    String100 `json:"place"`
	Country  String100 `json:"country"`
	City     String100 `json:"city"`
	Distance Int       `json:"distance"`
}

//easyjson:json
type Visit struct {
	Id        Int   `json:"id"`
	Location  Int   `json:"location"`
	User      Int   `json:"user"`
	VisitedAt Int64 `json:"visited_at"`
	Mark      Int   `json:"mark"`
}
