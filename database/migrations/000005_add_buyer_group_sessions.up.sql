-- Migration: Add buyer group session tables
-- Created: 2025-12-28

-- 1. Create buyer_group_sessions table
CREATE TABLE IF NOT EXISTS buyer_group_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_code TEXT NOT NULL UNIQUE,
    organizer_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_variant_id UUID NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    title TEXT,
    current_participants INTEGER DEFAULT 0,
    status TEXT DEFAULT 'open' CHECK (status IN ('open', 'locked', 'completed', 'cancelled', 'expired')),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_buyer_group_sessions_organizer ON buyer_group_sessions(organizer_user_id);
CREATE INDEX idx_buyer_group_sessions_status ON buyer_group_sessions(status);
CREATE INDEX idx_buyer_group_sessions_expires_at ON buyer_group_sessions(expires_at);

-- 2. Create buyer_group_members table  
CREATE TABLE IF NOT EXISTS buyer_group_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES buyer_group_sessions(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    status TEXT DEFAULT 'joined' CHECK (status IN ('joined', 'paid', 'cancelled')),
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- Prevent duplicate members in same session
    UNIQUE(session_id, user_id)
);

CREATE INDEX idx_buyer_group_members_session ON buyer_group_members(session_id);
CREATE INDEX idx_buyer_group_members_user ON buyer_group_members(user_id);
CREATE INDEX idx_buyer_group_members_status ON buyer_group_members(status);
