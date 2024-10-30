package tester

import "database/sql"

type Tester struct{}

func (t *Tester) NewDBConn(schemaName string) *sql.DB {
	// @@TODO: implement this shit
	return nil
}
