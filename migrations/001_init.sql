create extension if not exists "uuid-ossp";

create table votes (
    id uuid primary key default gen_random_uuid(),
    device_id text not null,
    breakfast smallint check (breakfast between 1 and 5),
    lunch smallint check (lunch between 1 and 5),
    dinner smallint check (dinner between 1 and 5),
    external_ip text not null,
    voted_at timestampz not null default now()
);

create unique index idx_votes_device_day on votes (device_id, DATE(voted_at));