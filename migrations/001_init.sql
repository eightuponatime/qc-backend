create extension if not exists "uuid-ossp";

create table votes (
    id uuid primary key default gen_random_uuid(),
    device_id text not null,
    phone_model text not null,
    browser text not null,
    breakfast smallint check (breakfast between 1 and 5),
    lunch smallint check (lunch between 1 and 5),
    dinner smallint check (dinner between 1 and 5),
    external_ip text not null,
    business_date date not null,
    breakfast_at timestamptz,
    lunch_at timestamptz,
    dinner_at timestamptz,
    created_at timestamptz not null default now()
);

create unique index idx_votes_device_business_date
on votes (device_id, business_date);