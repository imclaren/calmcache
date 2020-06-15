package dbcache

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/imclaren/sqldb"
	"github.com/imclaren/sqldb/sqlite"
)

// DB is the cache sql database struct
type DB struct {
	sync.RWMutex
	*sqldb.DB
}

// Init opens the cache sql database and creates the database tables if they do not already exist
func Init(DBPath string, ctx context.Context, cancel context.CancelFunc) (DB, error) {
	DB, err := Open(DBPath, ctx, cancel)
	if err != nil {
		return DB, err
	}
	err = DB.CreateTable()
	return DB, err
}

// Open opens the cache sql database 
func Open(DBPath string, ctx context.Context, cancel context.CancelFunc) (DB, error) {
	connectString := sqlite.ConnectString(DBPath, "UTC")
	return initDB(ctx, cancel, "sqlite", connectString)
}

func initDB(ctx context.Context, cancelFunc context.CancelFunc, dbType, connectString string) (DB, error) {
	db, err := sqldb.Init(ctx, cancelFunc, dbType, connectString)
	if err != nil {
		return DB{}, err
	}
	return DB{
		sync.RWMutex{}, &db,
	}, nil
}

// CreateTable creates the database table
func (db *DB) CreateTable() error {
	db.Lock()
	defer db.Unlock()

	switch db.Type {
	case "sqlite":
		_, err := db.Exec(`
    		CREATE TABLE IF NOT EXISTS cache (
	    		id INTEGER PRIMARY KEY,
	    		bucket TEXT,
	    		key TEXT,
	    		size INT,
	    		access_count INT,
	    		expires_at TIMESTAMP,
				created_at TIMESTAMP NULL DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
			    updated_at TIMESTAMP NULL DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW'))	    		
    			)
    		`)
		if err != nil {
			return err
		}
		// Create updated_at trigger for cache table
		_, err = db.Exec(`
			CREATE TRIGGER [update_cache_updated_at]
			    AFTER UPDATE
			    ON cache
			BEGIN
			    UPDATE cache SET updated_at=STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW') WHERE id=NEW.id;
			END;
		`)
		if err != nil && err.Error() != "trigger [update_cache_updated_at] already exists" {
			return err
		}
	case "postgres":
		_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS cache (
				id BIGSERIAL PRIMARY KEY, 
				bucket TEXT,
				key TEXT, 
				size BIGINT, 
				access_count BIGINT, 
				expires_at TIMESTAMP, 
				created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
			    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
			)
		`)
		if err != nil {
			return err
		}
		// Create or replace update_updated_at_column function
		// Note we only need to do this once for all of the tables that we update
		_, err = db.Exec(`
			CREATE OR REPLACE FUNCTION update_updated_at_column()
			RETURNS TRIGGER AS $$
			BEGIN
			   NEW.updated_at = now(); 
			   RETURN NEW;
			END;
			$$ language 'plpgsql';
		`)
		if err != nil {
			return err
		}
		// Create updated_at trigger for cache table
		_, err = db.Exec(`
			CREATE TRIGGER update_cache_updated_at 
				BEFORE UPDATE
				ON cache 
				FOR EACH ROW 
				EXECUTE PROCEDURE update_updated_at_column();
		`)
		if err != nil && err.Error() != "pq: trigger \"update_cacheupdated_at\" for relation \"cache\" already exists" {
			return err
		}
	default:
		return fmt.Errorf("Create table error: database type not implemented: %s", db.Type)
	}

	// Add indexes
	indexSlice := []struct {
		col  string
		isUnique  bool
	}{
		{"key", true},
		{"size", false},
		{"access_count", false},
		//{"expires_at", false},
		//{"created_at", false},
		{"updated_at", false},
	}
	for _, in := range indexSlice {
		SQLString, err := indexSQLString(db.Type, "btree", in.isUnique, "cache", []string{in.col}, "")
		if err != nil {
			return err
		}
		_, err = db.Exec(SQLString)
		if err != nil {
			return err
		}
	}
	return nil
}

func indexSQLString(dbType, indexType string, isUnique bool, tableName string, indexColumns []string, whereString string) (string, error) {
	uniqueString := ""
	if isUnique {
		uniqueString = "UNIQUE"
	}
	uColumns := strings.Join(indexColumns, "_")
	commaColumns := strings.Join(indexColumns, ", ")
	switch dbType {
	case "sqlite":
		return fmt.Sprintf("CREATE %s INDEX IF NOT EXISTS %s_%s_idx ON %s (%s) %s", uniqueString, tableName, uColumns, tableName, commaColumns, whereString), nil
	case "postgres":
		return fmt.Sprintf("CREATE %s INDEX IF NOT EXISTS %s_%s_idx ON %s USING %s (%s) %s", uniqueString, tableName, uColumns, tableName, indexType, commaColumns, whereString), nil
	default:
		return "", fmt.Errorf("Create table index error: database type not implemented: %s", dbType)
	}
} 

// DropTable drops the database table
func (db *DB) DropTable() error {
	db.Lock()
	defer db.Unlock()

	_, err := db.Exec("DROP TABLE IF EXISTS cache")
	return err
}
