package testdata

// User represents a simple user entity with basic fields.
type User struct {
	ID    []int
	Name  []string
	Email []string
}

// Product represents a product with ID, Name, and Price.
type Product struct {
	ID    []int
	Name  []string
	Price []float64
}

// Company represents a company entity containing a list of employees.
type Company struct {
	Name      []string
	Employees []*User
}

// Address represents a physical address.
type Address struct {
	Street  []string
	City    []string
	ZipCode []*string
}

// UserProfile represents a user's profile information, including nested address.
type UserProfile struct {
	Name    []string
	Age     []int
	Address []Address
}

// Employee represents an employee with a profile and skills.
type Employee struct {
	ID      []int
	Profile []*UserProfile
	Skills  []string
}
