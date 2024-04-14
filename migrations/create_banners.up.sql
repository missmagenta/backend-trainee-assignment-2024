SET statement_timeout = 0;

--bun:split

create schema if not exists banners;

create table if not exists banners.banners
(
    id         serial primary key       not null,

    feature_id int                      not null,
    content    text                     not null,
    is_active  bool                     not null,

    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null
);

create table if not exists banners.tags
(
    banner_id  int not null references banners.banners (id) on delete cascade,
    feature_id int not null,
    tag_id     int not null,
    primary key (feature_id, tag_id)
)
