package models

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Session struct {
	ID        int       `gorm:"primaryKey"`
	Date      time.Time `gorm:"not null;index"       binding:"required"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	UserID    int       `gorm:"not null;index"`
	Climbs    []Climb   `gorm:"foreignKey:SessionID"` // not a column, is a relationship
	Notes     string    `                                               json:"notes"`
}

type SessionSummary struct {
	ID   int       `gorm:"primaryKey"`
	Date time.Time `gorm:"not null;index" binding:"required"`
	Load float64
}

func FindAllSessions(db *gorm.DB, userID, offset, limit int) (*[]Session, error) {
	var sessions []Session

	// Execute the query and check for errors
	err := db.
		Preload("Climbs").
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&sessions).Error

	if err != nil {
		log.Printf("Database error: %v", err) // Log the error if query fails
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Error finding sessions for user with id %d", userID),
		)
	}

	return &sessions, nil
}

func (s *Session) FindByID(db *gorm.DB, userID, sessionID, offset, limit int) error {
	err := db.
		// climbs can be paginated within a session
		Preload("Climbs", func(db *gorm.DB) *gorm.DB {
			return db.Limit(limit).Offset(offset).Order("created_at DESC")
		}).
		Where("user_id = ? AND id = ?", userID, sessionID).
		Take(s).Error

	if err != nil {
		log.Printf("Database error: %v", err) // Log the error if query fails
		return errors.Wrap(
			err,
			fmt.Sprintf("Error finding sessions with id %d for user with id %d", sessionID, userID),
		)
	}

	return nil
}

func (s *Session) FindByDate(db *gorm.DB, userID int, sessionDate time.Time) error {
	err := db.
		Preload("Climbs").
		Where("user_id = ? AND date = ?", userID, sessionDate).
		Take(s).Error

	if err != nil {
		log.Printf("Database error: %v", err) // Log the error if query fails
		return errors.Wrap(
			err,
			fmt.Sprintf(
				"Error finding session on date %s for user with id %d",
				sessionDate,
				userID,
			),
		)
	}

	return nil
}

// session summaries
func FindSessionSummaries(
	db *gorm.DB,
	userID int,
	startDate, endDate time.Time,
) (*[]SessionSummary, error) {
	var summaries []SessionSummary

	err := db.Table("sessions").
		Select("sessions.id, sessions.date, SUM(climbs.load) as load").
		Joins("JOIN climbs on climbs.session_id = sessions.id").
		Where("sessions.user_id = ? AND sessions.date BETWEEN ? AND ?", userID, startDate, endDate).
		Group("sessions.id").
		Scan(&summaries).Error

	if err != nil {
		log.Printf("Database error: %v", err) // Log the error if query fails
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Error finding session summaries for user with id %d", userID),
		)
	}

	return &summaries, nil
}
