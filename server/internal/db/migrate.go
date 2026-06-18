package db

import "database/sql"

func Migrate(database *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			nickname TEXT NOT NULL,
			role TEXT NOT NULL,
			avatar TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS souls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			birth_info TEXT NOT NULL,
			death_info TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL,
			merit_score INTEGER NOT NULL DEFAULT 0,
			karma_score INTEGER NOT NULL DEFAULT 0,
			memory_residue INTEGER NOT NULL DEFAULT 0,
			relationship_risk TEXT NOT NULL DEFAULT 'low',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS life_book_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			soul_id INTEGER NOT NULL,
			expected_lifespan INTEGER NOT NULL,
			actual_lifespan INTEGER NOT NULL,
			death_reason TEXT NOT NULL,
			fate_level TEXT NOT NULL,
			calamity_count INTEGER NOT NULL DEFAULT 0,
			locked INTEGER NOT NULL DEFAULT 0,
			risk_flag TEXT NOT NULL DEFAULT 'normal',
			updated_at TEXT NOT NULL,
			FOREIGN KEY (soul_id) REFERENCES souls(id)
		);`,
		`CREATE TABLE IF NOT EXISTS capture_tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			soul_id INTEGER NOT NULL,
			assignee_id INTEGER,
			status TEXT NOT NULL,
			location TEXT NOT NULL,
			scheduled_time TEXT NOT NULL,
			actual_time TEXT NOT NULL DEFAULT '',
			exception_type TEXT NOT NULL DEFAULT '',
			exception_note TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (soul_id) REFERENCES souls(id),
			FOREIGN KEY (assignee_id) REFERENCES users(id)
		);`,
		`CREATE TABLE IF NOT EXISTS reincarnation_orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			soul_id INTEGER NOT NULL,
			target_realm TEXT NOT NULL,
			priority INTEGER NOT NULL DEFAULT 0,
			soup_status TEXT NOT NULL,
			quota_type TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL,
			review_note TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (soul_id) REFERENCES souls(id)
		);`,
		`CREATE TABLE IF NOT EXISTS approvals (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL,
			target_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			status TEXT NOT NULL,
			applicant_id INTEGER,
			approver_id INTEGER,
			reason TEXT NOT NULL DEFAULT '',
			result_note TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (applicant_id) REFERENCES users(id),
			FOREIGN KEY (approver_id) REFERENCES users(id)
		);`,
		`CREATE TABLE IF NOT EXISTS hell_floors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			floor_no INTEGER NOT NULL UNIQUE,
			name TEXT NOT NULL,
			capacity INTEGER NOT NULL,
			occupancy INTEGER NOT NULL,
			equipment_status TEXT NOT NULL,
			load_level TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS hell_sentences (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			soul_id INTEGER NOT NULL,
			floor_id INTEGER NOT NULL,
			crime_type TEXT NOT NULL,
			sentence_days INTEGER NOT NULL,
			pain_level INTEGER NOT NULL,
			review_status TEXT NOT NULL,
			equipment_id TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (soul_id) REFERENCES souls(id),
			FOREIGN KEY (floor_id) REFERENCES hell_floors(id)
		);`,
		`CREATE TABLE IF NOT EXISTS wish_tickets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			wish_type TEXT NOT NULL,
			requester_name TEXT NOT NULL,
			incense_amount INTEGER NOT NULL DEFAULT 0,
			merit_score INTEGER NOT NULL DEFAULT 0,
			status TEXT NOT NULL,
			assigned_deity TEXT NOT NULL DEFAULT '',
			result_note TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS soup_inventory (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			stock INTEGER NOT NULL,
			warning_stock INTEGER NOT NULL,
			formula_note TEXT NOT NULL DEFAULT '',
			updated_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS soup_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			soul_id INTEGER NOT NULL,
			order_id INTEGER,
			inventory_id INTEGER NOT NULL,
			dose INTEGER NOT NULL,
			memory_after INTEGER NOT NULL,
			operator_id INTEGER,
			created_at TEXT NOT NULL,
			FOREIGN KEY (soul_id) REFERENCES souls(id),
			FOREIGN KEY (order_id) REFERENCES reincarnation_orders(id),
			FOREIGN KEY (inventory_id) REFERENCES soup_inventory(id),
			FOREIGN KEY (operator_id) REFERENCES users(id)
		);`,
		`CREATE TABLE IF NOT EXISTS alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			level TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			handled INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			actor_id INTEGER,
			action TEXT NOT NULL,
			target_type TEXT NOT NULL,
			target_id INTEGER NOT NULL,
			before_json TEXT NOT NULL DEFAULT '',
			after_json TEXT NOT NULL DEFAULT '',
			note TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			FOREIGN KEY (actor_id) REFERENCES users(id)
		);`,
	}

	for _, statement := range statements {
		if _, err := database.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}
