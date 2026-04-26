create extension if not exists "uuid-ossp";

create table votes (
    id uuid primary key default uuid_generate_v4(),
    device_id text not null,
    shift_type text not null check (shift_type in ('day', 'night')),
    phone_model text not null,
    browser text not null,
    external_ip text not null,
    business_date date not null,
    created_at timestamptz not null default now()
);
create unique index idx_votes_device_business_date
on votes (device_id, business_date);

create table vote_items (
    id uuid primary key default uuid_generate_v4(),
    vote_id uuid not null references votes(id) on delete cascade,
    meal_type text not null check (meal_type in ('breakfast','lunch','dinner')),
    rating smallint check (rating between 1 and 5),
    review text,
    created_at timestamptz not null default now()
);
create unique index idx_vote_items_unique
on vote_items (vote_id, meal_type);

create table sent_reports (
    id uuid primary key default uuid_generate_v4(),
    period_start date not null,
    period_end date not null,
    sent_at timestamptz not null default now()
);
create unique index idx_sent_reports_period_unique
on sent_reports (period_start, period_end);

create table analytics_access_codes (
    id uuid primary key default uuid_generate_v4(),
    code_hash text not null unique,
    valid_from date not null,
    valid_until date not null,
    created_at timestamptz not null default now()
);
