-- Newe 数据库初始化脚本
-- 创建数据库
CREATE DATABASE IF NOT EXISTS newe CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE newe;

-- 创建管理员用户
CREATE USER IF NOT EXISTS 'neweuser'@'%' IDENTIFIED BY 'newe123';
GRANT ALL PRIVILEGES ON newe.* TO 'neweuser'@'%';
FLUSH PRIVILEGES;

-- 创建基础表（如果需要）
-- 这里可以添加项目所需的基础表结构

-- 示例：创建系统配置表
CREATE TABLE IF NOT EXISTS sys_config (
    id VARCHAR(36) PRIMARY KEY,
    config_key VARCHAR(100) NOT NULL UNIQUE,
    config_value TEXT,
    config_desc VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入默认配置
INSERT IGNORE INTO sys_config (id, config_key, config_value, config_desc) VALUES
(UUID(), 'system_name', 'Newe Management System', '系统名称'),
(UUID(), 'system_version', '1.0.0', '系统版本'),
(UUID(), 'admin_email', 'admin@newe.com', '管理员邮箱');

-- 创建索引
CREATE INDEX idx_config_key ON sys_config(config_key);
CREATE INDEX idx_created_at ON sys_config(created_at);