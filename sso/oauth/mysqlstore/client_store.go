package mysqlstore

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

type ClientStore struct {
	db        *gorm.DB
	tableName string

	initTableDisabled bool
	maxLifetime       time.Duration
	maxOpenConns      int
	maxIdleConns      int
}

// ClientStoreItem data item
type ClientStoreItem struct {
	ID     string `db:"id"`
	Secret string `db:"secret"`
	Domain string `db:"domain"`
	Data   string `db:"data"`
}

// NewClientStore creates PostgreSQL store instance
func NewClientStore(db *gorm.DB, options ...ClientStoreOption) (*ClientStore, error) {

	store := &ClientStore{
		db:           db,
		tableName:    "oauth2_clients",
		maxLifetime:  time.Hour * 2,
		maxOpenConns: 50,
		maxIdleConns: 25,
	}

	for _, o := range options {
		o(store)
	}

	var err error
	if !store.initTableDisabled {
		err = store.initTable()
	}

	if err != nil {
		return store, err
	}

	// store.db.SetMaxOpenConns(store.maxOpenConns)
	// store.db.SetMaxIdleConns(store.maxIdleConns)
	// store.db.SetConnMaxLifetime(store.maxLifetime)

	return store, err
}

func (s *ClientStore) initTable() error {

	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id VARCHAR(255) NOT NULL PRIMARY KEY,
		secret VARCHAR(255) NOT NULL,
		domain VARCHAR(255) NOT NULL,
		data TEXT NOT NULL	
	  );
`, s.tableName)

	res := s.db.Exec(query)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *ClientStore) toClientInfo(data string) (oauth2.ClientInfo, error) {
	var cm models.Client
	err := jsoniter.Unmarshal([]byte(data), &cm)
	return &cm, err
}
func (s *ClientStore) Set(ctx context.Context, id string, client *models.Client) error {
	return s.Create(ctx, client)
}

// GetByID retrieves and returns client information by id
func (s *ClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	if id == "" {
		return nil, nil
	}

	var item ClientStoreItem
	res := s.db.Raw(fmt.Sprintf("SELECT * FROM %s WHERE id = ?", s.tableName), id).Scan(&item)
	switch {
	case res.Error == sql.ErrNoRows:
		return nil, nil
	case res.Error != nil:
		return nil, res.Error
	}

	return s.toClientInfo(item.Data)
}

// Create creates and stores the new client information
func (s *ClientStore) Create(ctx context.Context, info oauth2.ClientInfo) error {
	data, err := jsoniter.Marshal(info)
	if err != nil {
		return err
	}
	INSERT_QUERY := fmt.Sprintf("INSERT INTO %s (id, secret, domain, data) VALUES (?,?,?,?)", s.tableName)

	res := s.db.Exec(INSERT_QUERY,
		info.GetID(),
		info.GetSecret(),
		info.GetDomain(),
		string(data))
	if res.Error != nil {
		return res.Error
	}
	return nil
}
