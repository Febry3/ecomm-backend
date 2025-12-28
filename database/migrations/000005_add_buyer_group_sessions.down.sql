-- Rollback: Remove buyer group session tables

-- Drop indexes first
DROP INDEX IF EXISTS idx_buyer_group_members_status;
DROP INDEX IF EXISTS idx_buyer_group_members_user;
DROP INDEX IF EXISTS idx_buyer_group_members_session;
DROP INDEX IF EXISTS idx_buyer_group_sessions_expires_at;
DROP INDEX IF EXISTS idx_buyer_group_sessions_status;
DROP INDEX IF EXISTS idx_buyer_group_sessions_organizer;

-- Drop tables (members first due to FK constraint)
DROP TABLE IF EXISTS buyer_group_members;
DROP TABLE IF EXISTS buyer_group_sessions;
