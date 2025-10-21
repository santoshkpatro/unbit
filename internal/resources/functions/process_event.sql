REATE OR REPLACE FUNCTION increment_counter()
RETURNS void AS $$
BEGIN
    UPDATE counters SET value = value + 1 WHERE name = 'visits';
END;
$$ LANGUAGE plpgsql;