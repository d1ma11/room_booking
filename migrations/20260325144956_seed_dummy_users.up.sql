INSERT INTO users (id, email, role, password, created_at)
SELECT '00000000-0000-0000-0000-000000000001',
       'admin@dummy.local',
       'admin',
       NULL,
       NOW() WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE id = '00000000-0000-0000-0000-000000000001'
);

INSERT INTO users (id, email, role, password, created_at)
SELECT '00000000-0000-0000-0000-000000000002',
       'user@dummy.local',
       'user',
       NULL,
       NOW() WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE id = '00000000-0000-0000-0000-000000000002'
);