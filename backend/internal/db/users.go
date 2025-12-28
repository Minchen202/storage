package db

import (
	"database/sql"
	"p2p-storage/backend/pkg/models"
	"time"

	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func CreateUser(db *sql.DB, user *models.User) error {
	log := logrus.WithFields(logrus.Fields{
		"function": "CreateUser",
		"email":    user.Email,
	})
	log.Debug("Attempting to create user in database...")

	start := time.Now()
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at, storage_contributed, storage_quota, total_uptime_seconds, last_online`

	err := db.QueryRow(
		query,
		user.Email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.StorageContributed,
		&user.StorageQuota,
		&user.TotalUptimeSeconds,
		&user.LastOnline,
	)

	duration := time.Since(start)
	if err != nil {
		log.WithError(err).WithField("query_ms", duration.Milliseconds()).Error("Failed to create user in database")
		return err
	}

	log.WithFields(logrus.Fields{
		"user_id":  user.ID,
		"query_ms": duration.Milliseconds(),
	}).Debug("User created successfully in database")
	return nil
}

func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	log := logrus.WithFields(logrus.Fields{
		"function": "GetUserByEmail",
		"email":    email,
	})
	log.Debug("Looking up user by email in database...")

	start := time.Now()
	user := &models.User{}
	query := `SELECT id, email, password_hash, created_at, storage_contributed, storage_quota, total_uptime_seconds, last_online FROM users WHERE email = $1`
	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.StorageContributed,
		&user.StorageQuota,
		&user.TotalUptimeSeconds,
		&user.LastOnline,
	)

	duration := time.Since(start)
	if err != nil {
		if err == sql.ErrNoRows {
			log.WithField("query_ms", duration.Milliseconds()).Debug("User not found in database")
		} else {
			log.WithError(err).WithField("query_ms", duration.Milliseconds()).Error("Database error while looking up user by email")
		}
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"user_id":  user.ID,
		"query_ms": duration.Milliseconds(),
	}).Debug("User found in database")
	return user, nil
}

func GetUserStats(db *sql.DB, userID string) (*models.User, error) {
	log := logrus.WithFields(logrus.Fields{
		"function": "GetUserStats",
		"user_id":  userID,
	})
	log.Debug("Looking up user stats in database...")

	start := time.Now()
	user := &models.User{}
	query := `SELECT id, email, created_at, storage_contributed, storage_quota, total_uptime_seconds, last_online FROM users WHERE id = $1`
	err := db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.StorageContributed,
		&user.StorageQuota,
		&user.TotalUptimeSeconds,
		&user.LastOnline,
	)

	duration := time.Since(start)
	if err != nil {
		if err == sql.ErrNoRows {
			log.WithField("query_ms", duration.Milliseconds()).Debug("User not found in database")
		} else {
			log.WithError(err).WithField("query_ms", duration.Milliseconds()).Error("Database error while looking up user stats")
		}
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"email":    user.Email,
		"query_ms": duration.Milliseconds(),
	}).Debug("User stats retrieved successfully")
	return user, nil
}
