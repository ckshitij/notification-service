CREATE TABLE IF NOT EXISTS notifications (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,

  channel ENUM('email', 'slack', 'in-app') NOT NULL,

  template_id BIGINT NOT NULL,

  recipient JSON NOT NULL,
  template_kv JSON NOT NULL,

  status VARCHAR(20) NOT NULL,

  scheduled_at DATETIME NULL,
  sent_at DATETIME NULL,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    ON UPDATE CURRENT_TIMESTAMP,

  -- Indexes
  INDEX idx_status (status),
  INDEX idx_scheduled (scheduled_at),
  INDEX idx_template_id (template_id),

  -- Foreign key constraint
  CONSTRAINT fk_notifications_template
    FOREIGN KEY (template_id)
    REFERENCES templates(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
) ENGINE=InnoDB;
