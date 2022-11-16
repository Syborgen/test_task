SELECT * FROM tech_windows WHERE duration <@ '[2022-02-01 00:00:00,2022-03-30 23:59:59]'::tsrange;

SELECT id_object, count(duration) AS windows_count, avg(get_interval(duration)) AS average_duration FROM tech_windows GROUP BY id_object;

SELECT id_object, count(duration) AS windows_count, avg(get_interval(duration)) AS average_duration FROM tech_windows WHERE duration <@ '[2022-02-01 00:00:00,2022-03-30 23:59:59]'::tsrange GROUP BY id_object ORDER BY average_duration DESC;

CREATE INDEX average_window_duration ON tech_windows ((get_interval(duration)) ASC);