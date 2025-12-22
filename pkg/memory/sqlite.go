package memory

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/openai/openai-go/v3/responses"
)

const (
	// DefaultDBPath is the default path for the SQLite database, using in-memory storage.
	DefaultDBPath = ":memory:"
	// DefaultSessionsTable is the default table name for session metadata.
	DefaultSessionsTable = "agent_sessions"
	// DefaultMessagesTable is the default table name for message data.
	DefaultMessagesTable = "agent_messages"
)

// SQLiteSessionConfig holds configuration for creating a new SQLiteSession.
type SQLiteSessionConfig struct {
	// SessionID is the unique identifier for the session (required).
	SessionID string
	// DBPath is the path to the SQLite database file; defaults to ":memory:" for in-memory storage.
	DBPath string
	// SessionsTable is the table name for session metadata; defaults to "agent_sessions".
	SessionsTable string
	// MessagesTable is the table name for message data; defaults to "agent_messages".
	MessagesTable string
}

// Validate checks if the configuration is valid.
func (c *SQLiteSessionConfig) Validate() error {
	if c.SessionID == "" {
		return ErrSessionNotFound
	}
	// Additional validation for table names could be added here to prevent invalid SQL identifiers.
	return nil
}

var _ Session = (*SQLiteSession)(nil)

// SQLiteSession implements a session storage using SQLite.
// It manages conversation history with thread-safety using a mutex.
// The session can be backed by an on-disk or in-memory database.
type SQLiteSession struct {
	sessionID     string
	db            *sql.DB
	sessionsTable string
	messagesTable string
	isMemoryDB    bool
	mu            sync.Mutex
}

// NewSQLiteSession creates and initializes a new SQLiteSession based on the provided configuration.
// It opens the database connection, initializes the schema, and ensures the session entry exists.
// The context is used for database operations during initialization.
func NewSQLiteSession(ctx context.Context, config SQLiteSessionConfig) (*SQLiteSession, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	if config.DBPath == "" {
		config.DBPath = DefaultDBPath
	}
	if config.SessionsTable == "" {
		config.SessionsTable = DefaultSessionsTable
	}
	if config.MessagesTable == "" {
		config.MessagesTable = DefaultMessagesTable
	}

	db, err := sql.Open("sqlite3", config.DBPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDatabaseOpen, err)
	}

	s := &SQLiteSession{
		sessionID:     config.SessionID,
		db:            db,
		sessionsTable: config.SessionsTable,
		messagesTable: config.MessagesTable,
		isMemoryDB:    config.DBPath == ":memory:",
	}

	if err := s.initDB(ctx); err != nil {
		db.Close()
		return nil, err
	}

	// Ensure session entry exists upon creation.
	if err := s.ensureSessionExists(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return s, nil
}

