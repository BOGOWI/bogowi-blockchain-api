package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the global database instance
type DB struct {
	conn *sql.DB
	mu   sync.RWMutex
}

// NFTMapping represents the mapping between blockchain token ID and Datakyte NFT ID
type NFTMapping struct {
	ID            int64
	TokenID       uint64
	DatakyteNFTID string
	Network       string
	ContractAddr  string
	OwnerAddress  string
	BookingID     string
	EventID       string
	Status        string
	MetadataURI   string
	ImageURL      string
	TxHash        string
	MintedAt      string
	RedeemedAt    *string
}

var (
	instance *DB
	once     sync.Once
)

// GetDB returns the singleton database instance
func GetDB() *DB {
	once.Do(func() {
		var err error
		instance, err = NewDB("")
		if err != nil {
			panic(fmt.Sprintf("failed to initialize database: %v", err))
		}
	})
	return instance
}

// NewDB creates a new database connection
func NewDB(dbPath string) (*DB, error) {
	if dbPath == "" {
		// Default path - in data directory relative to binary
		homeDir, _ := os.UserHomeDir()
		dbPath = filepath.Join(homeDir, ".bogowi", "nft_mappings.db")

		// Create directory if it doesn't exist
		dbDir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// Open database connection
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	conn.SetMaxOpenConns(1) // SQLite doesn't handle concurrent writes well
	conn.SetMaxIdleConns(1)

	db := &DB{conn: conn}

	// Initialize schema
	if err := db.initSchema(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

// initSchema creates the necessary tables
func (db *DB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS nft_token_mappings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		token_id INTEGER NOT NULL,
		datakyte_nft_id TEXT NOT NULL,
		network TEXT NOT NULL,
		contract_address TEXT NOT NULL,
		owner_address TEXT NOT NULL,
		booking_id TEXT,
		event_id TEXT,
		status TEXT DEFAULT 'active',
		metadata_uri TEXT,
		image_url TEXT,
		tx_hash TEXT NOT NULL,
		minted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		redeemed_at DATETIME,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(token_id, network, contract_address)
	);

	CREATE INDEX IF NOT EXISTS idx_datakyte_id ON nft_token_mappings(datakyte_nft_id);
	CREATE INDEX IF NOT EXISTS idx_owner ON nft_token_mappings(owner_address);
	CREATE INDEX IF NOT EXISTS idx_status ON nft_token_mappings(status);
	CREATE INDEX IF NOT EXISTS idx_network ON nft_token_mappings(network);
	`

	_, err := db.conn.Exec(schema)
	return err
}

// SaveNFTMapping stores the mapping between token ID and Datakyte NFT ID
func (db *DB) SaveNFTMapping(mapping *NFTMapping) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	query := `
	INSERT INTO nft_token_mappings (
		token_id, datakyte_nft_id, network, contract_address, 
		owner_address, booking_id, event_id, status, 
		metadata_uri, image_url, tx_hash
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(token_id, network, contract_address) 
	DO UPDATE SET
		datakyte_nft_id = excluded.datakyte_nft_id,
		owner_address = excluded.owner_address,
		status = excluded.status,
		metadata_uri = excluded.metadata_uri,
		image_url = excluded.image_url,
		updated_at = CURRENT_TIMESTAMP
	`

	_, err := db.conn.Exec(query,
		mapping.TokenID,
		mapping.DatakyteNFTID,
		mapping.Network,
		mapping.ContractAddr,
		mapping.OwnerAddress,
		mapping.BookingID,
		mapping.EventID,
		mapping.Status,
		mapping.MetadataURI,
		mapping.ImageURL,
		mapping.TxHash,
	)

	return err
}

// GetDatakyteID retrieves the Datakyte NFT ID for a given token
func (db *DB) GetDatakyteID(tokenID uint64, network string) (string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var datakyteID string
	query := `
	SELECT datakyte_nft_id 
	FROM nft_token_mappings 
	WHERE token_id = ? AND network = ?
	LIMIT 1
	`

	err := db.conn.QueryRow(query, tokenID, network).Scan(&datakyteID)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("no mapping found for token %d on %s", tokenID, network)
	}
	return datakyteID, err
}

// GetNFTMapping retrieves the full mapping for a token
func (db *DB) GetNFTMapping(tokenID uint64, network string) (*NFTMapping, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	mapping := &NFTMapping{}
	query := `
	SELECT id, token_id, datakyte_nft_id, network, contract_address,
		   owner_address, booking_id, event_id, status, metadata_uri,
		   image_url, tx_hash, minted_at, redeemed_at
	FROM nft_token_mappings
	WHERE token_id = ? AND network = ?
	LIMIT 1
	`

	err := db.conn.QueryRow(query, tokenID, network).Scan(
		&mapping.ID,
		&mapping.TokenID,
		&mapping.DatakyteNFTID,
		&mapping.Network,
		&mapping.ContractAddr,
		&mapping.OwnerAddress,
		&mapping.BookingID,
		&mapping.EventID,
		&mapping.Status,
		&mapping.MetadataURI,
		&mapping.ImageURL,
		&mapping.TxHash,
		&mapping.MintedAt,
		&mapping.RedeemedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no mapping found for token %d on %s", tokenID, network)
	}
	return mapping, err
}

// UpdateNFTStatus updates the status of an NFT
func (db *DB) UpdateNFTStatus(tokenID uint64, network string, status string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	query := `
	UPDATE nft_token_mappings 
	SET status = ?, updated_at = CURRENT_TIMESTAMP
	WHERE token_id = ? AND network = ?
	`

	result, err := db.conn.Exec(query, status, tokenID, network)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no mapping found for token %d on %s", tokenID, network)
	}

	return nil
}

// UpdateNFTRedemption marks an NFT as redeemed
func (db *DB) UpdateNFTRedemption(tokenID uint64, network string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	query := `
	UPDATE nft_token_mappings 
	SET status = 'redeemed', 
		redeemed_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
	WHERE token_id = ? AND network = ?
	`

	result, err := db.conn.Exec(query, tokenID, network)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no mapping found for token %d on %s", tokenID, network)
	}

	return nil
}

// GetUserNFTs retrieves all NFTs owned by a specific address
func (db *DB) GetUserNFTs(ownerAddress string, network string) ([]NFTMapping, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	query := `
	SELECT id, token_id, datakyte_nft_id, network, contract_address,
		   owner_address, booking_id, event_id, status, metadata_uri,
		   image_url, tx_hash, minted_at, redeemed_at
	FROM nft_token_mappings
	WHERE owner_address = ? AND network = ?
	ORDER BY minted_at DESC
	`

	rows, err := db.conn.Query(query, ownerAddress, network)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []NFTMapping
	for rows.Next() {
		var m NFTMapping
		err := rows.Scan(
			&m.ID,
			&m.TokenID,
			&m.DatakyteNFTID,
			&m.Network,
			&m.ContractAddr,
			&m.OwnerAddress,
			&m.BookingID,
			&m.EventID,
			&m.Status,
			&m.MetadataURI,
			&m.ImageURL,
			&m.TxHash,
			&m.MintedAt,
			&m.RedeemedAt,
		)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, m)
	}

	return mappings, rows.Err()
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}
