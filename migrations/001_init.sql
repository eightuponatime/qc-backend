create extension if not exists "uuid-ossp";

create table votes (
    id uuid primary key default gen_random_uuid(),
    device_id text not null,
    phone_model text not null,
    browser text not null,
    external_ip text not null,
    business_date date not null,
    created_at timestamptz not null default now()
);

create unique index idx_votes_device_business_date
on votes (device_id, business_date);

create table vote_items (
    id uuid primary key default gen_random_uuid(),
    vote_id uuid not null references votes(id) on delete cascade,
    meal_type text not null check (meal_type in ('breakfast','lunch','dinner')),
    rating smallint check (rating between 1 and 5),
    review text,
    created_at timestamptz not null default now()
);

create unique index idx_vote_items_unique
on vote_items (vote_id, meal_type);