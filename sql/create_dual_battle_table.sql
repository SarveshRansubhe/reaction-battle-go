CREATE TABLE IF NOT EXISTS public.dual_battle (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id),
    red_scores REAL[] NOT NULL,
    blue_scores REAL[] NOT NULL,
    red_penalties REAL[] NOT NULL,
    blue_penalties REAL[] NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    is_completed boolean DEFAULT false
);
