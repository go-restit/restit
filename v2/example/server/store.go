package server

// Storable are information thats storable
type Storable interface {
	GetID() interface{}
	GetType() string
	SetID(string)
}

// NewStore returns a blank Store
func NewStore() *Store {
	return &Store{}
}

// Store is a dummy storage device that
// stores Storable
type Store map[string]map[interface{}]interface{}

// List returns a list of Storabes of the given type
func (s *Store) List(t string) (l []interface{}) {
	l = make([]interface{}, 0)
	if _, ok := (*s)[t]; !ok {
		return
	}
	for _, v := range (*s)[t] {
		l = append(l, v)
	}
	return
}

// Put append / overwrite the given Storable into the store
func (s *Store) Put(v Storable) {
	if _, ok := (*s)[v.GetType()]; !ok {
		(*s)[v.GetType()] = make(map[interface{}]interface{})
	}
	(*s)[v.GetType()][v.GetID()] = v
}

// Get retrieve given Storable from the store
func (s *Store) Get(t string, id interface{}) (ret interface{}) {
	var ok bool
	if _, ok = (*s)[t]; !ok {
		return nil
	} else if ret, ok = (*s)[t][id]; !ok {
		return nil
	}
	return
}

// Delete removes Storable of given type and ID
func (s *Store) Delete(t string, id interface{}) {
	if _, ok := (*s)[t]; !ok {
		return
	} else if _, ok := (*s)[t][id]; !ok {
		return
	}
	delete((*s)[t], id)
}
