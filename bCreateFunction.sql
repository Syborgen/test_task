CREATE OR REPLACE FUNCTION generate(objects_count integer, windows_count integer) RETURNS void AS $$
    DECLARE
        period_with_one_tech_window_duration INTERVAL := INTERVAL '1 year'/windows_count;
        start_of_current_period TIMESTAMP;
        end_of_current_period TIMESTAMP;
        duration_as_text TEXT;
    BEGIN
        DELETE FROM objects WHERE true;
        FOR i IN 1..objects_count LOOP 
            INSERT INTO objects(name, clock) VALUES('object ' || i, i%27 - 12);
            
            FOR j IN 0..windows_count-1 LOOP
                start_of_current_period := TIMESTAMP '2022-01-01 00:00:00' + period_with_one_tech_window_duration * j + (random() * (period_with_one_tech_window_duration - INTERVAL '23 hours') + INTERVAL '1 hour');
                end_of_current_period := start_of_current_period + random() * INTERVAL '24 hours';
                duration_as_text := TEXT '[' || to_char(start_of_current_period,'YYYY-MM-DD HH24:MI:SS') || ', ' || to_char(end_of_current_period,'YYYY-MM-DD HH24:MI:SS') || TEXT ')';

                INSERT INTO tech_windows(id_object, duration)
                    VALUES(
                        (SELECT id FROM objects WHERE name='object ' || i), 
                        tsrange (duration_as_text)
                    );
            END LOOP;
        END LOOP;
    END
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_interval(tsrange) RETURNS INTERVAL AS $$
    BEGIN
        RETURN upper($1) - lower($1);
    END
$$ LANGUAGE plpgsql IMMUTABLE;