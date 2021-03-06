package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateIssueStatusesTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec(`CREATE TABLE IF NOT EXISTS public.issue_statuses(
    id serial NOT NULL,
    title varchar(120) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
    )`).Error

	return
}

// Rollback the database to previous version
func (Migration) RollbackIssueStatusesTable() (e error) {
	db := database.NewGORMInstance()
	defer db.Close()
	e = db.Exec("DROP TABLE IF EXISTS public.issue_statuses;").Error

	return
}
