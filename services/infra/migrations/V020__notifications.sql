CREATE TABLE IF NOT EXISTS notification_state (
    user_id UUID NOT NULL,
    notif_id VARCHAR(100) NOT NULL,
    state VARCHAR(20) NOT NULL DEFAULT 'Unread',
    snooze_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, notif_id)
);

CREATE TABLE IF NOT EXISTS notification_preferences (
    user_id UUID NOT NULL,
    category VARCHAR(50) NOT NULL,
    in_app BOOLEAN NOT NULL DEFAULT TRUE,
    push BOOLEAN NOT NULL DEFAULT TRUE,
    email BOOLEAN NOT NULL DEFAULT FALSE,
    min_priority VARCHAR(20) NOT NULL DEFAULT 'P4',
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, category)
);