// initDB initializes the database schema if it does not exist.
// It creates tables for sessions and messages, along with necessary indexes.
func (s *SQLiteSession) initDB(ctx context.Context) error {
	// Use parameterized queries where possible, but table names require formatting.
	// Assuming table names are trusted (set via config).
	createSessionsTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			session_id TEXT PRIMARY KEY,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`, s.sessionsTable)

	createMessagesTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id TEXT NOT NULL,
			message_data TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (session_id) REFERENCES %s (session_id) ON DELETE CASCADE
		)`, s.messagesTable, s.sessionsTable)

	createIndex := fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS idx_%s_session_id 
		ON %s (session_id, created_at)`, s.messagesTable, s.messagesTable)

	_, err := s.db.ExecContext(ctx, createSessionsTable+";"+createMessagesTable+";"+createIndex)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDatabaseInit, err)
	}
	return nil
}

// ensureSessionExists inserts the session entry if it does not already exist.
func (s *SQLiteSession) ensureSessionExists(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.ExecContext(ctx,
		fmt.Sprintf(`INSERT OR IGNORE INTO %s (session_id) VALUES (?)`, s.sessionsTable),
		s.sessionID,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrOperationFailed, err)
	}
	return nil
}

// SessionID returns the unique identifier for this session.
func (s *SQLiteSession) SessionID() string {
	return s.sessionID
}

// GetItems retrieves the conversation history items for this session.
// If limit <= 0, all items are returned in chronological order (oldest first).
// If limit > 0, the latest 'limit' items are returned in chronological order.
// Items are deserialized from JSON stored in the database.
func (s *SQLiteSession) GetItems(ctx context.Context, limit int) ([]responses.ResponseInputItemUnionParam, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var query string
	var args []any
	if limit > 0 {
		query = fmt.Sprintf(`SELECT message_data FROM %s WHERE session_id = ? ORDER BY created_at DESC LIMIT ?`, s.messagesTable)
		args = []any{s.sessionID, limit}
	} else {
		query = fmt.Sprintf(`SELECT message_data FROM %s WHERE session_id = ? ORDER BY created_at ASC`, s.messagesTable)
		args = []any{s.sessionID}
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOperationFailed, err)
	}
	defer rows.Close()

	var items []responses.ResponseInputItemUnionParam
	for rows.Next() {
		var messageData string
		if err := rows.Scan(&messageData); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrOperationFailed, err)
		}
		var item responses.ResponseInputItemUnionParam
		if err := json.Unmarshal([]byte(messageData), &item); err != nil {
			// Log or handle invalid JSON; for now, skip.
			continue
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOperationFailed, err)
	}

	// If limited, reverse to chronological order (oldest first).
	if limit > 0 {
		for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
			items[i], items[j] = items[j], items[i]
		}
	}

	return items, nil
}

// AddItems appends new items to the conversation history.
// It uses a transaction to ensure atomicity and updates the session's updated_at timestamp.
func (s *SQLiteSession) AddItems(ctx context.Context, items []responses.ResponseInputItemUnionParam) error {
	if len(items) == 0 {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTransactionFailed, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Prepare insert statement for messages.
	insertQuery := fmt.Sprintf(`INSERT INTO %s (session_id, message_data) VALUES (?, ?)`, s.messagesTable)
	stmt, err := tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrOperationFailed, err)
	}
	defer stmt.Close()

	for _, item := range items {
		data, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidItemData, err)
		}
		if _, err := stmt.ExecContext(ctx, s.sessionID, string(data)); err != nil {
			return fmt.Errorf("%w: %w", ErrOperationFailed, err)
		}
	}

	// Update session timestamp.
	updateQuery := fmt.Sprintf(`UPDATE %s SET updated_at = ? WHERE session_id = ?`, s.sessionsTable)
	if _, err := tx.ExecContext(ctx, updateQuery, time.Now(), s.sessionID); err != nil {
		return fmt.Errorf("%w: %w", ErrOperationFailed, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: %w", ErrTransactionFailed, err)
	}
	return nil
}

// PopItem removes and returns the most recent item from the session.
// If no items exist, it returns nil without error.
// It uses SQLite's RETURNING clause to fetch the popped item.
func (s *SQLiteSession) PopItem(ctx context.Context) (*responses.ResponseInputItemUnionParam, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var messageData string
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = (
			SELECT id FROM %s
			WHERE session_id = ?
			ORDER BY created_at DESC
			LIMIT 1
		)
		RETURNING message_data`, s.messagesTable, s.messagesTable)

	err := s.db.QueryRowContext(ctx, query, s.sessionID).Scan(&messageData)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOperationFailed, err)
	}

	var item responses.ResponseInputItemUnionParam
	if err := json.Unmarshal([]byte(messageData), &item); err != nil {
		// Corrupted data; treat as no item.
		return nil, nil
	}

	return &item, nil
}

// ClearSession removes all items and the session metadata for this session.
func (s *SQLiteSession) ClearSession(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTransactionFailed, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Deleting from messages cascades due to foreign key.
	deleteMessages := fmt.Sprintf(`DELETE FROM %s WHERE session_id = ?`, s.messagesTable)
	if _, err := tx.ExecContext(ctx, deleteMessages, s.sessionID); err != nil {
		return fmt.Errorf("%w: %w", ErrOperationFailed, err)
	}

	deleteSession := fmt.Sprintf(`DELETE FROM %s WHERE session_id = ?`, s.sessionsTable)
	if _, err := tx.ExecContext(ctx, deleteSession, s.sessionID); err != nil {
		return fmt.Errorf("%w: %w", ErrOperationFailed, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: %w", ErrTransactionFailed, err)
	}
	return nil
}

// Close closes the underlying database connection.
// It should be called when the session is no longer needed.
func (s *SQLiteSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db == nil {
		return nil
	}
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("%w: %w", ErrDatabaseClose, err)
	}
	s.db = nil
	return nil
}
