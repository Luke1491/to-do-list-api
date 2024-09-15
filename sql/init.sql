CREATE TABLE IF NOT EXISTS todo_lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS todo_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_id UUID REFERENCES todo_lists(id) ON DELETE CASCADE,
    description VARCHAR(255) NOT NULL,
    is_checked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
