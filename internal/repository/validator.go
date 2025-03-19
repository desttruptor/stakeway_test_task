package repository

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"stakeway_test_task/internal/models"
	"time"
)

type ValidatorRepository struct {
	db *sql.DB
}

func NewValidatorRepository(dbPath string) (*ValidatorRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	err = createTables(db)
	if err != nil {
		return nil, err
	}

	return &ValidatorRepository{db: db}, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS validator_requests (
			id TEXT PRIMARY KEY,
			num_validators INTEGER,
			fee_recipient TEXT,
			status TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			error_message TEXT
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS validator_keys (
			id TEXT PRIMARY KEY,
			request_id TEXT,
			key TEXT,
			fee_recipient TEXT,
			FOREIGN KEY (request_id) REFERENCES validator_requests (id)
		)
	`)
	return err
}

func (r *ValidatorRepository) CreateRequest(request *models.ValidatorRequest) error {
	_, err := r.db.Exec(
		"INSERT INTO validator_requests (id, num_validators, fee_recipient, status, created_at, updated_at, error_message) VALUES (?, ?, ?, ?, ?, ?, ?)",
		request.ID, request.NumValidators, request.FeeRecipient, request.Status, time.Now(), time.Now(), request.ErrorMessage,
	)
	return err
}

func (r *ValidatorRepository) GetRequestByID(id string) (*models.ValidatorRequest, error) {
	row := r.db.QueryRow("SELECT id, num_validators, fee_recipient, status, created_at, updated_at, error_message FROM validator_requests WHERE id = ?", id)

	var req models.ValidatorRequest
	var status string
	err := row.Scan(&req.ID, &req.NumValidators, &req.FeeRecipient, &status, &req.CreatedAt, &req.UpdatedAt, &req.ErrorMessage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("request not found")
		}
		return nil, err
	}

	req.Status = models.Status(status)
	return &req, nil
}

func (r *ValidatorRepository) UpdateRequestStatus(id string, status models.Status, errorMessage string) error {
	_, err := r.db.Exec(
		"UPDATE validator_requests SET status = ?, updated_at = ?, error_message = ? WHERE id = ?",
		status, time.Now(), errorMessage, id,
	)
	return err
}

func (r *ValidatorRepository) SaveValidatorKey(key *models.ValidatorKey) error {
	_, err := r.db.Exec(
		"INSERT INTO validator_keys (id, request_id, key, fee_recipient) VALUES (?, ?, ?, ?)",
		key.ID, key.RequestID, key.Key, key.FeeRecipient,
	)
	return err
}

func (r *ValidatorRepository) GetKeysByRequestID(requestID string) ([]string, error) {
	rows, err := r.db.Query("SELECT key FROM validator_keys WHERE request_id = ?", requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return keys, nil
}

func (r *ValidatorRepository) CheckHealth() error {
	return r.db.Ping()
}

func (r *ValidatorRepository) Close() error {
	return r.db.Close()
}
