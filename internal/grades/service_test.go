package grades_test

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"traininggrades/internal/grades"
)

func TestServiceLifecycle(t *testing.T) {
	var submitted []string
	service := grades.NewService(
		grades.NewMemoryStore(),
		grades.WithAfterSubmit(func(result grades.SubmissionResult) {
			submitted = append(submitted, result.Student.ID)
		}),
	)

	first, err := service.Submit(grades.Student{ID: "T002", Name: "李强", Homework1: 76, Homework2: 84, FinalExam: 80})
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if first.Average != "80.00" {
		t.Fatalf("Submit() average = %s, want 80.00", first.Average)
	}
	_, err = service.Submit(grades.Student{ID: "T001", Name: "张敏", Homework1: 88, Homework2: 91, FinalExam: 94})
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if !reflect.DeepEqual(submitted, []string{"T002", "T001"}) {
		t.Fatalf("submitted IDs = %v", submitted)
	}

	student, err := service.Query("T002")
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if student.Name != "李强" {
		t.Fatalf("Query() name = %s", student.Name)
	}

	updated, err := service.Modify("T002", grades.Student{Name: "李强", Homework1: 90, Homework2: 92, FinalExam: 94})
	if err != nil {
		t.Fatalf("Modify() error = %v", err)
	}
	if updated.Average().StringFixed(2) != "92.00" {
		t.Fatalf("Modify() average = %s", updated.Average().StringFixed(2))
	}

	sorted, err := service.Sort(grades.SortByAverage, true)
	if err != nil {
		t.Fatalf("Sort() error = %v", err)
	}
	if got := []string{sorted[0].ID, sorted[1].ID}; !reflect.DeepEqual(got, []string{"T002", "T001"}) {
		t.Fatalf("Sort() IDs = %v", got)
	}

	if err := service.Delete("T001"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := service.Query("T001"); !errors.Is(err, grades.ErrStudentNotFound) {
		t.Fatalf("Query() error = %v", err)
	}
}

func TestServiceRejectsInvalidAndDuplicateRecords(t *testing.T) {
	service := grades.NewService(grades.NewMemoryStore(), grades.WithAfterSubmit(func(grades.SubmissionResult) {}))
	invalid := grades.Student{ID: "T001", Name: "周宁", Homework1: 101, Homework2: 80, FinalExam: 90}
	if _, err := service.Submit(invalid); err == nil {
		t.Fatal("Submit() error = nil")
	}
	valid := grades.Student{ID: "T001", Name: "周宁", Homework1: 70, Homework2: 80, FinalExam: 90}
	if _, err := service.Submit(valid); err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if _, err := service.Submit(valid); !errors.Is(err, grades.ErrDuplicateID) {
		t.Fatalf("Submit() error = %v", err)
	}
}

func TestServiceSavesAndLoadsDeterministically(t *testing.T) {
	service := grades.NewService(grades.NewMemoryStore(), grades.WithAfterSubmit(func(grades.SubmissionResult) {}))
	students := []grades.Student{
		{ID: "T020", Name: "陈曦", Homework1: 72, Homework2: 84, FinalExam: 90},
		{ID: "T010", Name: "赵峰", Homework1: 96, Homework2: 88, FinalExam: 91},
	}
	for _, student := range students {
		if _, err := service.Submit(student); err != nil {
			t.Fatalf("Submit() error = %v", err)
		}
	}

	path := filepath.Join(t.TempDir(), "grades.json")
	if err := service.Save(path); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	loaded := grades.NewService(grades.NewMemoryStore())
	if err := loaded.Load(path); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	want := []grades.Student{students[1], students[0]}
	if got := loaded.List(); !reflect.DeepEqual(got, want) {
		t.Fatalf("List() = %#v, want %#v", got, want)
	}
}

func TestServiceLoadsRepositoryFixture(t *testing.T) {
	service := grades.NewService(grades.NewMemoryStore())
	if err := service.Load(filepath.Join("..", "..", "fixtures", "students.json")); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	students := service.List()
	if len(students) != 2 || students[0].ID != "T001" || students[1].ID != "T002" {
		t.Fatalf("List() = %#v", students)
	}
}
