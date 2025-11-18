package testdata

type User struct {
	ID    []int
	Name  []string
	Email []string
}

type Product struct {
	ID    []int
	Name  []string
	Price []float64
}

type Company struct {
	Name      []string
	Employees []*User
}

type Address struct {
	Street  []string
	City    []string
	ZipCode []*string
}

type UserProfile struct {
	Name    []string
	Age     []int
	Address []Address
}

type Employee struct {
	ID      []int
	Profile []*UserProfile
	Skills  []string
}
