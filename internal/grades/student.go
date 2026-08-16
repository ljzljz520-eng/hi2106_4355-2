package grades

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	ErrStudentNotFound = errors.New("student not found")
	ErrDuplicateID     = errors.New("student ID already exists")
)

type Student struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Homework1 int    `json:"homework1"`
	Homework2 int    `json:"homework2"`
	FinalExam int    `json:"final_exam"`
}

func (s Student) Validate() error {
	if strings.TrimSpace(s.ID) == "" {
		return errors.New("student ID is required")
	}
	if strings.TrimSpace(s.Name) == "" {
		return errors.New("student name is required")
	}
	for label, score := range map[string]int{
		"homework one": s.Homework1,
		"homework two": s.Homework2,
		"final exam":   s.FinalExam,
	} {
		if score < 0 || score > 100 {
			return fmt.Errorf("%s score must be between 0 and 100", label)
		}
	}
	return nil
}

func (s Student) Average() decimal.Decimal {
	total := int64(s.Homework1 + s.Homework2 + s.FinalExam)
	return decimal.NewFromInt(total).Div(decimal.NewFromInt(3)).Round(2)
}

func normalizeStudent(student Student) Student {
	student.ID = strings.TrimSpace(student.ID)
	student.Name = strings.TrimSpace(student.Name)
	return student
}
