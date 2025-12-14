package repository

import (
	"testing"
	"time"
	
	"ms_template/internal/domain"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PostgresRepoTestSuite struct {
	suite.Suite
	repo *Postgres
}

func TestPostgresRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresRepoTestSuite))
}

func (s *PostgresRepoTestSuite) SetupTest() {
	s.repo = NewPostgresRepo()
}

func (s *PostgresRepoTestSuite) TearDownTest() {
	// Очистка данных после каждого теста
	// В реальной реализации нужно очистить notes
}

func (s *PostgresRepoTestSuite) TestNewPostgresRepo() {
	// Act & Assert
	assert.NotNil(s.T(), s.repo)
}

func (s *PostgresRepoTestSuite) TestGetNotes_EmptyRepository() {
	// Act
	notes := s.repo.GetNotes()

	// Assert
	assert.Empty(s.T(), notes)
	assert.Len(s.T(), notes, 0)
}

func (s *PostgresRepoTestSuite) TestAddNote_SingleNote() {
	// Arrange
	note := domain.Note{
		ID:        "test-id-1",
		Title:     "Test Note",
		Content:   "Test Content",
		CreatedAt: time.Now().UTC(),
		UserID:    "user-1",
	}

	// Act
	s.repo.AddNote(note)
	notes := s.repo.GetNotes()

	// Assert
	assert.Len(s.T(), notes, 1)
	assert.Equal(s.T(), note.ID, notes[0].ID)
	assert.Equal(s.T(), note.Title, notes[0].Title)
	assert.Equal(s.T(), note.Content, notes[0].Content)
	assert.Equal(s.T(), note.UserID, notes[0].UserID)
}

func (s *PostgresRepoTestSuite) TestAddNote_MultipleNotes() {
	// Arrange
	notesToAdd := []domain.Note{
		{
			ID:        "test-id-1",
			Title:     "Note 1",
			Content:   "Content 1",
			CreatedAt: time.Now(),
			UserID:    "user-1",
		},
		{
			ID:        "test-id-2",
			Title:     "Note 2",
			Content:   "Content 2",
			CreatedAt: time.Now().Add(time.Hour),
			UserID:    "user-2",
		},
		{
			ID:        "test-id-3",
			Title:     "Note 3",
			Content:   "Content 3",
			CreatedAt: time.Now().Add(2 * time.Hour),
			UserID:    "user-1",
		},
	}

	// Act
	for _, note := range notesToAdd {
		s.repo.AddNote(note)
	}
	notes := s.repo.GetNotes()

	// Assert
	assert.Len(s.T(), notes, 3)
	
	// Проверяем, что все заметки добавлены
	noteMap := make(map[string]domain.Note)
	for _, note := range notes {
		noteMap[note.ID] = note
	}
	
	for _, expectedNote := range notesToAdd {
		actualNote, exists := noteMap[expectedNote.ID]
		assert.True(s.T(), exists)
		assert.Equal(s.T(), expectedNote.Title, actualNote.Title)
		assert.Equal(s.T(), expectedNote.Content, actualNote.Content)
		assert.Equal(s.T(), expectedNote.UserID, actualNote.UserID)
	}
}

func (s *PostgresRepoTestSuite) TestAddNote_UpdateExistingNote() {
	// Arrange
	initialNote := domain.Note{
		ID:        "test-id",
		Title:     "Initial Title",
		Content:   "Initial Content",
		CreatedAt: time.Now(),
		UserID:    "user-1",
	}

	updatedNote := domain.Note{
		ID:        "test-id", // Тот же ID
		Title:     "Updated Title",
		Content:   "Updated Content",
		CreatedAt: time.Now().Add(time.Hour),
		UserID:    "user-1",
	}

	// Act
	s.repo.AddNote(initialNote)
	s.repo.AddNote(updatedNote) // Перезаписываем
	notes := s.repo.GetNotes()

	// Assert
	assert.Len(s.T(), notes, 1) // Все еще одна запись
	assert.Equal(s.T(), "Updated Title", notes[0].Title)
	assert.Equal(s.T(), "Updated Content", notes[0].Content)
}

func (s *PostgresRepoTestSuite) TestAddNote_ConcurrentAccess() {
	// Arrange
	numGoroutines := 100
	notesPerGoroutine := 10
	done := make(chan bool, numGoroutines)

	// Act
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			for j := 0; j < notesPerGoroutine; j++ {
				note := domain.Note{
					ID:      s.generateNoteID(goroutineID, j),
					Title:   "Note",
					Content: "Content",
					UserID:  "user",
				}
				s.repo.AddNote(note)
			}
			done <- true
		}(i)
	}

	// Ждем завершения всех горутин
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Assert
	notes := s.repo.GetNotes()
	assert.Len(s.T(), notes, numGoroutines*notesPerGoroutine)
}

func (s *PostgresRepoTestSuite) generateNoteID(goroutineID, noteIndex int) string {
	return string(rune('A'+goroutineID)) + string(rune('0'+noteIndex))
}

func (s *PostgresRepoTestSuite) TestGetNotes_ConcurrentReadWrite() {
	// Arrange
	done := make(chan bool)
	readComplete := make(chan bool)
	
	// Запускаем горутину для записи
	go func() {
		for i := 0; i < 1000; i++ {
			note := domain.Note{
				ID:      string(rune('A' + i%26)),
				Title:   "Note",
				Content: "Content",
				UserID:  "user",
			}
			s.repo.AddNote(note)
			time.Sleep(time.Microsecond)
		}
		done <- true
	}()
	
	// Запускаем горутину для чтения
	go func() {
		for i := 0; i < 100; i++ {
			notes := s.repo.GetNotes()
			_ = len(notes) // Просто читаем
			time.Sleep(time.Millisecond)
		}
		readComplete <- true
	}()
	
	// Assert - проверяем, что нет гонки данных
	<-readComplete
	<-done
	// Если тест не падает с data race - всё хорошо
}