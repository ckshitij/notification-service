INSERT IGNORE INTO templates (
  name,
  description,
  channel,
  type,
  active_version,
  created_by,
  updated_by
)
VALUES (
  'onboard_user',
  'Welcome email for new users',
  'email',
  'system',
  1,
  0,
  0
);

INSERT ignore INTO template_versions (
  template_id,
  version,
  subject,
  body,
  is_active
)
VALUES (
  LAST_INSERT_ID(),
  1,
  'Welcome to {{.AppName}}',
  'Hi {{.UserName}},\n\nWelcome to {{.AppName}}!',
  TRUE
);


INSERT ignore INTO templates (
  name,
  description,
  channel,
  type,
  active_version,
  created_by,
  updated_by
)
VALUES (
  'onboard_user',
  'Welcome slack message for new users',
  'slack',
  'system',
  1,
  0,
  0
);

INSERT ignore INTO template_versions (
  template_id,
  version,
  body,
  is_active
)
VALUES (
  LAST_INSERT_ID(),
  1,
  'Welcome *{{.UserName}}* to *{{.AppName}}* 🎉',
  TRUE
);
