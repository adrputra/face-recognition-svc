-- Migration: Create role table
-- Description: Table for storing role information

CREATE TABLE IF NOT EXISTS `role` (
    `id` VARCHAR(255) NOT NULL PRIMARY KEY,
    `role_name` VARCHAR(255) NOT NULL,
    `role_desc` VARCHAR(500) DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_by` VARCHAR(255) DEFAULT NULL,
    `updated_by` VARCHAR(255) DEFAULT NULL,
    `is_active` BOOLEAN DEFAULT TRUE,
    INDEX `idx_role_name` (`role_name`),
    INDEX `idx_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

