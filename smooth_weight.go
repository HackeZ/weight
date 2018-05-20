package weight

import (
	"fmt"
	"sync"
)

var _ Weight = &SW{}

// SW smooth weighted round-robin implement
type SW struct {
	mtx   sync.RWMutex
	elems []*smoothElem
}

// smoothElem element of smooth weight round-robin implement
type smoothElem struct {
	weight          int
	currentWeight   int
	effectiveWeight int

	name string
	elem interface{}
}

func (se *smoothElem) String() string {
	return fmt.Sprintf("%+v", *se)
}

// Add a naming element with weight
func (sw *SW) Add(name string, elem interface{}, weight int) error {
	sw.mtx.Lock()

	if sw.elems == nil {
		sw.elems = make([]*smoothElem, 0, 4)
	}

	for _, elem := range sw.elems {
		if elem.name == name {
			sw.mtx.Unlock()
			return fmt.Errorf("element of the same name already exist")
		}
	}

	sw.elems = append(sw.elems, &smoothElem{
		weight:          weight,
		currentWeight:   0,
		effectiveWeight: weight,
		name:            name,
		elem:            elem,
	})

	sw.mtx.Unlock()
	return nil
}

// Remove a single element by its name and return
// error if empty elements or element not exist
func (sw *SW) Remove(name string) error {
	sw.mtx.Lock()

	if sw.elems == nil || len(sw.elems) == 0 {
		sw.mtx.Unlock()
		return fmt.Errorf("smooth weight list is empty")
	}

	for idx, elem := range sw.elems {
		if elem.name == name {
			sw.elems = append(sw.elems[0:idx-1], sw.elems[idx+1:]...)

			sw.mtx.Unlock()
			return nil
		}
	}

	sw.mtx.Unlock()
	return fmt.Errorf("element not exist")
}

// Update weight of the element dynamically
func (sw *SW) Update(name string, weight int) error {
	sw.mtx.Lock()

	if sw.elems == nil || len(sw.elems) == 0 {
		sw.mtx.Unlock()
		return fmt.Errorf("smooth weight list is empty")
	}

	for _, elem := range sw.elems {
		if elem.name == name {
			elem.weight = weight
			elem.effectiveWeight = weight
			elem.currentWeight = 0

			sw.mtx.Unlock()
			return nil
		}
	}

	sw.mtx.Unlock()
	return fmt.Errorf("element not exist")
}

// Total return total of weight in this smooth round-robin weighted
func (sw *SW) Total() int {
	sw.mtx.Lock()

	if sw.elems == nil || len(sw.elems) == 0 {
		sw.mtx.Unlock()
		return 0
	}

	var total int
	for _, elem := range sw.elems {
		total += elem.weight
	}

	sw.mtx.Unlock()
	return total
}

// Close all element using fn
func (sw *SW) Close(fn func(interface{}) error) error {
	sw.mtx.Lock()
	defer sw.mtx.Unlock()

	if fn == nil {
		return nil
	}

	for _, elem := range sw.elems {
		err := fn(elem)
		if err != nil {
			return err
		}
	}

	return nil
}

// Next pick up next element under smooth round-robin weight balancing
func (sw *SW) Next() interface{} {
	sw.mtx.RLock()
	defer sw.mtx.RUnlock()

	if sw.elems == nil || len(sw.elems) == 0 {
		return nil
	}

	var total int
	next := &smoothElem{}
	for _, elem := range sw.elems {
		total += elem.effectiveWeight
		elem.currentWeight += elem.effectiveWeight

		if elem.effectiveWeight < elem.weight { // automatic recovery
			elem.effectiveWeight++
		}

		if next == nil || next.currentWeight < elem.currentWeight {
			next = elem
		}
	}

	next.currentWeight -= total
	return next.elem
}
