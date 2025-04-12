-- Current Date and Time (UTC): 2025-04-11 08:00:04
-- Current User's Login: Anton-Bondarchuk
-- Description:
-- Admin panel and auth functionality migration

-- Create users table
CREATE TABLE users (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    login VARCHAR(255) NOT NULL UNIQUE,
    role_flags INTEGER NOT NULL DEFAULT 1, -- Bit mask for roles:
                                          -- 1 = Regular User (0001)
                                          -- 2 = Admin       (0010)
                                          -- 4 = Teacher   (0100)
                                          -- 8 = Blocked     (1000)
);

-- Create verification_codes table
CREATE TABLE verification_codes (
    id UUID PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code VARCHAR(10) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create quizzes table
CREATE TABLE quizzes (
    id UUID PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    title VARCHAR(255) NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create questions table
CREATE TABLE questions (
    id UUID PRIMARY KEY,
    quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    time_limit INTEGER NOT NULL DEFAULT 30, -- in seconds
    points INTEGER NOT NULL DEFAULT 100,
    position INTEGER NOT NULL DEFAULT 0, -- for ordering questions
);

-- Create options table for question answers
CREATE TABLE options (
    id UUID PRIMARY KEY,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    
    text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    position INTEGER NOT NULL DEFAULT 0 -- for ordering options
);

-- Create game_sessions table for active games
CREATE TABLE game_sessions (
    id UUID PRIMARY KEY,
    quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    host_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    join_code VARCHAR(8) NOT NULL UNIQUE,
    status_flags INTEGER NOT NULL DEFAULT 1, -- Bit mask for game status:
                                           -- 1 = Waiting  (0001)
                                           -- 2 = Active   (0010)
                                           -- 4 = Paused   (0100)
                                           -- 8 = Finished (1000)
    current_question_index INTEGER DEFAULT 0,
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
);

-- Create participants table for users in a game
CREATE TABLE participants (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    login VARCHAR(50) NOT NULL,
    
    score INTEGER NOT NULL DEFAULT 0,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create answers table for participant responses
CREATE TABLE answers (
    id UUID PRIMARY KEY,
    participant_id UUID NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    option_id UUID REFERENCES options(id) ON DELETE SET NULL,
    
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    response_time_ms INTEGER, -- response time in milliseconds
    points_awarded INTEGER NOT NULL DEFAULT 0,
    answered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_questions_quiz_id ON questions(quiz_id);
CREATE INDEX idx_options_question_id ON options(question_id);
CREATE INDEX idx_game_sessions_host ON game_sessions(host_id);
CREATE INDEX idx_participants_session ON participants(session_id);
CREATE INDEX idx_answers_participant ON answers(participant_id);
CREATE INDEX idx_answers_question ON answers(question_id);
CREATE INDEX idx_users_login ON users(login);
CREATE INDEX idx_verification_codes_user_id ON verification_codes(user_id);