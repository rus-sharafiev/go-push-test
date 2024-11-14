CREATE TABLE adverts (
    id serial PRIMARY KEY,
    email text UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE teams (
    id serial PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE adverts_in_teams (
    team_id integer REFERENCES teams (id) ON DELETE CASCADE,
    advert_id integer REFERENCES adverts (id) ON DELETE CASCADE,
    parent_advert_id integer REFERENCES adverts (id) ON DELETE CASCADE,
    CONSTRAINT advert_team PRIMARY KEY (advert_id, team_id)
);

CREATE TABLE users (
    id serial PRIMARY KEY,
    advert_id integer REFERENCES adverts (id) ON DELETE CASCADE,
    has_push_subscription boolean DEFAULT TRUE,
    push_token text UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    installed_at TIMESTAMP WITH TIME ZONE,
    registered_at TIMESTAMP WITH TIME ZONE,
    deposit_made_at TIMESTAMP WITH TIME ZONE,
    last_active_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX installed_at_idx ON users (installed_at) WHERE has_push_subscription = TRUE;
CREATE INDEX registered_at_idx ON users (registered_at) WHERE has_push_subscription = TRUE;
CREATE INDEX deposit_made_at_idx ON users (deposit_made_at) WHERE has_push_subscription = TRUE;
CREATE INDEX last_active_at_idx ON users (last_active_at) WHERE has_push_subscription = TRUE;
CREATE INDEX last_push_sent_at_idx ON users (last_push_sent_at) WHERE has_push_subscription = TRUE;

CREATE TABLE one_time_pushes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    advert_id integer REFERENCES adverts (id) ON DELETE CASCADE,
    team_id integer REFERENCES teams (id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    scheduled_time TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE schedule_pushes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    advert_id integer REFERENCES adverts (id) ON DELETE CASCADE,
    team_id integer REFERENCES teams (id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    monday TIME WITH TIME ZONE,
    tuesday TIME WITH TIME ZONE,
    wednesday TIME WITH TIME ZONE,
    thursday TIME WITH TIME ZONE,
    friday TIME WITH TIME ZONE,
    saturday TIME WITH TIME ZONE,
    sunday TIME WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE event_pushes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    advert_id integer REFERENCES adverts (id) ON DELETE CASCADE,
    team_id integer REFERENCES teams (id) ON DELETE CASCADE,
    data JSONB NOT NULL,    
    install boolean DEFAULT false,
    reg boolean DEFAULT false,
    dep boolean DEFAULT false,
    four_hours boolean DEFAULT false,
    half_day boolean DEFAULT false,
    full_day boolean DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TYPE push_status AS ENUM ('SENT', 'REJECTED');
CREATE TABLE sent_pushes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    push_status push_status,
    user_id integer REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);