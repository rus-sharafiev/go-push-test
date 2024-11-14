package push

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id                  *int       `db:"id"`
	AdvertId            *int       `db:"advert_id"`
	TeamId              *int       `db:"team_id"`
	HasPushSubscription *bool      `db:"has_push_subscription"`
	PushToken           *string    `db:"push_token"`
	CreatedAt           *time.Time `db:"created_at"`
	UpdatedAt           *time.Time `db:"updated_at"`
	LastActiveAt        *time.Time `db:"last_active_at"`
	LastPushSentAt      *time.Time `db:"last_push_sent_at"`
}

type Push struct {
	Id       *uuid.UUID         `db:"id"`
	AdvertId *int               `db:"advert_id"`
	Data     *map[string]string `db:"data"`
}

type UserWithToken struct {
	Id        *int    `db:"id"`
	PushToken *string `db:"push_token"`
}

type Message struct {
	Data   map[string]string `json:"data"`
	Tokens []string          `json:"tokens"`
}

type EventPush struct {
	Id        *uuid.UUID         `db:"id"`
	AdvertId  *int               `db:"advert_id"`
	Data      *map[string]string `db:"data"`
	Install   *bool              `db:"install"`
	Reg       *bool              `db:"reg"`
	Dep       *bool              `db:"dep"`
	FourHours *bool              `db:"four_hours"`
	HalfDay   *bool              `db:"half_day"`
	FullDay   *bool              `db:"full_day"`
}
