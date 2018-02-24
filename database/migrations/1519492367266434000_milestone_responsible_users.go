package migrations

import "gogit/database"

// Migrate the database to a new version
func (Migration) MigrateMilestoneResponsibleUsers() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.milestone_responsible_users (
    id bigserial NOT NULL,
    milestone_id bigint NOT NULL,
    user_id bigint NOT NULL,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT milestone_responsible_users_pkey PRIMARY KEY (id),
    CONSTRAINT milestone_responsible_users_milestone_id_milestones_id FOREIGN KEY (milestone_id)
        REFERENCES public.milestones (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT milestone_responsible_users_user_id_users_id FOREIGN KEY (user_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT milestone_responsible_users_created_by_id_users_id FOREIGN KEY (user_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackMilestoneResponsibleUsers() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`DROP TABLE IF EXISTS public.milestone_responsible_users;`)

	return
}