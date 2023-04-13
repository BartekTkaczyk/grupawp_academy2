package test

import (
	academy "github.com/grupawp/akademia-programowania/Golang/zadania/academy2"
	"github.com/grupawp/akademia-programowania/Golang/zadania/academy2/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type GradeStudentCaseData struct {
	name           string
	yearFn         uint8
	finalGradeFn   int
	getFnErr       error
	expectedResult error
}

func extractStudentsName(students []GradeStudentCaseData) []string {
	studentsNames := make([]string, len(students))

	index := 0
	for _, student := range students {
		studentsNames[index] = student.name
		index++
	}

	return studentsNames
}

func TestGradeStudent(t *testing.T) {
	testCases := map[string]GradeStudentCaseData{
		"Student not found": {
			yearFn:         2,
			finalGradeFn:   4,
			getFnErr:       academy.ErrStudentNotFound,
			expectedResult: nil,
		},
		"Invalid grade": {
			yearFn:         2,
			finalGradeFn:   0,
			getFnErr:       nil,
			expectedResult: academy.ErrInvalidGrade,
		},
		"Year sentence": {
			yearFn:         uint8(2),
			finalGradeFn:   4,
			getFnErr:       nil,
			expectedResult: nil,
		},
		"College graduation": {
			yearFn:         3,
			finalGradeFn:   4,
			getFnErr:       nil,
			expectedResult: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			student := mocks.NewStudent(t)

			repo.On("Get", mock.Anything).Return(student, tc.getFnErr)

			switch name {
			case "Invalid grade":
				student.On("FinalGrade").Return(tc.finalGradeFn)
			case "Year sentence":
				student.On("FinalGrade").Return(tc.finalGradeFn)
				student.On("Year").Return(tc.yearFn)
				student.On("Name").Return(name)
				repo.On("Save", name, tc.yearFn+1).Return(nil)
			case "College graduation":
				student.On("FinalGrade").Return(tc.finalGradeFn)
				student.On("Year").Return(tc.yearFn)
				student.On("Name").Return(name)
				repo.On("Graduate", name).Return(nil)
			}

			assert.Equal(t, tc.expectedResult, academy.GradeStudent(repo, name))
		})
	}
}

func TestGardeYear_CompleteYear(t *testing.T) {
	testCases := []GradeStudentCaseData{
		{
			name:           "John3",
			yearFn:         3,
			finalGradeFn:   1,
			getFnErr:       nil,
			expectedResult: nil,
		},
		{
			name:           "John4",
			yearFn:         3,
			finalGradeFn:   4,
			getFnErr:       nil,
			expectedResult: nil,
		},
	}

	names := extractStudentsName(testCases)
	repo := mocks.NewRepository(t)

	repo.On("List", uint8(3)).Return(names, nil)
	for _, tc := range testCases {
		student := mocks.NewStudent(t)

		repo.On("Get", tc.name).Return(student, nil)
		student.On("FinalGrade").Return(tc.finalGradeFn)

		student.On("Year").Return(tc.yearFn)
		student.On("Name").Return(tc.name)

		if tc.finalGradeFn == 1 {
			repo.On("Save", tc.name, tc.yearFn).Return(nil)
		} else {
			repo.On("Graduate", tc.name).Return(nil)
		}
	}

	assert.Equal(t, nil, academy.GradeYear(repo, 3))
}

func TestGradeYear_InvalidGrade(t *testing.T) {
	testCases := []GradeStudentCaseData{
		{
			name:           "John1",
			yearFn:         3,
			finalGradeFn:   4,
			getFnErr:       academy.ErrStudentNotFound,
			expectedResult: nil,
		},
		{
			name:           "John3",
			yearFn:         3,
			finalGradeFn:   4,
			getFnErr:       nil,
			expectedResult: nil,
		},
		{
			name:           "John4",
			yearFn:         3,
			finalGradeFn:   4,
			getFnErr:       nil,
			expectedResult: nil,
		},
		{
			name:           "John2",
			yearFn:         3,
			finalGradeFn:   0,
			getFnErr:       nil,
			expectedResult: academy.ErrInvalidGrade,
		},
	}

	studentsNames := extractStudentsName(testCases)
	repo := mocks.NewRepository(t)

	repo.On("List", uint8(3)).Return(studentsNames, nil)

	for _, tc := range testCases {
		student := mocks.NewStudent(t)

		if tc.getFnErr != nil {
			repo.On("Get", tc.name).Return(nil, tc.getFnErr)
			continue
		}
		repo.On("Get", tc.name).Return(student, tc.getFnErr)
		student.On("FinalGrade").Return(tc.finalGradeFn)

		if tc.finalGradeFn < 1 || tc.finalGradeFn > 5 {
			continue
		}

		student.On("Name").Return(tc.name)
		student.On("Year").Return(tc.yearFn)

		if tc.finalGradeFn == 1 {
			repo.On("Save", tc.name, tc.yearFn).Return(nil)
		}

		student.On("Name").Return(tc.name)
		repo.On("Graduate", tc.name).Return(nil)
	}

	assert.Equal(t, academy.ErrInvalidGrade, academy.GradeYear(repo, uint8(3)))
}

func TestGradeYear_StudentNotFound(t *testing.T) {
	testCases := []GradeStudentCaseData{
		{
			name:           "John1",
			yearFn:         2,
			finalGradeFn:   4,
			getFnErr:       nil,
			expectedResult: nil,
		},
		{
			name:           "John3",
			yearFn:         2,
			finalGradeFn:   4,
			getFnErr:       nil,
			expectedResult: nil,
		},
		{
			name:           "John4",
			yearFn:         2,
			finalGradeFn:   4,
			getFnErr:       academy.ErrStudentNotFound,
			expectedResult: nil,
		},
	}

	studentsNames := extractStudentsName(testCases)
	repo := mocks.NewRepository(t)

	repo.On("List", uint8(2)).Return(studentsNames, nil)

	for _, tc := range testCases {
		student := mocks.NewStudent(t)

		if tc.getFnErr != nil {
			repo.On("Get", tc.name).Return(nil, tc.getFnErr)
			continue
		}
		repo.On("Get", tc.name).Return(student, tc.getFnErr)
		student.On("FinalGrade").Return(tc.finalGradeFn)

		student.On("Name").Return(tc.name)
		student.On("Year").Return(tc.yearFn)
		if tc.finalGradeFn == 1 {
			repo.On("Save", tc.name, tc.yearFn).Return(nil)
		}
		repo.On("Save", tc.name, tc.yearFn+1).Return(nil)

	}

	assert.Equal(t, nil, academy.GradeYear(repo, 2))
}
