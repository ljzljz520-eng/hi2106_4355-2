package grades

import (
	"fmt"
	"sort"
	"strings"
)

type SortField string

const (
	SortByID        SortField = "id"
	SortByName      SortField = "name"
	SortByHomework1 SortField = "homework1"
	SortByHomework2 SortField = "homework2"
	SortByFinalExam SortField = "final"
	SortByAverage   SortField = "average"
)

type ServiceOption func(*Service)

func WithAfterSubmit(callback SubmitCallback) ServiceOption {
	return func(service *Service) {
		service.executor = NewExtensionExecutor(Extensions{AfterSubmit: callback})
	}
}

type Service struct {
	store    Store
	executor ExtensionExecutor
}

func NewService(store Store, options ...ServiceOption) *Service {
	if store == nil {
		store = NewMemoryStore()
	}
	service := &Service{
		store:    store,
		executor: NewExtensionExecutor(Extensions{}),
	}
	for _, option := range options {
		option(service)
	}
	return service
}

func (s *Service) Submit(student Student) (SubmissionResult, error) {
	student = normalizeStudent(student)
	if err := student.Validate(); err != nil {
		return SubmissionResult{}, err
	}
	if err := s.store.Create(student); err != nil {
		return SubmissionResult{}, err
	}
	result := SubmissionResult{
		Student: student,
		Average: student.Average().StringFixed(2),
	}
	s.executor.ExecuteAfterSubmit(result)
	return result, nil
}

func (s *Service) List() []Student {
	return s.store.List()
}

func (s *Service) Query(id string) (Student, error) {
	return s.store.Get(strings.TrimSpace(id))
}

func (s *Service) Sort(field SortField, descending bool) ([]Student, error) {
	students := s.store.List()
	if !validSortField(field) {
		return nil, fmt.Errorf("unsupported sort field %q", field)
	}
	sort.SliceStable(students, func(i, j int) bool {
		comparison := compareStudents(students[i], students[j], field)
		if comparison == 0 {
			comparison = strings.Compare(students[i].ID, students[j].ID)
		}
		if descending {
			return comparison > 0
		}
		return comparison < 0
	})
	return students, nil
}

func (s *Service) Delete(id string) error {
	return s.store.Delete(strings.TrimSpace(id))
}

func (s *Service) Modify(id string, student Student) (Student, error) {
	student.ID = strings.TrimSpace(id)
	student = normalizeStudent(student)
	if err := student.Validate(); err != nil {
		return Student{}, err
	}
	if err := s.store.Update(student); err != nil {
		return Student{}, err
	}
	return student, nil
}

func (s *Service) Save(path string) error {
	return saveStudents(path, s.store.List())
}

func (s *Service) Load(path string) error {
	students, err := loadStudents(path)
	if err != nil {
		return err
	}
	return s.store.ReplaceAll(students)
}

func validSortField(field SortField) bool {
	switch field {
	case SortByID, SortByName, SortByHomework1, SortByHomework2, SortByFinalExam, SortByAverage:
		return true
	default:
		return false
	}
}

func compareStudents(left, right Student, field SortField) int {
	switch field {
	case SortByID:
		return strings.Compare(left.ID, right.ID)
	case SortByName:
		return strings.Compare(left.Name, right.Name)
	case SortByHomework1:
		return compareInt(left.Homework1, right.Homework1)
	case SortByHomework2:
		return compareInt(left.Homework2, right.Homework2)
	case SortByFinalExam:
		return compareInt(left.FinalExam, right.FinalExam)
	case SortByAverage:
		return left.Average().Cmp(right.Average())
	default:
		return 0
	}
}

func compareInt(left, right int) int {
	if left < right {
		return -1
	}
	if left > right {
		return 1
	}
	return 0
}
