package logic

import "sync"

type VarStore struct {
	data map[string]int
	mu   sync.RWMutex
}

func NewVarStore() *VarStore {
	return &VarStore{
		data: make(map[string]int, 20),
	}
}

func (s *VarStore) Set(name string, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[name] = value
}

func (s *VarStore) Get(name string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[name]
	return val, ok
}
