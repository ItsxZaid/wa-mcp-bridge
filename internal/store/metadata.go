package store

import "database/sql"

type sqliteMetadataStore struct {
	db *sql.DB
}

type MetadataStore interface {
	GetValueByKey (key string) (*MetaData, error)
	SetValueByKey (key string, value []byte) error
	DeleteByKey (key string) error
}

type MetaData struct {
	ID        int
	Key string
	Value  []byte
}

func (m *sqliteMetadataStore) GetValueByKey(key string) (*MetaData, error) {
	return nil, sql.ErrNoRows
}

func (m *sqliteMetadataStore) SetValueByKey(key string, value []byte) error {						
	return nil
}

func (m *sqliteMetadataStore) DeleteByKey(key string) error	 {
	return nil
}

// TODOSSSSSSSSSSSSSSSSSSS