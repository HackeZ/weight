package weight

// Weight interface of implement round-robin weight balancing algorithm
type Weight interface {

	// Add a naming element with weight
	Add(name string, elem interface{}, weight int) error

	// Next return existed element by balancing rule
	Next() interface{}

	// Remove an element by its name, and return error if
	// this element not exist.
	Remove(name string) error

	// Update current weight for element by name, and
	// return error if this element not exist
	Update(name string, weight int) error

	// Total return total weight of this Weight
	Total() int

	// Close weight balancing using close function
	Close(fn func(interface{}) error) error
}
