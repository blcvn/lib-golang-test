package mysqlstore

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

type TokenStore struct {
	db                *gorm.DB
	tableName         string
	gcDisabled        bool
	gcInterval        time.Duration
	ticker            *time.Ticker
	initTableDisabled bool
	maxLifetime       time.Duration
	maxOpenConns      int
	maxIdleConns      int
}

// TokenStoreItem data item
type TokenStoreItem struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	ExpiredAt time.Time `db:"expired_at"`
	Code      string    `db:"code"`
	Access    string    `db:"access"`
	Refresh   string    `db:"refresh"`
	Data      string    `db:"data"`
}

// NewTokenStore creates PostgreSQL store instance
func NewTokenStore(db *gorm.DB, options ...TokenStoreOption) (*TokenStore, error) {

	store := &TokenStore{
		db:           db,
		tableName:    "oauth2_tokens",
		gcInterval:   10 * time.Minute,
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

	if !store.gcDisabled {
		store.ticker = time.NewTicker(store.gcInterval)
		go store.gc()
	}

	// store.db.SetMaxOpenConns(store.maxOpenConns)
	// store.db.SetMaxIdleConns(store.maxIdleConns)
	// store.db.SetConnMaxLifetime(store.maxLifetime)

	return store, err
}

// Close close the store
func (s *TokenStore) Close() error {
	if !s.gcDisabled {
		s.ticker.Stop()
	}
	return nil
}

func (s *TokenStore) gc() {
	for range s.ticker.C {
		s.clean()
	}
}

func (s *TokenStore) initTable() error {

	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		code VARCHAR(255),
		access VARCHAR(255) NOT NULL,
		refresh VARCHAR(255) NOT NULL,
		data TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		KEY access_k(access),
		KEY refresh_k (refresh),
		KEY expired_at_k (expired_at),
		KEY code_k (code)
	  );
`, s.tableName)

	res := s.db.Exec(query)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (s *TokenStore) clean() {

	now := time.Now().Unix()
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE expired_at<=? OR (code='' AND access='' AND refresh='')", s.tableName)

	var count int64
	res := s.db.Raw(query, now).Row().Scan(&count)
	if res.Error != nil || count == 0 {
		if res.Error != nil {
			log.Println(res.Error())
		}
		return
	}

	res1 := s.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE expired_at<=? OR (code='' AND access='' AND refresh='')", s.tableName), now)
	if res1.Error != nil {
		log.Println(res1.Error.Error())
	}
}

// Create create and store the new token information
func (s *TokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	buf, _ := jsoniter.Marshal(info)
	item := &TokenStoreItem{
		Data:      string(buf),
		CreatedAt: time.Now(),
	}

	if code := info.GetCode(); code != "" {
		item.Code = code
		item.ExpiredAt = info.GetCodeCreateAt().Add(info.GetCodeExpiresIn())
	} else {
		item.Access = info.GetAccess()
		item.ExpiredAt = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn())

		if refresh := info.GetRefresh(); refresh != "" {
			item.Refresh = info.GetRefresh()
			item.ExpiredAt = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn())
		}
	}

	fmt.Print(item.CreatedAt)

	res := s.db.Exec(
		fmt.Sprintf("INSERT INTO %s (created_at, expired_at, code, access, refresh, data) VALUES (?,?,?,?,?,?)", s.tableName),
		item.CreatedAt,
		item.ExpiredAt,
		item.Code,
		item.Access,
		item.Refresh,
		item.Data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// RemoveByCode delete the authorization code
func (s *TokenStore) RemoveByCode(ctx context.Context, code string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE code=? LIMIT 1", s.tableName)
	res := s.db.Exec(query, code)
	if res.Error != nil && res.Error == sql.ErrNoRows {
		return nil
	}
	return res.Error
}

// RemoveByAccess use the access token to delete the token information
func (s *TokenStore) RemoveByAccess(ctx context.Context, access string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE access=? LIMIT 1", s.tableName)
	res := s.db.Exec(query, access)
	if res.Error != nil && res.Error == sql.ErrNoRows {
		return nil
	}
	return res.Error
}

// RemoveByRefresh use the refresh token to delete the token information
func (s *TokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE refresh=? LIMIT 1", s.tableName)
	res := s.db.Exec(query, refresh)
	if res.Error != nil && res.Error == sql.ErrNoRows {
		return nil
	}
	return res.Error
}

func (s *TokenStore) toTokenInfo(data string) oauth2.TokenInfo {
	var tm models.Token
	jsoniter.Unmarshal([]byte(data), &tm)
	return &tm
}

// GetByCode use the authorization code for token information data
func (s *TokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	if code == "" {
		return nil, nil
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE code=? LIMIT 1", s.tableName)
	var item TokenStoreItem
	res := s.db.Raw(query, code).Scan(&item)
	switch {
	case res.Error == sql.ErrNoRows:
		return nil, nil
	case res.Error != nil:
		return nil, res.Error
	}

	return s.toTokenInfo(item.Data), nil
}

// GetByAccess use the access token for token information data
func (s *TokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	if access == "" {
		return nil, nil
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE access=? LIMIT 1", s.tableName)
	var item TokenStoreItem
	res := s.db.Raw(query, access).Scan(&item)
	switch {
	case res.Error == sql.ErrNoRows:
		return nil, nil
	case res.Error != nil:
		return nil, res.Error
	}
	return s.toTokenInfo(item.Data), nil
}

// GetByRefresh use the refresh token for token information data
func (s *TokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	if refresh == "" {
		return nil, nil
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE refresh=? LIMIT 1", s.tableName)
	var item TokenStoreItem
	res := s.db.Raw(query, refresh).Scan(&item)
	switch {
	case res.Error == sql.ErrNoRows:
		return nil, nil
	case res.Error != nil:
		return nil, res.Error
	}
	return s.toTokenInfo(item.Data), nil
}
