package push

import (
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// - EVENTS -----------------------------------------------------------------------

func (s service) getEventPushes() ([]EventPush, error) {
	query := `
		SELECT id, advert_id, data, install, reg, dep, four_hours, half_day, full_day FROM event_pushes
		WHERE install = TRUE OR reg = TRUE OR dep = TRUE OR four_hours = TRUE OR half_day = TRUE OR full_day = TRUE;
	`
	rows, _ := s.db.Query(&query)
	pushes, err := pgx.CollectRows(rows, pgx.RowToStructByName[EventPush])
	if err != nil {
		var empty []EventPush
		return empty, err
	}

	return pushes, nil
}

func (s service) getUserTokensAndSendByEvent(push *EventPush, event string, e chan error) {
	query := `
		SELECT count(*) FROM users 
		WHERE advert_id = $1 AND has_push_subscription = TRUE
		AND ` + event + ` <= current_timestamp 
		AND ` + event + ` > (current_timestamp - '1 minute'::interval);
	`
	var qty int
	err := s.db.QueryRow(&query, push.AdvertId).Scan(&qty)
	if err != nil {
		e <- err
		return
	}

	if qty == 0 {
		return
	}

	take := 5
	if Config.BatchSize != nil {
		take = *Config.BatchSize
	}
	batches := qty/take + 1

	for i := 0; i < batches; i++ {
		go func() {
			query := `
				SELECT push_token FROM users 
				WHERE advert_id = @id AND has_push_subscription = TRUE
				AND ` + event + ` <= current_timestamp 
				AND ` + event + ` > (current_timestamp - '1 minute'::interval)
				LIMIT @take OFFSET @skip;
			`
			args := pgx.NamedArgs{
				"id":   push.AdvertId,
				"take": take,
				"skip": take * i,
			}
			rows, _ := s.db.Query(&query, args)
			tokens, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
				var n string
				err := row.Scan(&n)
				return n, err
			})
			if err != nil {
				e <- err
				return
			}

			if err := s.SendBatch(*push.Data, tokens); err != nil {
				e <- err
				return
			}

			e <- nil
		}()
	}
}

// - SCHEDULE ---------------------------------------------------------------------

func (s service) getOneTimePushes(c chan []Push, e chan error) {
	query := `
		SELECT id, advert_id, data FROM one_time_pushes 
		WHERE scheduled_time <= current_timestamp 
		AND scheduled_time > (current_timestamp - @interval::interval);
	`
	interval := 60
	if Config.Interval != nil {
		interval = *Config.Interval
	}
	args := pgx.NamedArgs{
		"interval": strconv.Itoa(interval) + " minute",
	}
	rows, _ := s.db.Query(&query, args)
	pushes, err := pgx.CollectRows(rows, pgx.RowToStructByName[Push])
	if err != nil {
		e <- err
		return
	}
	c <- pushes
}

func (s service) getSchedulePushes(c chan []Push, e chan error) {
	weekDay := strings.ToLower(time.Now().Weekday().String())
	query := `
		SELECT id, advert_id, data FROM schedule_pushes 
		WHERE ` + weekDay + ` IS NOT NULL
		AND ` + weekDay + ` <= current_time 
		AND ` + weekDay + ` > (current_time - @interval::interval);
	`
	interval := 60
	if Config.Interval != nil {
		interval = *Config.Interval
	}
	args := pgx.NamedArgs{
		"interval": strconv.Itoa(interval) + " minute",
	}
	rows, _ := s.db.Query(&query, args)
	pushes, err := pgx.CollectRows(rows, pgx.RowToStructByName[Push])
	if err != nil {
		e <- err
		return
	}
	c <- pushes
}

// --------------------------------------------------------------------------------

func (s service) getUsersQty(advertId *int) (int, error) {
	query := `
		SELECT count(*) FROM users 
		WHERE advert_id = $1 AND has_push_subscription = TRUE;
	`
	var qty int
	err := s.db.QueryRow(&query, advertId).Scan(&qty)
	if err != nil {
		return 0, err
	}
	return qty, nil
}

func (s service) getUserTokensAndSend(push Push, take int, skip int, e chan error) {
	query := `
		SELECT push_token FROM users 
		WHERE advert_id = 1 AND has_push_subscription = TRUE
		LIMIT @take OFFSET @skip;
	`
	args := pgx.NamedArgs{
		"take": take,
		"skip": skip,
	}
	rows, _ := s.db.Query(&query, args)
	tokens, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
		var n string
		err := row.Scan(&n)
		return n, err
	})
	if err != nil {
		e <- err
		return
	}

	if err := s.SendBatch(*push.Data, tokens); err != nil {
		e <- err
		return
	}

	e <- nil
}
