INSERT INTO users (name, email, password, salt, is_verified, created, last_login)
SELECT 'AgentTester', 'agent_tester@empoweredpixels.io', 'LXGLsrFubjvSnMyIJozR+X9Ta9vNMeqUT7DFiai2CT0=', 'AGENT_TESTER_SALT_BASE64_==', true, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'agent_tester@empoweredpixels.io'
);
