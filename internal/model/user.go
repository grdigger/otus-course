package model

type User struct {
	ID       int
	Name     *string
	Surname  *string
	Age      *int
	GenderId *int
	CityId   *int
	Gender   *string
	City     *string
}

func (u *User) GetID() int {
	return u.ID
}

func (u *User) GetName() string {
	if u.Name != nil {
		return *u.Name
	}
	return ""
}

func (u *User) GetSurname() string {
	if u.Surname != nil {
		return *u.Surname
	}
	return ""
}

func (u *User) GetAge() int {
	if u.Age != nil {
		return *u.Age
	}
	return 0
}

func (u *User) GetGenderId() int {
	if u.GenderId != nil {
		return *u.GenderId
	}
	return 0
}

func (u *User) GetCityId() int {
	if u.CityId != nil {
		return *u.CityId
	}
	return 0
}

func (u *User) GetGender() string {
	if u.Gender != nil {
		return *u.Gender
	}
	return ""
}

func (u *User) GetCity() string {
	if u.City != nil {
		return *u.City
	}
	return ""
}

func (u *User) IsEmpty() bool {
	return u.ID == 0
}
