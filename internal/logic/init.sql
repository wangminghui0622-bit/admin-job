/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

CREATE TABLE IF NOT EXISTS `sys_tasks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL COMMENT 'Create Time | 创建日期',
  `updated_at` timestamp NOT NULL COMMENT 'Update Time | 修改日期',
  `status` tinyint unsigned DEFAULT '1' COMMENT 'Status 1: normal 2: ban | 状态 1 正常 2 禁用',
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'Task Name | 任务名称',
  `task_group` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'Task Group | 任务分组',
  `cron_expression` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'Cron expression | 定时任务表达式',
  `pattern` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'Cron Pattern | 任务的模式 （用于区分和确定要执行的任务）',
  `payload` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT 'The data used in cron (JSON string) | 任务需要的数据(JSON 字符串)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `task_pattern` (`pattern`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Task Configuration Table | 任务配置表';

CREATE TABLE IF NOT EXISTS `sys_task_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `started_at` timestamp NOT NULL COMMENT 'Task Started Time | 任务启动时间',
  `finished_at` timestamp NOT NULL COMMENT 'Task Finished Time | 任务完成时间',
  `result` tinyint unsigned NOT NULL COMMENT 'The Task Process Result | 任务执行结果',
  `task_task_logs` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='Task Log Table | 任务日志表';

INSERT IGNORE INTO `sys_tasks`
  (`created_at`, `updated_at`, `status`, `name`, `task_group`, `cron_expression`, `pattern`, `payload`)
VALUES
  (NOW(), NOW(), 1, 'hello_world', 'base', '@every 60s', 'hello_world', '{"name": "Mike (DPTask 60s)"}');

INSERT IGNORE INTO `sys_apis`
  (`id`, `created_at`, `updated_at`, `path`, `description`, `api_group`, `service_name`, `method`, `is_required`)
VALUES
  (3001, NOW(), NOW(), '/sys-api/task/create', 'apiDesc.createTask', 'task', 'Job', 'POST', 0),
  (3002, NOW(), NOW(), '/sys-api/task/update', 'apiDesc.updateTask', 'task', 'Job', 'POST', 0),
  (3003, NOW(), NOW(), '/sys-api/task/delete', 'apiDesc.deleteTask', 'task', 'Job', 'POST', 0),
  (3004, NOW(), NOW(), '/sys-api/task/list', 'apiDesc.getTaskList', 'task', 'Job', 'POST', 0),
  (3005, NOW(), NOW(), '/sys-api/task', 'apiDesc.getTaskById', 'task', 'Job', 'POST', 0),
  (3006, NOW(), NOW(), '/sys-api/task_log/create', 'apiDesc.createTaskLog', 'task_log', 'Job', 'POST', 0),
  (3007, NOW(), NOW(), '/sys-api/task_log/update', 'apiDesc.updateTaskLog', 'task_log', 'Job', 'POST', 0),
  (3008, NOW(), NOW(), '/sys-api/task_log/delete', 'apiDesc.deleteTaskLog', 'task_log', 'Job', 'POST', 0),
  (3009, NOW(), NOW(), '/sys-api/task_log/list', 'apiDesc.getTaskLogList', 'task_log', 'Job', 'POST', 0),
  (3010, NOW(), NOW(), '/sys-api/task_log', 'apiDesc.getTaskLogById', 'task_log', 'Job', 'POST', 0);

INSERT IGNORE INTO `casbin_rules`
  (`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`)
VALUES
  (8001, 'p', '001', '/sys-api/task/create', 'POST', '', '', ''),
  (8002, 'p', '001', '/sys-api/task/update', 'POST', '', '', ''),
  (8003, 'p', '001', '/sys-api/task/delete', 'POST', '', '', ''),
  (8004, 'p', '001', '/sys-api/task/list', 'POST', '', '', ''),
  (8005, 'p', '001', '/sys-api/task', 'POST', '', '', ''),
  (8006, 'p', '001', '/sys-api/task_log/create', 'POST', '', '', ''),
  (8007, 'p', '001', '/sys-api/task_log/update', 'POST', '', '', ''),
  (8008, 'p', '001', '/sys-api/task_log/delete', 'POST', '', '', ''),
  (8009, 'p', '001', '/sys-api/task_log/list', 'POST', '', '', ''),
  (8010, 'p', '001', '/sys-api/task_log', 'POST', '', '', '');



-- menu registration: task management (parent = core SystemManagement menu id 3)
INSERT IGNORE INTO `sys_menus` VALUES
(210,'2026-10-15 06:49:53','2026-10-15 06:49:53',8,2,1,'/task','TaskManagement','','/sys/task/index',0,'Job',NULL,'route.taskManagement','ic:baseline-access-alarm',1,0,0,0,'',0,0,0,20,'',3);

-- role 1 is authorized for the task management menu (resolved by path)
INSERT IGNORE INTO `role_menus` (`role_id`,`menu_id`)
SELECT 1, `id` FROM `sys_menus` WHERE `path` IN ('/task');

/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
