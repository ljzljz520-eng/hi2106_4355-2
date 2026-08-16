package cli_test

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"traininggrades/internal/cli"
	"traininggrades/internal/grades"
)

func TestRunnerShowsFixtureRecords(t *testing.T) {
	service := grades.NewService(grades.NewMemoryStore())
	if err := service.Load(filepath.Join("..", "..", "fixtures", "students.json")); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	var output bytes.Buffer
	runner := cli.NewRunner(strings.NewReader("2\n0\n"), &output, service, "grades.json")
	if err := runner.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	for _, expected := range []string{"T001", "张敏", "91.00", "T002", "李强", "80.00"} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("output does not contain %q: %s", expected, output.String())
		}
	}
}

func TestRunnerValidatesScoreInput(t *testing.T) {
	service := grades.NewService(grades.NewMemoryStore(), grades.WithAfterSubmit(func(grades.SubmissionResult) {}))
	input := strings.NewReader("1\nT003\n王芳\n-1\n88\n101\n89\n90\n0\n")
	var output bytes.Buffer
	runner := cli.NewRunner(input, &output, service, "grades.json")
	if err := runner.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if count := strings.Count(output.String(), "成绩必须是 0 到 100 的整数"); count != 2 {
		t.Fatalf("validation message count = %d, want 2", count)
	}
	if !strings.Contains(output.String(), "录入成功: T003 王芳 平均分 89.00") {
		t.Fatalf("output = %s", output.String())
	}
}

func TestRunnerRecordsStudent(t *testing.T) {
	service := grades.NewService(grades.NewMemoryStore())
	input := strings.NewReader("1\nT003\n王芳\n88\n89\n90\n0\n")
	var output bytes.Buffer
	runner := cli.NewRunner(input, &output, service, "grades.json")
	if err := runner.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !strings.Contains(output.String(), "录入成功: T003 王芳 平均分 89.00") {
		t.Fatalf("output = %s", output.String())
	}
}
