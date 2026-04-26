-- Demo seed for analytics/report preview
-- Period: 2026-04-01 .. 2026-04-15
-- Creates 10 demo users per day with breakfast/lunch and partial dinner votes.

begin;

-- Clean only previously generated demo data.
delete from vote_items
where vote_id in (
    select id
    from votes
    where device_id like 'demo-shift-user-%'
);

delete from votes
where device_id like 'demo-shift-user-%';

with period_dates as (
    select generate_series(date '2026-04-01', date '2026-04-15', interval '1 day')::date as business_date
),
devices as (
    select *
    from (
        values
            (1, 'demo-shift-user-01', 'Win32', 'Chrome'),
            (2, 'demo-shift-user-02', 'iPhone', 'Safari'),
            (3, 'demo-shift-user-03', 'Android', 'Chrome'),
            (4, 'demo-shift-user-04', 'MacIntel', 'Safari'),
            (5, 'demo-shift-user-05', 'Linux x86_64', 'Firefox'),
            (6, 'demo-shift-user-06', 'Win32', 'Edge'),
            (7, 'demo-shift-user-07', 'Android', 'Chrome'),
            (8, 'demo-shift-user-08', 'iPhone', 'Safari'),
            (9, 'demo-shift-user-09', 'MacIntel', 'Chrome'),
            (10, 'demo-shift-user-10', 'Win32', 'Chrome')
    ) as t(idx, device_id, phone_model, browser)
)
insert into votes (
    id,
    device_id,
    shift_type,
    phone_model,
    browser,
    external_ip,
    business_date,
    created_at
)
select
    uuid_generate_v5(uuid_ns_url(), format('vote-%s-%s', p.business_date, d.device_id)),
    d.device_id,
    case when d.idx <= 5 then 'day' else 'night' end,
    d.phone_model,
    d.browser,
    '::1',
    p.business_date,
    (p.business_date::text || ' 08:00:00+05')::timestamptz + make_interval(mins => d.idx * 3)
from period_dates p
cross join devices d;

with period_dates as (
    select generate_series(date '2026-04-01', date '2026-04-15', interval '1 day')::date as business_date
),
devices as (
    select *
    from (
        values
            (1, 'demo-shift-user-01'),
            (2, 'demo-shift-user-02'),
            (3, 'demo-shift-user-03'),
            (4, 'demo-shift-user-04'),
            (5, 'demo-shift-user-05'),
            (6, 'demo-shift-user-06'),
            (7, 'demo-shift-user-07'),
            (8, 'demo-shift-user-08'),
            (9, 'demo-shift-user-09'),
            (10, 'demo-shift-user-10')
    ) as t(idx, device_id)
)
insert into vote_items (
    id,
    vote_id,
    meal_type,
    rating,
    review,
    created_at
)
select
    uuid_generate_v5(uuid_ns_url(), format('vote-item-%s-%s-breakfast', p.business_date, d.device_id)),
    uuid_generate_v5(uuid_ns_url(), format('vote-%s-%s', p.business_date, d.device_id)),
    'breakfast',
    (((extract(day from p.business_date)::int + d.idx) % 3) + 3)::smallint,
    case
        when (((extract(day from p.business_date)::int + d.idx) % 3) + 3) = 5 and d.idx % 2 = 0 then 'Очень вкусный завтрак'
        when (((extract(day from p.business_date)::int + d.idx) % 3) + 3) = 5 then 'Все понравилось'
        when (((extract(day from p.business_date)::int + d.idx) % 3) + 3) = 4 and d.idx % 3 = 0 then 'Нормально, без замечаний'
        when (((extract(day from p.business_date)::int + d.idx) % 3) + 3) = 3 and d.idx % 2 = 1 then 'Можно лучше'
        else null
    end,
    (p.business_date::text || ' 08:05:00+05')::timestamptz + make_interval(mins => d.idx * 2)
from period_dates p
cross join devices d;

