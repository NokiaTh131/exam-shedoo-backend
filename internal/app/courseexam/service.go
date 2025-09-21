package courseexam

import (
	"fmt"
	"shedoo-backend/internal/models"
	"shedoo-backend/internal/repositories"
)

type CourseExamService struct {
    repo *repositories.CourseExamRepository
}

func NewCourseExamService(repo *repositories.CourseExamRepository) *CourseExamService {
    return &CourseExamService{repo: repo}
}

func (s *CourseExamService) CreateExam(exam *models.CourseExam) (*models.CourseExam, error) {
    
    if exam.CourseCode == "" || exam.LecSection == "" || exam.LabSection == "" {
        return nil, fmt.Errorf("courseCode, lecSection, labSection are required")
    }
    
    existing, err := s.repo.GetByCourseSections(exam.CourseCode, exam.LecSection, exam.LabSection)
    if err == nil && existing != nil {
        return nil, fmt.Errorf("exam already exists for courseCode %s, lecSection %s, labSection %s", exam.CourseCode, exam.LecSection, exam.LabSection)
    }
    
    if err := s.repo.Create(exam); err != nil {
        return nil, err
    }
    return exam, nil
}

func (s *CourseExamService) GetExamByCourseSections(courseCode, lecSection, labSection string) (*models.CourseExam, error) {
    return s.repo.GetByCourseSections(courseCode, lecSection, labSection)
}

func (s *CourseExamService) UpdateExam(id uint, updates *models.CourseExam) error {
    return s.repo.UpdateByID(id, updates)
}

func (s *CourseExamService) FindByID(id uint) (*models.CourseExam, error) {
    return s.repo.FindByID(id)
}