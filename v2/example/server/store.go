package server

// Storable are information thats storable
type Storable interface {
	GetID() interface{}
	GetType() string
	SetID(string)
}

// NewStore returns a blank Store
func NewStore() *Store {
	return &Store{
		orders: make(map[string][]interface{}),
		values: make(map[string]map[interface{}]interface{}),
	}
}

// Store is a dummy storage device that
// stores Storable
type Store struct {
	orders map[string][]interface{}
	values map[string]map[interface{}]interface{}
}

// List returns a list of Storabes of the given type
func (s *Store) List(t string) (l []interface{}) {
	l = make([]interface{}, 0)
	if _, ok := (*s).values[t]; !ok {
		return
	}
	for _, k := range (*s).orders[t] {
		if _, ok := (*s).values[t][k]; ok {
			l = append(l, (*s).values[t][k])
		}
	}
	return
}

// Put append / overwrite the given Storable into the store
func (s *Store) Put(v Storable) {

	// if there is not type, create type
	if _, ok := (*s).values[v.GetType()]; !ok {
		(*s).orders[v.GetType()] = make([]interface{}, 0, 1)
		(*s).values[v.GetType()] = make(map[interface{}]interface{})
	}

	// if there is no old value, add to order
	if _, ok := (*s).values[v.GetType()][v.GetID()]; !ok {
		(*s).orders[v.GetType()] = append((*s).orders[v.GetType()], v.GetID())
	}

	// assign value
	(*s).values[v.GetType()][v.GetID()] = v
}

// Get retrieve given Storable from the store
func (s *Store) Get(t string, id interface{}) (ret interface{}) {
	var ok bool
	if _, ok = (*s).values[t]; !ok {
		return nil
	} else if ret, ok = (*s).values[t][id]; !ok {
		return nil
	}
	return
}

// Delete removes Storable of given type and ID
func (s *Store) Delete(t string, id interface{}) {
	if _, ok := (*s).values[t]; !ok {
		return
	} else if _, ok := (*s).values[t][id]; !ok {
		return
	}
	delete((*s).values[t], id)
	// TODO: remove order
}
