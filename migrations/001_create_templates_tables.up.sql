CREATE TABLE IF NOT EXISTS templates (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,

  name VARCHAR(100) NOT NULL,
  description TEXT,

  channel ENUM('email', 'slack', 'in_app') NOT NULL,
  type ENUM('system', 'user') NOT NULL,

  active_version INT NOT NULL DEFAULT 1,

  created_by BIGINT NOT NULL,
  updated_by BIGINT NOT NULL,

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  UNIQUE KEY uniq_template (name, type, channel)
);

CREATE TABLE IF NOT EXISTS template_versions (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,

  template_id BIGINT NOT NULL,
  version INT NOT NULL,

  subject VARCHAR(255),
  body TEXT NOT NULL,

  is_active BOOLEAN NOT NULL DEFAULT FALSE,

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  UNIQUE KEY uniq_template_version (template_id, version),
  INDEX idx_template_active (template_id, is_active),

  CONSTRAINT fk_template_version
    FOREIGN KEY (template_id)
    REFERENCES templates(id)
    ON DELETE CASCADE
);
