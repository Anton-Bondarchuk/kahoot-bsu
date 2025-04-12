-- Drop indexes
DROP INDEX IF EXISTS idx_verification_codes_user_id;
DROP INDEX IF EXISTS idx_users_login;
DROP INDEX IF EXISTS idx_answers_question;
DROP INDEX IF EXISTS idx_answers_participant;
DROP INDEX IF EXISTS idx_participants_session;
DROP INDEX IF EXISTS idx_game_sessions_host;
DROP INDEX IF EXISTS idx_options_question_uuid;
DROP INDEX IF EXISTS idx_questions_quiz_uuid;

-- Drop tables in reverse order to handle dependencies
DROP TABLE IF EXISTS answers;
DROP TABLE IF EXISTS participants;
DROP TABLE IF EXISTS game_sessions;
DROP TABLE IF EXISTS options;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS quizzes;
DROP TABLE IF EXISTS verification_codes;
DROP TABLE IF EXISTS users;