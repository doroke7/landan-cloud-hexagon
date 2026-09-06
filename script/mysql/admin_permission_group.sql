# 傾印（Dump）資料表 tx-admin_permission_groups
# ------------------------------------------------------------
# domain: internal/domain/admin_permission_group.go
#   Id uint64 / ParentId uint64 / Key string / Name string / CreatedAt / UpdatedAt / DeletedAt
#   gorm NamingStrategy：AdminPermissionGroup -> admin_permission_groups，前綴 CONFIG.MYSQL.PREFIX（tx-）
#   軟刪除以 deleted_at = '2038-01-19 03:14:07' 代表「未刪除」
#   parent_id = 0 代表頂層節點（虛擬 root 的直接子節點）

DROP TABLE IF EXISTS `tx-admin_permission_groups`;

CREATE TABLE `tx-admin_permission_groups` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '父節點 id，0 = 頂層',
  `key` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NOT NULL DEFAULT '2038-01-19 03:14:07',
  PRIMARY KEY (`id`),
  UNIQUE KEY `k-da` (`key`,`deleted_at`),
  KEY `p-da` (`parent_id`,`deleted_at`),
  KEY `n-da` (`name`,`deleted_at`),
  KEY `da` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `tx-admin_permission_groups` WRITE;
/*!40000 ALTER TABLE `tx-admin_permission_groups` DISABLE KEYS */;

INSERT INTO `tx-admin_permission_groups` (`id`, `parent_id`, `key`, `name`, `created_at`, `updated_at`, `deleted_at`)
VALUES
	-- 頂層分類
	(1,0,'game','遊戲','2026-01-01 09:00:00','2026-01-01 09:00:00','2038-01-19 03:14:07'),
	(2,0,'table','桌台','2026-01-01 09:00:00','2026-01-01 09:00:00','2038-01-19 03:14:07'),
	(3,0,'permission','權限','2026-01-01 09:00:00','2026-01-01 09:00:00','2038-01-19 03:14:07'),
	(4,0,'system','系統','2026-01-01 09:00:00','2026-01-01 09:00:00','2038-01-19 03:14:07'),

	-- game 底下
	(5,1,'game.type','遊戲類型','2026-01-01 09:00:00','2026-01-01 09:00:00','2038-01-19 03:14:07'),
	(6,1,'game.list','遊戲清單','2026-01-01 09:00:00','2026-01-01 09:00:00','2038-01-19 03:14:07'),

	-- table 底下
	(7,2,'table.list','桌台清單','2026-01-01 09:00:00','2026-01-01 09:00:00','2038-01-19 03:14:07'),
	(8,2,'table.record','桌台紀錄','2026-01-01 09:00:00','2026-01-01 09:00:00','2038-01-19 03:14:07'),

	-- permission 底下
	(9,3,'permission.admin_user','後台帳號','2026-01-01 09:00:00','2026-01-01 09:00:00','2038-01-19 03:14:07'),
	(10,3,'permission.admin_role','角色','2026-01-01 09:00:00','2026-01-01 09:00:00','2038-01-19 03:14:07'),
	(11,3,'permission.admin_permission','權限明細','2026-01-01 09:00:00','2026-01-01 09:00:00','2038-01-19 03:14:07'),
	(12,3,'permission.admin_permission_group','權限分類','2026-01-01 09:00:00','2026-01-01 09:00:00','2038-01-19 03:14:07'),

	-- system 底下
	(13,4,'system.config','參數設定','2026-01-01 09:00:00','2026-01-01 09:00:00','2038-01-19 03:14:07'),
	(14,4,'system.log','操作日誌','2026-01-01 09:00:00','2026-01-01 09:00:00','2038-01-19 03:14:07'),

	-- 巢狀第三層範例：game.type 底下再細分
	(15,5,'game.type.detail','遊戲類型-詳情','2026-01-01 09:00:00','2026-01-01 09:00:00','2038-01-19 03:14:07'),
	(16,5,'game.type.edit','遊戲類型-編輯','2026-01-01 09:00:00','2026-01-01 09:00:00','2038-01-19 03:14:07'),

	-- 軟刪除範例（deleted_at 非哨兵值）
	(17,4,'system.deprecated','已停用分類','2026-01-01 09:00:00','2026-06-01 09:00:00','2026-06-01 09:00:00');

/*!40000 ALTER TABLE `tx-admin_permission_groups` ENABLE KEYS */;
UNLOCK TABLES;
