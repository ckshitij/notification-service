INSERT IGNORE INTO templates (
  name,
  description,
  channel,
  type,
  subject,
  body,
  created_by,
  updated_by
)
VALUES (
  'onboard_user',
  'Welcome email for new users',
  'email',
  'system',
  'Welcome to {{.AppName}}',
  'Hi {{.UserName}},\n\nWelcome to {{.AppName}}!',
  0,
  0
);

INSERT IGNORE INTO templates (
  name,
  description,
  channel,
  type,
  body,
  created_by,
  updated_by
)
VALUES (
  'onboard_user',
  'Welcome slack message for new users',
  'slack',
  'system',
  'Welcome *{{.UserName}}* to *{{.AppName}}* 🎉',
  0,
  0
);