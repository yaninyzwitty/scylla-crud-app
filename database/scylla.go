package database

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
)

const (
	maxRetries = 60              // Max retry attempts
	retryDelay = 1 * time.Second // 1 second delay between retries
)

func NewDatabaseConnection(hosts string) (*gocqlx.Session, error) {
	// Separate hosts by commas
	cluster := gocql.NewCluster(hosts)
	cluster.Keyspace = "music"
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second

	var session gocqlx.Session
	var err error

	// Retry connection logic
	for retries := 0; retries < maxRetries; retries++ {
		session, err = gocqlx.WrapSession(cluster.CreateSession())
		if err == nil {
			fmt.Println("Successfully connected to the database.")
			break
		}
		fmt.Printf("Failed to connect to database (attempt %d/%d): %v\n", retries+1, maxRetries, err)
		time.Sleep(retryDelay)
	}

	// If we failed to connect after max retries, return the error
	if err != nil {
		return nil, fmt.Errorf("could not establish database connection after %d attempts: %w", maxRetries, err)
	}

	// Create keyspace if it doesn't exist
	err = session.ExecStmt(`CREATE KEYSPACE IF NOT EXISTS music WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`)
	if err != nil {
		return nil, fmt.Errorf("failed to create keyspace: %w", err)
	}

	// Create persons table if it doesn't exist
	err = session.ExecStmt(`CREATE TABLE IF NOT EXISTS songs (id UUID PRIMARY KEY , title TEXT, album TEXT, artist TEXT, tags SET<TEXT>, data BLOB );
`)
	if err != nil {
		return nil, fmt.Errorf("failed to create persons table: %w", err)
	}

	return &session, nil
}
