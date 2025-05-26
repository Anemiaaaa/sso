INSERT INTO apps (id, name, secret)
VALUES ('1', 'test_app', 'test_secret')
ON CONFLICT DO NOTHING;