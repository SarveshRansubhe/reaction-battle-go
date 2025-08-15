CREATE TABLE IF NOT EXISTS public.single_battle (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id),
    scores REAL[] NOT NULL,
    penalties REAL[] NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    is_completed boolean DEFAULT false
);