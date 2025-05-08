-- Create notifications table
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('email', 'sms', 'push')),
    sent_at TIMESTAMPTZ DEFAULT now()
);
