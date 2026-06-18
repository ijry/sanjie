package db

import (
	"database/sql"
	"time"
)

func Seed(database *sql.DB) error {
	var count int
	if err := database.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	now := time.Now().Format(time.RFC3339)
	tx, err := database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statements := []struct {
		query string
		args  []interface{}
	}{
		{`INSERT INTO users (username, nickname, role, avatar, created_at) VALUES
			('yudi', '玉帝', 'jade_emperor', '', ?),
			('yanwang', '阎王', 'yanwang', '', ?),
			('panguan', '判官', 'judge', '', ?),
			('heiwu', '黑无常', 'impermanence', '', ?),
			('mengpo', '孟婆', 'mengpo', '', ?)`, []interface{}{now, now, now, now, now}},
		{`INSERT INTO souls (name, birth_info, death_info, status, merit_score, karma_score, memory_residue, relationship_risk, created_at, updated_at) VALUES
			('张三', '甲子年三月初八 午时', '阳寿将尽', 'pending_capture', 12, 37, 0, 'low', ?, ?),
			('李四', '乙丑年六月十二 子时', '待投胎', 'reincarnation_waiting', 88, 4, 12, 'medium', ?, ?),
			('王五', '丙寅年九月廿一 酉时', '审判完成', 'hell', 3, 91, 0, 'low', ?, ?),
			('赵六', '丁卯年正月十五 辰时', '阳寿异常', 'alive', 66, 6, 0, 'high', ?, ?)`, []interface{}{now, now, now, now, now, now, now, now}},
		{`INSERT INTO life_book_records (soul_id, expected_lifespan, actual_lifespan, death_reason, fate_level, calamity_count, locked, risk_flag, updated_at) VALUES
			(1, 73, 73, '正常到期', '普通命', 3, 0, 'normal', ?),
			(2, 81, 81, '功德圆满', '富贵命', 1, 0, 'normal', ?),
			(3, 58, 58, '恶业反噬', '坎坷命', 7, 0, 'warning', ?),
			(4, 69, 96, '疑似被加寿', '关系命', 2, 1, 'critical', ?)`, []interface{}{now, now, now, now}},
		{`INSERT INTO capture_tasks (soul_id, assignee_id, status, location, scheduled_time, created_at, updated_at) VALUES
			(1, 4, 'pending', '人间东市第三巷', ?, ?, ?)`, []interface{}{now, now, now}},
		{`INSERT INTO reincarnation_orders (soul_id, target_realm, priority, soup_status, quota_type, status, review_note, created_at, updated_at) VALUES
			(2, 'human', 88, 'pending', '普通胎', 'queued', '功德较高但记忆残留偏高', ?, ?)`, []interface{}{now, now}},
		{`INSERT INTO hell_floors (floor_no, name, capacity, occupancy, equipment_status, load_level, updated_at) VALUES
			(1, '拔舌地狱', 1000, 940, 'warning', 'critical', ?),
			(2, '剪刀地狱', 900, 421, 'normal', 'normal', ?),
			(3, '铁树地狱', 800, 688, 'normal', 'warning', ?)`, []interface{}{now, now, now}},
		{`INSERT INTO hell_sentences (soul_id, floor_id, crime_type, sentence_days, pain_level, review_status, equipment_id, created_at, updated_at) VALUES
			(3, 1, '口舌是非', 360, 7, 'none', 'tongue-001', ?, ?)`, []interface{}{now, now}},
		{`INSERT INTO wish_tickets (title, wish_type, requester_name, incense_amount, merit_score, status, assigned_deity, created_at, updated_at) VALUES
			('求暴富但不想上班', 'wealth', '刘某', 18, 1, 'pending', '', ?, ?),
			('求考试刚好及格', 'exam', '陈某', 6, 12, 'routed', '文昌帝君', ?, ?)`, []interface{}{now, now, now, now}},
		{`INSERT INTO soup_inventory (name, stock, warning_stock, formula_note, updated_at) VALUES
			('标准孟婆汤', 320, 80, '适合普通投胎批次', ?),
			('加浓孟婆汤', 48, 60, '适合记忆残留较高魂魄', ?)`, []interface{}{now, now}},
		{`INSERT INTO alerts (level, title, content, handled, created_at) VALUES
			('critical', '拔舌地狱容量过高', '当前容量达到 94%，建议开启跨层分流。', 0, ?),
			('warning', '加浓孟婆汤库存不足', '库存低于预警线，下一批投胎可能受影响。', 0, ?),
			('critical', '生死簿高危变更', '赵六阳寿记录已冻结，疑似关系户操作。', 0, ?)`, []interface{}{now, now, now}},
	}

	for _, statement := range statements {
		if _, err := tx.Exec(statement.query, statement.args...); err != nil {
			return err
		}
	}
	return tx.Commit()
}
