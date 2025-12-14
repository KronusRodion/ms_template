package usecase

import (
	"testing"
	"time"

	"ms_template/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/google/uuid"
)

// Мок репозитория для тестирования юзкейса
type MockNoteRepository struct {
	mock.Mock
}

func (m *MockNoteRepository) AddNote(note domain.Note) string {
	args := m.Called(note)
	return args.String(0)
}

func (m *MockNoteRepository) GetNotes() []domain.Note {
	args := m.Called()
	return args.Get(0).([]domain.Note)
}

type BasicUsecaseTestSuite struct {
	suite.Suite
	mockRepo *MockNoteRepository
	usecase  NoteUsecase
}

func TestBasicUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(BasicUsecaseTestSuite))
}

func (s *BasicUsecaseTestSuite) SetupTest() {
	s.mockRepo = new(MockNoteRepository)
	s.usecase = NewBasic(s.mockRepo)
}

func (s *BasicUsecaseTestSuite) TestNewBasic() {
	// Arrange
	repo := new(MockNoteRepository)

	// Act
	usecaseInstance := NewBasic(repo)

	// Assert
	assert.NotNil(s.T(), usecaseInstance)
	assert.IsType(s.T(), &Basic{}, usecaseInstance)
}

func (s *BasicUsecaseTestSuite) TestGetNotes() {
	// Arrange
	expectedNotes := []domain.Note{
		{
			ID:        "note-1",
			Title:     "Test Note 1",
			Content:   "Content 1",
			CreatedAt: time.Now(),
			UserID:    "user-1",
		},
		{
			ID:        "note-2",
			Title:     "Test Note 2",
			Content:   "Content 2",
			CreatedAt: time.Now().Add(time.Hour),
			UserID:    "user-2",
		},
	}

	s.mockRepo.On("GetNotes").Return(expectedNotes)

	// Act
	notes := s.usecase.GetNotes("user-1")

	// Assert
	assert.Len(s.T(), notes, 2)
	assert.Equal(s.T(), expectedNotes[0].ID, notes[0].ID)
	assert.Equal(s.T(), expectedNotes[1].ID, notes[1].ID)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *BasicUsecaseTestSuite) TestGetNotes_EmptyResult() {
	// Arrange
	s.mockRepo.On("GetNotes").Return([]domain.Note{})

	// Act
	notes := s.usecase.GetNotes("user-1")

	// Assert
	assert.Empty(s.T(), notes)
	assert.Len(s.T(), notes, 0)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *BasicUsecaseTestSuite) TestAddNote() {
	// Arrange
	note := domain.Note{
		Title:   "Test Title",
		Content: "Test Content",
		UserID:  "user-1",
	}
	
	// Используем mock.Anything чтобы не проверять точные значения ID и CreatedAt
	s.mockRepo.On("AddNote", mock.AnythingOfType("domain.Note")).
		Return("generated-id").
		Run(func(args mock.Arguments) {
			// Проверяем что ID и CreatedAt установлены
			noteArg := args.Get(0).(domain.Note)
			assert.NotEmpty(s.T(), noteArg.ID)
			assert.NotZero(s.T(), noteArg.CreatedAt)
			assert.Equal(s.T(), note.Title, noteArg.Title)
			assert.Equal(s.T(), note.Content, noteArg.Content)
			assert.Equal(s.T(), note.UserID, noteArg.UserID)
		})

	// Act
	resultID := s.usecase.AddNote(note)

	// Assert
	assert.NotEmpty(s.T(), resultID)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *BasicUsecaseTestSuite) TestAddNote_WithExistingID() {
	// Arrange
	existingID := "existing-id"
	note := domain.Note{
		ID:       existingID,
		Title:    "Test Title",
		Content:  "Test Content",
		UserID:   "user-1",
		CreatedAt: time.Now(), // Уже установлено
	}
	
	// Используем mock.MatchedBy для более точной проверки
	s.mockRepo.On("AddNote", mock.MatchedBy(func(n domain.Note) bool {
		// Проверяем что ID перезаписан новым UUID
		_, err := uuid.Parse(n.ID)
		return err == nil && n.ID != existingID &&
			n.Title == note.Title &&
			n.Content == note.Content &&
			n.UserID == note.UserID &&
			!n.CreatedAt.IsZero()
	})).Return("new-id")

	// Act
	resultID := s.usecase.AddNote(note)

	// Assert
	assert.NotEmpty(s.T(), resultID)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *BasicUsecaseTestSuite) TestAddNote_WithDifferentUsers() {
	// Arrange
	notes := []domain.Note{
		{
			Title:   "User 1 Note",
			Content: "Content",
			UserID:  "user-1",
		},
		{
			Title:   "User 2 Note",
			Content: "Content",
			UserID:  "user-2",
		},
	}

	for _, note := range notes {
		s.mockRepo.On("AddNote", mock.MatchedBy(func(n domain.Note) bool {
			return n.Title == note.Title &&
				n.Content == note.Content &&
				n.UserID == note.UserID &&
				!(n.ID == "") &&
				!n.CreatedAt.IsZero()
		})).Return("generated-id")
	}

	// Act & Assert
	for _, note := range notes {
		resultID := s.usecase.AddNote(note)
		assert.NotEmpty(s.T(), resultID)
	}

	s.mockRepo.AssertExpectations(s.T())
}

