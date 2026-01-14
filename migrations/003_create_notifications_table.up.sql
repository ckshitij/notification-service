CREATE TABLE IF NOT EXISTS notifications (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,

  channel ENUM('email', 'slack', 'in_app') NOT NULL,

  template_version_id BIGINT NOT NULL,

  recipient JSON NOT NULL,
  payload JSON NOT NULL,

  status VARCHAR(20) NOT NULL,

  scheduled_at DATETIME NULL,
  sent_at DATETIME NULL,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    ON UPDATE CURRENT_TIMESTAMP,

  INDEX idx_status (status),
  INDEX idx_scheduled (scheduled_at),
  INDEX idx_template_version (template_version_id)
);
