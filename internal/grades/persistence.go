package grades

import (
	"encoding/json"
	"fmt"
	"os"
)

func saveStudents(path string, students []Student) error {
	data, err := json.MarshalIndent(students, "", "  ")
	if err != nil {
		return fmt.Errorf("encode grades: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write grades: %w", err)
	}
	return nil
}

func loadStudents(path string) ([]Student, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read grades: %w", err)
	}
	var students []Student
	if err := json.Unmarshal(data, &students); err != nil {
		return nil, fmt.Errorf("decode grades: %w", err)
	}
	return students, nil
}