func (s *BasicUsecaseTestSuite) TestGetNotes_UserIDParameterIgnored() {
	// Arrange
	allNotes := []domain.Note{
		{ID: "1", UserID: "user-1"},
		{ID: "2", UserID: "user-2"},
		{ID: "3", UserID: "user-3"},
	}

	s.mockRepo.On("GetNotes").Return(allNotes).Times(3)

	// Act
	notes1 := s.usecase.GetNotes("user-1")
	notes2 := s.usecase.GetNotes("user-2")
	notes3 := s.usecase.GetNotes("")

	// Assert
	assert.Len(s.T(), notes1, 3)
	assert.Len(s.T(), notes2, 3)
	assert.Len(s.T(), notes3, 3)
	s.mockRepo.AssertExpectations(s.T())
}

// Тест на граничные условия
func (s *BasicUsecaseTestSuite) TestAddNote_EmptyFields() {
	// Arrange
	testCases := []struct {
		name string
		note domain.Note
	}{
		{
			name: "Empty Title",
			note: domain.Note{
				Title:   "",
				Content: "Content",
				UserID:  "user-1",
			},
		},
		{
			name: "Empty Content",
			note: domain.Note{
				Title:   "Title",
				Content: "",
				UserID:  "user-1",
			},
		},
		{
			name: "Empty UserID",
			note: domain.Note{
				Title:   "Title",
				Content: "Content",
				UserID:  "",
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			s.mockRepo.On("AddNote", mock.MatchedBy(func(n domain.Note) bool {
				return n.Title == tc.note.Title &&
					n.Content == tc.note.Content &&
					n.UserID == tc.note.UserID &&
					!(n.ID == "") &&
					!n.CreatedAt.IsZero()
			})).Return("generated-id")

			// Act
			resultID := s.usecase.AddNote(tc.note)

			// Assert
			assert.NotEmpty(t, resultID)
		})
	}

	s.mockRepo.AssertExpectations(s.T())
}

func (s *BasicUsecaseTestSuite) TestAddNote_ReturnsGeneratedID() {
	// Arrange
	note := domain.Note{
		Title:   "Test Title",
		Content: "Test Content",
		UserID:  "user-1",
	}
	
	expectedID := "test-generated-id"
	s.mockRepo.On("AddNote", mock.AnythingOfType("domain.Note")).Return(expectedID)

	// Act
	resultID := s.usecase.AddNote(note)

	// Assert
	assert.Equal(s.T(), expectedID, resultID)
	s.mockRepo.AssertExpectations(s.T())
}