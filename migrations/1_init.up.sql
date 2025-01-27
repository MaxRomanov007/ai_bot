CREATE TABLE IF NOT EXISTS user_roles(
    user_role_id BIGSERIAL PRIMARY KEY,
    user_role_name VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_user_role_name ON user_roles(user_role_name);

INSERT INTO user_roles (user_role_name) VALUES
    ('admin'),
    ('user')
ON CONFLICT (user_role_name) DO NOTHING;

CREATE TABLE IF NOT EXISTS user_statuses(
    user_status_id BIGSERIAL PRIMARY KEY,
    user_status_name VARCHAR(255)
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_user_status_name ON user_statuses(user_status_name);

INSERT INTO user_statuses (user_status_name) VALUES
    ('authorized'),
    ('unauthorized'),
    ('blocked')
ON CONFLICT (user_status_name) DO NOTHING;

CREATE TABLE IF NOT EXISTS users(
    user_id BIGINT PRIMARY KEY,
    username VARCHAR(255),
    chat_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    status_id BIGINT NOT NULL,

    FOREIGN KEY (role_id) REFERENCES user_roles(user_role_id) ON DELETE CASCADE,
    FOREIGN KEY (status_id) REFERENCES user_statuses(user_status_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);
CREATE INDEX IF NOT EXISTS idx_users_status_id ON users(status_id);

CREATE TABLE IF NOT EXISTS message_roles(
    message_role_id BIGSERIAL PRIMARY KEY,
    message_role_name VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_message_role_name ON message_roles(message_role_name);

INSERT INTO message_roles (message_role_name) VALUES
    ('user'),
    ('assistant')
ON CONFLICT (message_role_name) DO NOTHING;

CREATE TABLE IF NOT EXISTS messages (
    message_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    content TEXT NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES message_roles(message_role_id) ON DELETE CASCADE
);