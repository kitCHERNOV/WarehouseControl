-- Migration: 00002_create_triggers
-- Description: Create triggers for automatic history logging (ANTI-PATTERN for educational purposes)

-- Enable extension for JSONB operations if not already enabled
CREATE EXTENSION IF NOT EXISTS plpgsql;

-- Function to log item changes
-- This function is called by triggers after INSERT, UPDATE, DELETE operations
CREATE OR REPLACE FUNCTION log_item_change()
RETURNS TRIGGER AS $$
DECLARE
    user_context TEXT;
    user_id TEXT := 'system';
    user_name TEXT := 'system';
BEGIN
    -- Get current user from application context (set via SET LOCAL)
    -- Format: "user_id:username"
    BEGIN
        user_context := current_setting('app.current_user', true);
        IF user_context IS NOT NULL THEN
            -- Parse user_context (format: "user_id:username")
            user_id := split_part(user_context, ':', 1);
            user_name := split_part(user_context, ':', 2);
        END IF;
    EXCEPTION WHEN OTHERS THEN
        -- If setting not found, use system as default
        user_id := 'system';
        user_name := 'system';
    END;

    IF TG_OP = 'INSERT' THEN
        -- Log new item creation
        INSERT INTO item_history (
            item_id,
            action,
            old_name,
            new_name,
            old_quantity,
            new_quantity,
            user_id,
            user_name,
            changed_at
        ) VALUES (
            NEW.id,
            'INSERT',
            NULL,
            NEW.name,
            NULL,
            NEW.quantity,
            user_id,
            user_name,
            CURRENT_TIMESTAMP
        );
        RETURN NEW;

    ELSIF TG_OP = 'UPDATE' THEN
        -- Log item update (only if values actually changed)
        IF (OLD.name IS DISTINCT FROM NEW.name) OR (OLD.quantity IS DISTINCT FROM NEW.quantity) THEN
            INSERT INTO item_history (
                item_id,
                action,
                old_name,
                new_name,
                old_quantity,
                new_quantity,
                user_id,
                user_name,
                changed_at
            ) VALUES (
                NEW.id,
                'UPDATE',
                OLD.name,
                NEW.name,
                OLD.quantity,
                NEW.quantity,
                user_id,
                user_name,
                CURRENT_TIMESTAMP
            );
        END IF;
        RETURN NEW;

    ELSIF TG_OP = 'DELETE' THEN
        -- Log item deletion
        INSERT INTO item_history (
            item_id,
            action,
            old_name,
            new_name,
            old_quantity,
            new_quantity,
            user_id,
            user_name,
            changed_at
        ) VALUES (
            OLD.id,
            'DELETE',
            OLD.name,
            NULL,
            OLD.quantity,
            NULL,
            user_id,
            user_name,
            CURRENT_TIMESTAMP
        );
        RETURN OLD;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for INSERT operations
DROP TRIGGER IF EXISTS trigger_items_insert ON items;
CREATE TRIGGER trigger_items_insert
    AFTER INSERT ON items
    FOR EACH ROW
    EXECUTE FUNCTION log_item_change();

-- Create trigger for UPDATE operations
DROP TRIGGER IF EXISTS trigger_items_update ON items;
CREATE TRIGGER trigger_items_update
    AFTER UPDATE ON items
    FOR EACH ROW
    EXECUTE FUNCTION log_item_change();

-- Create trigger for DELETE operations
DROP TRIGGER IF EXISTS trigger_items_delete ON items;
CREATE TRIGGER trigger_items_delete
    AFTER DELETE ON items
    FOR EACH ROW
    EXECUTE FUNCTION log_item_change();

-- Add comment documenting this as an anti-pattern
COMMENT ON FUNCTION log_item_change() IS '
ANTI-PATTERN: This function implements history logging via database triggers.
Educational purpose only - demonstrates why triggers should be avoided for business logic.

Problems with this approach:
1. Business logic hidden in database
2. Hard to test (requires real DB)
3. Migration difficulties
4. Database vendor lock-in
5. Synchronous execution (blocks operations)
6. Cannot add application-level logic

Better alternatives:
- Application-level logging in service layer
- Event Sourcing pattern
- Audit Log microservice
- Change Data Capture (CDC) tools like Debezium
';
