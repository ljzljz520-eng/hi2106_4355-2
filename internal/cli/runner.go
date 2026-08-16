package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"traininggrades/internal/grades"
)

type Runner struct {
	service  *grades.Service
	scanner  *bufio.Scanner
	out      io.Writer
	dataPath string
}

func NewRunner(input io.Reader, output io.Writer, service *grades.Service, dataPath string) *Runner {
	return &Runner{
		service:  service,
		scanner:  bufio.NewScanner(input),
		out:      output,
		dataPath: dataPath,
	}
}

func (r *Runner) Run() error {
	for {
		r.printMenu()
		choice, err := r.readLine("请选择操作: ")
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if choice == "0" {
			fmt.Fprintln(r.out, "已退出")
			return nil
		}
		if err := r.execute(choice); err != nil {
			fmt.Fprintf(r.out, "操作失败: %v\n", err)
		}
	}
}

func (r *Runner) execute(choice string) error {
	switch choice {
	case "1":
		return r.submit()
	case "2":
		r.printStudents(r.service.List())
		return nil
	case "3":
		return r.query()
	case "4":
		return r.sort()
	case "5":
		return r.delete()
	case "6":
		return r.modify()
	case "7":
		if err := r.service.Save(r.dataPath); err != nil {
			return err
		}
		fmt.Fprintf(r.out, "已保存到 %s\n", r.dataPath)
		return nil
	case "8":
		if err := r.service.Load(r.dataPath); err != nil {
			return err
		}
		fmt.Fprintf(r.out, "已从 %s 读取\n", r.dataPath)
		return nil
	default:
		return errors.New("无效的菜单选项")
	}
}

func (r *Runner) submit() error {
	student, err := r.readStudent("")
	if err != nil {
		return err
	}
	result, err := r.service.Submit(student)
	if err != nil {
		return err
	}
	fmt.Fprintf(r.out, "录入成功: %s %s 平均分 %s\n", result.Student.ID, result.Student.Name, result.Average)
	return nil
}

func (r *Runner) query() error {
	id, err := r.readLine("学员编号: ")
	if err != nil {
		return err
	}
	student, err := r.service.Query(id)
	if err != nil {
		return err
	}
	r.printStudents([]grades.Student{student})
	return nil
}

func (r *Runner) sort() error {
	field, err := r.readLine("排序字段(id/name/homework1/homework2/final/average): ")
	if err != nil {
		return err
	}
	direction, err := r.readLine("降序排列(y/N): ")
	if err != nil {
		return err
	}
	students, err := r.service.Sort(grades.SortField(field), strings.EqualFold(direction, "y"))
	if err != nil {
		return err
	}
	r.printStudents(students)
	return nil
}

func (r *Runner) delete() error {
	id, err := r.readLine("要删除的学员编号: ")
	if err != nil {
		return err
	}
	if err := r.service.Delete(id); err != nil {
		return err
	}
	fmt.Fprintf(r.out, "已删除学员 %s\n", strings.TrimSpace(id))
	return nil
}

func (r *Runner) modify() error {
	id, err := r.readLine("要修改的学员编号: ")
	if err != nil {
		return err
	}
	student, err := r.readStudent(id)
	if err != nil {
		return err
	}
	updated, err := r.service.Modify(id, student)
	if err != nil {
		return err
	}
	fmt.Fprintf(r.out, "已修改学员 %s\n", updated.ID)
	return nil
}

func (r *Runner) readStudent(existingID string) (grades.Student, error) {
	id := existingID
	var err error
	if id == "" {
		id, err = r.readLine("学员编号: ")
		if err != nil {
			return grades.Student{}, err
		}
	}
	name, err := r.readLine("姓名: ")
	if err != nil {
		return grades.Student{}, err
	}
	homework1, err := r.readScore("作业一成绩: ")
	if err != nil {
		return grades.Student{}, err
	}
	homework2, err := r.readScore("作业二成绩: ")
	if err != nil {
		return grades.Student{}, err
	}
	finalExam, err := r.readScore("结课测试成绩: ")
	if err != nil {
		return grades.Student{}, err
	}
	return grades.Student{
		ID:        id,
		Name:      name,
		Homework1: homework1,
		Homework2: homework2,
		FinalExam: finalExam,
	}, nil
}

func (r *Runner) readScore(prompt string) (int, error) {
	for {
		value, err := r.readLine(prompt)
		if err != nil {
			return 0, err
		}
		score, err := strconv.Atoi(value)
		if err == nil && score >= 0 && score <= 100 {
			return score, nil
		}
		fmt.Fprintln(r.out, "成绩必须是 0 到 100 的整数")
	}
}

func (r *Runner) readLine(prompt string) (string, error) {
	fmt.Fprint(r.out, prompt)
	if !r.scanner.Scan() {
		if err := r.scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return strings.TrimSpace(r.scanner.Text()), nil
}

func (r *Runner) printMenu() {
	fmt.Fprintln(r.out, "\n培训作业成绩系统")
	fmt.Fprintln(r.out, "1. 录入  2. 查看  3. 查询  4. 排序")
	fmt.Fprintln(r.out, "5. 删除  6. 修改  7. 保存  8. 读取  0. 退出")
}

func (r *Runner) printStudents(students []grades.Student) {
	if len(students) == 0 {
		fmt.Fprintln(r.out, "暂无学员成绩")
		return
	}
	fmt.Fprintln(r.out, "编号\t姓名\t作业一\t作业二\t结课测试\t平均分")
	for _, student := range students {
		fmt.Fprintf(
			r.out,
			"%s\t%s\t%d\t%d\t%d\t%s\n",
			student.ID,
			student.Name,
			student.Homework1,
			student.Homework2,
			student.FinalExam,
			student.Average().StringFixed(2),
		)
	}
}
