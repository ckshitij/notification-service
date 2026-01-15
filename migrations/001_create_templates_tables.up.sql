CREATE TABLE IF NOT EXISTS templates (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,

  name VARCHAR(100) NOT NULL,
  description TEXT,

  channel ENUM('email', 'slack', 'in-app') NOT NULL,
  type ENUM('system', 'user') NOT NULL,

  subject VARCHAR(255),
  body TEXT NOT NULL,

  is_active BOOLEAN NOT NULL DEFAULT TRUE,

  created_by BIGINT NOT NULL,
  updated_by BIGINT NOT NULL,

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  UNIQUE KEY uniq_template (name, type, channel),

  KEY idx_channel_type (channel, type),
  KEY idx_created_by (created_by),
  KEY idx_updated_by (updated_by),
  KEY idx_updated_at (updated_at),
  KEY idx_active (is_active)
);
