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
			('mengpo', '孟婆', 'mengpo', '', ?),
			('juling', '巨灵神', 'heavenly_general', '', ?),
			('chenghuang', '城隍爷', 'city_god', '', ?),
			('xuntian', '巡天审计使', 'auditor', '', ?)`, []interface{}{now, now, now, now, now, now, now, now}},
		{`INSERT INTO souls (name, birth_info, death_info, status, merit_score, karma_score, memory_residue, relationship_risk, created_at, updated_at) VALUES
			('张三丰外卖版', '甲子年三月初三 午时', '夜跑误入鬼门关', 'pending_capture', 72, 18, 12, 'low', ?, ?),
			('李小卷', '乙丑年七月十五 子时', '熬夜写奏折', 'reincarnation_waiting', 91, 6, 25, 'medium', ?, ?),
			('王大胆', '丙寅年正月初九 酉时', '十八层体验券过期', 'hell', 21, 88, 4, 'high', ?, ?),
			('赵六六', '丁卯年二月初二 辰时', '阳寿异常延长', 'alive', 66, 9, 0, 'high', ?, ?),
			('钱多多', '戊辰年五月初五 巳时', '许愿过度反噬', 'pending_capture', 8, 45, 7, 'medium', ?, ?),
			('孙不睡', '己巳年六月十六 丑时', '连续加班三百年', 'reincarnation_waiting', 84, 12, 38, 'low', ?, ?),
			('周摸鱼', '庚午年八月廿八 申时', '工位灵魂出窍', 'captured', 49, 31, 19, 'medium', ?, ?),
			('吴插队', '辛未年九月初一 卯时', '试图插队投胎', 'reviewing', 17, 67, 11, 'high', ?, ?),
			('郑有礼', '壬申年十一月十一 亥时', '寿终正寝', 'reincarnation_waiting', 96, 2, 5, 'low', ?, ?),
			('冯造谣', '癸酉年腊月初八 未时', '造谣地府 KPI', 'hell', 14, 93, 3, 'low', ?, ?),
			('陈摸奖', '甲戌年二月十四 巳时', '连抽保底未出', 'pending_capture', 39, 42, 18, 'medium', ?, ?),
			('林迟到', '乙亥年三月廿二 辰时', '迟到撞上鬼门关打卡', 'pending_capture', 55, 20, 8, 'low', ?, ?),
			('黄加急', '丙子年四月初七 午时', '申请加急投胎失败', 'reincarnation_waiting', 63, 33, 16, 'medium', ?, ?),
			('何熬夜', '丁丑年五月十八 子时', '三界日报连夜赶稿', 'reincarnation_waiting', 77, 11, 29, 'low', ?, ?),
			('罗关系', '戊寅年六月初六 申时', '阳寿审批插队', 'alive', 22, 70, 2, 'high', ?, ?),
			('谢冲塔', '己卯年七月十九 酉时', '硬闯南天门直播', 'hell', 9, 84, 6, 'medium', ?, ?),
			('韩免汤', '庚辰年八月初八 卯时', '试图免喝孟婆汤', 'reviewing', 31, 48, 44, 'high', ?, ?),
			('朱排队', '辛巳年九月廿九 未时', '奈何桥排队睡过站', 'reincarnation_waiting', 68, 15, 21, 'low', ?, ?)`, []interface{}{now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now}},
		{`INSERT INTO life_book_records (soul_id, expected_lifespan, actual_lifespan, death_reason, fate_level, calamity_count, locked, risk_flag, updated_at) VALUES
			(1, 88, 73, '夜跑误入鬼门关', '小富即安', 3, 0, 'warning', ?),
			(2, 99, 99, '寿终正寝', '天选打工魂', 1, 0, 'normal', ?),
			(3, 66, 45, '主动试胆', '困难模式', 9, 1, 'critical', ?),
			(4, 69, 96, '疑似被加寿', '关系命', 2, 1, 'critical', ?),
			(5, 54, 54, '许愿过度反噬', '财来财去', 6, 0, 'warning', ?),
			(6, 82, 82, '连续加班三百年', '劳模命', 4, 0, 'normal', ?),
			(9, 101, 101, '寿终正寝', '积善之家', 0, 0, 'normal', ?),
			(11, 61, 47, '保底未出心态崩塌', '玄不救非', 8, 0, 'warning', ?),
			(12, 76, 75, '迟到撞上鬼门关打卡', '卡点命', 5, 0, 'warning', ?),
			(15, 62, 108, '阳寿审批插队', '关系命加强版', 2, 1, 'critical', ?),
			(16, 45, 30, '硬闯南天门直播', '作死命', 11, 1, 'critical', ?),
			(18, 88, 88, '奈何桥排队睡过站', '慢性子福报', 1, 0, 'normal', ?)`, []interface{}{now, now, now, now, now, now, now, now, now, now, now, now}},
		{`INSERT INTO capture_tasks (soul_id, assignee_id, status, location, scheduled_time, created_at, updated_at) VALUES
			(1, 4, 'pending', '人间西市烧烤摊', ?, ?, ?),
			(5, 4, 'running', '财神庙许愿池', ?, ?, ?),
			(2, 4, 'exception', '天庭加班区', ?, ?, ?),
			(7, 4, 'completed', '人间云办公间', ?, ?, ?),
			(8, 7, 'canceled', '鬼门关 VIP 通道', ?, ?, ?),
			(11, 4, 'pending', '人间抽卡广场 7 号机', ?, ?, ?),
			(12, 7, 'pending', '东岳庙后门考勤机', ?, ?, ?),
			(13, 4, 'running', '奈何桥加急窗口', ?, ?, ?),
			(14, 4, 'running', '三界日报排版间', ?, ?, ?),
			(17, 7, 'exception', '孟婆汤试饮台', ?, ?, ?),
			(15, 4, 'exception', '生死簿关系户窗口', ?, ?, ?),
			(16, 7, 'completed', '南天门直播通道', ?, ?, ?),
			(18, 4, 'completed', '奈何桥候车区', ?, ?, ?)`, []interface{}{now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now}},
		{`INSERT INTO reincarnation_orders (soul_id, target_realm, priority, soup_status, quota_type, status, review_note, created_at, updated_at) VALUES
			(2, '人间', 8, 'pending', '', 'queued', '', ?, ?),
			(1, '猫界', 6, 'issued', '普通胎', 'pending_soup', '建议保留方向感。', ?, ?),
			(6, '人间', 9, 'pending', '天才胎', 'reviewing', '记忆残留偏高，需复核。', ?, ?),
			(9, '天界', 10, 'issued', '福报胎', 'completed', '功德充足，已完成轮回。', ?, ?),
			(8, '人间', 1, 'pending', '', 'rejected', '插队证据确凿。', ?, ?),
			(7, '鱼界', 4, 'pending', '清闲胎', 'approved', '待分配胎位。', ?, ?),
			(13, '人间', 7, 'pending', '', 'queued', '一直催号，建议安排普通胎。', ?, ?),
			(14, '天界', 8, 'pending', '文书胎', 'reviewing', '疑似仍记得三界日报模板。', ?, ?),
			(18, '人间', 5, 'pending', '普通胎', 'pending_soup', '排队耐心好，优先保持平常心。', ?, ?),
			(11, '猫界', 3, 'issued', '橘猫胎', 'completed', '已完成，愿下世少抽卡。', ?, ?),
			(17, '人间', 2, 'failed', '', 'rejected', '拒喝孟婆汤，驳回重排。', ?, ?),
			(12, '人间', 6, 'pending', '准点胎', 'approved', '待孟婆汤窗口叫号。', ?, ?)`, []interface{}{now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now}},
		{`INSERT INTO approvals (type, target_id, title, status, applicant_id, approver_id, reason, result_note, created_at, updated_at) VALUES
			('lifespan_change', 1, '寿命变更申请', 'pending', 3, NULL, '救人一命，建议补阳寿。', '', ?, ?),
			('hell_review', 1, '地狱刑期复核', 'transferred', 2, NULL, '疑似判官手滑，多判三百年。', '转阎王二审', ?, ?),
			('reincarnation_jump', 5, '投胎插队申请', 'approved', 7, 1, '前世连续扶老人过马路。', '准予普通胎优先', ?, ?),
			('capture_exception', 3, '勾魂异常复盘', 'rejected', 4, 3, '疑似飞升导致任务异常。', '证据不足，继续蹲点', ?, ?),
			('soup_extra_dose', 2, '孟婆汤加量申请', 'pending', 5, NULL, '记忆残留 38%，还记得上家密码。', '', ?, ?),
			('capture_exception', 10, '魂魄逃逸复盘', 'pending', 7, NULL, '目标听到加浓孟婆汤后翻栏逃逸，需判定是否加派白无常。', '', ?, ?),
			('lifespan_change', 10, '阳寿异常延长', 'pending', 8, NULL, '罗关系阳寿多出 46 年，疑似香火充值走错账。', '', ?, ?),
			('hell_review', 3, '孽镜地狱刑期复核', 'transferred', 2, NULL, '朋友圈谣言日更是否按连载罪加罚，需阎王二审。', '转阎王复核', ?, ?),
			('reincarnation_jump', 7, '猫界胎位加塞', 'approved', 7, 1, '申请人称前世救过三只橘猫。', '准予猫界普通胎', ?, ?),
			('soup_extra_dose', 3, '孟婆汤加浓审批', 'approved', 5, 2, '记忆残留 44%，仍记得上辈子 Wi-Fi 密码。', '加浓两勺', ?, ?),
			('lifespan_change', 8, '抽卡阳寿补偿', 'rejected', 3, 1, '陈摸奖申请用十连保底抵阳寿。', '天条不认抽卡账单', ?, ?),
			('capture_exception', 11, '阳间抢救干扰', 'rejected', 4, 3, '阳间亲友烧香过猛，影响勾魂窗口。', '继续排队，不许插队', ?, ?)`, []interface{}{now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now, now}},
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