with period_dates as (
    select generate_series(date '2026-04-01', date '2026-04-15', interval '1 day')::date as business_date
),
devices as (
    select *
    from (
        values
            (1, 'demo-shift-user-01'),
            (2, 'demo-shift-user-02'),
            (3, 'demo-shift-user-03'),
            (4, 'demo-shift-user-04'),
            (5, 'demo-shift-user-05'),
            (6, 'demo-shift-user-06'),
            (7, 'demo-shift-user-07'),
            (8, 'demo-shift-user-08'),
            (9, 'demo-shift-user-09'),
            (10, 'demo-shift-user-10')
    ) as t(idx, device_id)
)
insert into vote_items (
    id,
    vote_id,
    meal_type,
    rating,
    review,
    created_at
)
select
    uuid_generate_v5(uuid_ns_url(), format('vote-item-%s-%s-lunch', p.business_date, d.device_id)),
    uuid_generate_v5(uuid_ns_url(), format('vote-%s-%s', p.business_date, d.device_id)),
    'lunch',
    (((extract(day from p.business_date)::int + d.idx * 2) % 5) + 1)::smallint,
    case
        when (((extract(day from p.business_date)::int + d.idx * 2) % 5) + 1) = 5 then 'Отличный обед'
        when (((extract(day from p.business_date)::int + d.idx * 2) % 5) + 1) = 4 and d.idx % 2 = 0 then 'Хорошо'
        when (((extract(day from p.business_date)::int + d.idx * 2) % 5) + 1) = 3 and d.idx % 3 = 0 then 'Средне'
        when (((extract(day from p.business_date)::int + d.idx * 2) % 5) + 1) = 2 then 'Еда была холодной'
        when (((extract(day from p.business_date)::int + d.idx * 2) % 5) + 1) = 1 then 'Совсем не понравилось'
        else null
    end,
    (p.business_date::text || ' 12:20:00+05')::timestamptz + make_interval(mins => d.idx * 3)
from period_dates p
cross join devices d;

with period_dates as (
    select generate_series(date '2026-04-01', date '2026-04-15', interval '1 day')::date as business_date
),
devices as (
    select *
    from (
        values
            (1, 'demo-shift-user-01'),
            (2, 'demo-shift-user-02'),
            (3, 'demo-shift-user-03'),
            (4, 'demo-shift-user-04'),
            (5, 'demo-shift-user-05'),
            (6, 'demo-shift-user-06'),
            (7, 'demo-shift-user-07'),
            (8, 'demo-shift-user-08'),
            (9, 'demo-shift-user-09'),
            (10, 'demo-shift-user-10')
    ) as t(idx, device_id)
)
insert into vote_items (
    id,
    vote_id,
    meal_type,
    rating,
    review,
    created_at
)
select
    uuid_generate_v5(uuid_ns_url(), format('vote-item-%s-%s-dinner', p.business_date, d.device_id)),
    uuid_generate_v5(uuid_ns_url(), format('vote-%s-%s', p.business_date, d.device_id)),
    'dinner',
    (((extract(day from p.business_date)::int + d.idx * 3) % 5) + 1)::smallint,
    case
        when (((extract(day from p.business_date)::int + d.idx * 3) % 5) + 1) = 5 and d.idx % 2 = 1 then 'Ужин отличный'
        when (((extract(day from p.business_date)::int + d.idx * 3) % 5) + 1) = 4 and d.idx % 2 = 0 then 'В целом хорошо'
        when (((extract(day from p.business_date)::int + d.idx * 3) % 5) + 1) = 3 then 'Нормально, но без восторга'
        when (((extract(day from p.business_date)::int + d.idx * 3) % 5) + 1) = 2 then 'Нужно улучшить качество ужина'
        when (((extract(day from p.business_date)::int + d.idx * 3) % 5) + 1) = 1 then 'Очень слабый ужин'
        else null
    end,
    (p.business_date::text || ' 18:10:00+05')::timestamptz + make_interval(mins => d.idx * 4)
from period_dates p
cross join devices d
where (extract(day from p.business_date)::int + d.idx) % 4 <> 0;

commit;
