INSERT INTO adverts(email) VALUES ('mail@mail.com');

INSERT INTO teams(id) VALUES (1);

INSERT INTO users(id, advert_id, has_push_subscription, push_token)
SELECT s.id, 1, random() < 0.5, gen_random_uuid()
FROM generate_series(1,20) AS s(id);

UPDATE users SET installed_at = current_timestamp WHERE id = 1;
UPDATE users SET registered_at = (current_timestamp + '1 minute'::interval) WHERE id = 3;
UPDATE users SET deposit_made_at = (current_timestamp + '2 minute'::interval) WHERE id = 8;

INSERT INTO one_time_pushes(advert_id, data, scheduled_time) 
VALUES (1, '{"title": "one_time_pushes", "notification": "text"}', current_timestamp);
UPDATE one_time_pushes SET scheduled_time = current_timestamp + interval '1 minute';

INSERT INTO schedule_pushes(advert_id, data, thursday) 
VALUES (1, '{"title": "schedule_pushes", "notification": "text"}', current_timestamp);
UPDATE schedule_pushes SET thursday = current_time;

INSERT INTO event_pushes(advert_id, data, install, reg, dep) 
VALUES (1, '{"title": "event_pushes", "notification": "install, reg, dep"}', TRUE, TRUE, TRUE);

INSERT INTO event_pushes(advert_id, data, four_hours, half_day, full_day) 
VALUES (1, '{"title": "event_pushes", "notification": "four_hours, half_day, full_day"}', TRUE, TRUE, TRUE);