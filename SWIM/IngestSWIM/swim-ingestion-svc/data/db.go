package data

import (
	"database/sql"
	"fmt"
	"github.build.ge.com/aviation-intelligent-airport/configuration-manager-svc/config"
	"github.com/jackmanlabs/errors"
	_ "github.com/lib/pq"
)

func db() (*sql.DB, error) {

	database := config.ActiveConfig.PostgreSQL.Database
	host := config.ActiveConfig.PostgreSQL.Host
	port := config.ActiveConfig.PostgreSQL.Port
	password := config.ActiveConfig.PostgreSQL.Password
	username := config.ActiveConfig.PostgreSQL.Username

	// SSL mode is not used on Predix.
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, errors.Stack(err)
	}

	return db, nil
}

// I'd prefer not to have this.
// In the future, let's look at Dwayne's method of updating the tables dynamically.
// Or hire a DBA that's willing to work around the Predix DB administration obstacles.
func Recreate() error {
	db, err := db()
	if err != nil {
		return errors.Stack(err)
	}

	q := `
-- Insert DB setup SQL here.
`

	_, err = db.Exec(q)
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}
