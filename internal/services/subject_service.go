package services

import (
	"eduqr-backend/internal/models"
	"eduqr-backend/internal/repositories"
	"errors"
)

type SubjectService struct {
	subjectRepo *repositories.SubjectRepository
}

func NewSubjectService(subjectRepo *repositories.SubjectRepository) *SubjectService {
	return &SubjectService{subjectRepo: subjectRepo}
}

// GetAllSubjects récupère toutes les matières
func (s *SubjectService) GetAllSubjects() ([]models.SubjectResponse, error) {
	subjects, err := s.subjectRepo.GetAllSubjects()
	if err != nil {
		return nil, err
	}

	responses := make([]models.SubjectResponse, len(subjects))
	for i, subject := range subjects {
		responses[i] = subject.ToSubjectResponse()
	}

	return responses, nil
}

// GetSubjectByID récupère une matière par son ID
func (s *SubjectService) GetSubjectByID(id uint) (*models.SubjectResponse, error) {
	subject, err := s.subjectRepo.GetSubjectByID(id)
	if err != nil {
		return nil, err
	}

	response := subject.ToSubjectResponse()
	return &response, nil
}

// CreateSubject crée une nouvelle matière
func (s *SubjectService) CreateSubject(req *models.CreateSubjectRequest) (*models.SubjectResponse, error) {
	// Vérifier si le nom existe déjà
	exists, err := s.subjectRepo.CheckSubjectExists(req.Name, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("une matière avec ce nom existe déjà")
	}

	// Créer la matière
	subject := &models.Subject{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	}

	err = s.subjectRepo.CreateSubject(subject)
	if err != nil {
		return nil, err
	}

	response := subject.ToSubjectResponse()
	return &response, nil
}

// UpdateSubject met à jour une matière
func (s *SubjectService) UpdateSubject(id uint, req *models.UpdateSubjectRequest) (*models.SubjectResponse, error) {
	// Vérifier si la matière existe
	subject, err := s.subjectRepo.GetSubjectByID(id)
	if err != nil {
		return nil, err
	}

	// Vérifier si le nouveau nom existe déjà (sauf pour cette matière)
	exists, err := s.subjectRepo.CheckSubjectExists(req.Name, &id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("une matière avec ce nom existe déjà")
	}

	// Mettre à jour les champs
	subject.Name = req.Name
	subject.Code = req.Code
	subject.Description = req.Description

	err = s.subjectRepo.UpdateSubject(subject)
	if err != nil {
		return nil, err
	}

	response := subject.ToSubjectResponse()
	return &response, nil
}

// DeleteSubject supprime une matière
func (s *SubjectService) DeleteSubject(id uint) error {
	// Vérifier si la matière existe
	_, err := s.subjectRepo.GetSubjectByID(id)
	if err != nil {
		return err
	}

	// Vérifier si la matière est utilisée dans des cours
	inUse, err := s.subjectRepo.CheckSubjectInUse(id)
	if err != nil {
		return err
	}
	if inUse {
		return errors.New("cette matière ne peut pas être supprimée car elle est utilisée dans des cours")
	}

	return s.subjectRepo.DeleteSubject(id)
}
