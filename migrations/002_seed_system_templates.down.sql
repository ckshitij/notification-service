use notification;

DELETE tv FROM template_versions tv
JOIN templates t ON tv.template_id = t.id
WHERE t.type = 'system'
  AND t.name IN ('onboard_user');

DELETE FROM templates
WHERE type = 'system'
  AND name IN ('onboard_user');
