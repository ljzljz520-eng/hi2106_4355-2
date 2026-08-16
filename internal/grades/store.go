package grades

import (
	"sort"
	"sync"
)

type Store interface {
	Create(Student) error
	Get(string) (Student, error)
	List() []Student
	Update(Student) error
	Delete(string) error
	ReplaceAll([]Student) error
}

type MemoryStore struct {
	mu       sync.RWMutex
	students map[string]Student
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{students: make(map[string]Student)}
}

func (s *MemoryStore) Create(student Student) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.students[student.ID]; exists {
		return ErrDuplicateID
	}
	s.students[student.ID] = student
	return nil
}

func (s *MemoryStore) Get(id string) (Student, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	student, exists := s.students[id]
	if !exists {
		return Student{}, ErrStudentNotFound
	}
	return student, nil
}

func (s *MemoryStore) List() []Student {
	s.mu.RLock()
	defer s.mu.RUnlock()
	students := make([]Student, 0, len(s.students))
	for _, student := range s.students {
		students = append(students, student)
	}
	sort.Slice(students, func(i, j int) bool {
		return students[i].ID < students[j].ID
	})
	return students
}

func (s *MemoryStore) Update(student Student) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.students[student.ID]; !exists {
		return ErrStudentNotFound
	}
	s.students[student.ID] = student
	return nil
}

func (s *MemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.students[id]; !exists {
		return ErrStudentNotFound
	}
	delete(s.students, id)
	return nil
}

func (s *MemoryStore) ReplaceAll(students []Student) error {
	replacement := make(map[string]Student, len(students))
	for _, student := range students {
		student = normalizeStudent(student)
		if err := student.Validate(); err != nil {
			return err
		}
		if _, exists := replacement[student.ID]; exists {
			return ErrDuplicateID
		}
		replacement[student.ID] = student
	}
	s.mu.Lock()
	s.students = replacement
	s.mu.Unlock()
	return nil
}
