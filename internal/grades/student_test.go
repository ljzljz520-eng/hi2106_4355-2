package grades_test

import (
	"testing"

	"traininggrades/internal/grades"
)

func TestStudentValidation(t *testing.T) {
	tests := []struct {
		name    string
		student grades.Student
		valid   bool
	}{
		{
			name:    "valid boundary scores",
			student: grades.Student{ID: "T001", Name: "周宁", Homework1: 0, Homework2: 100, FinalExam: 60},
			valid:   true,
		},
		{
			name:    "missing identifier",
			student: grades.Student{Name: "周宁", Homework1: 70, Homework2: 80, FinalExam: 90},
		},
		{
			name:    "score below range",
			student: grades.Student{ID: "T001", Name: "周宁", Homework1: -1, Homework2: 80, FinalExam: 90},
		},
		{
			name:    "score above range",
			student: grades.Student{ID: "T001", Name: "周宁", Homework1: 70, Homework2: 101, FinalExam: 90},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.student.Validate()
			if test.valid && err != nil {
				t.Fatalf("Validate() error = %v", err)
			}
			if !test.valid && err == nil {
				t.Fatal("Validate() error = nil")
			}
		})
	}
}

func TestStudentAverageUsesExactDecimalResult(t *testing.T) {
	student := grades.Student{ID: "T001", Name: "周宁", Homework1: 80, Homework2: 81, FinalExam: 81}
	if got := student.Average().StringFixed(2); got != "80.67" {
		t.Fatalf("Average() = %s, want 80.67", got)
	}
}
