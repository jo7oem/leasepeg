create table users
(
    id           uuid                                   not null
        primary key,
    user_name    varchar(32)                            not null
        constraint users_username_key
            unique,
    created_at   timestamp with time zone default now() not null,
    display_name text                                   not null,
    updated_at   timestamp with time zone default now() not null
);

comment on table users is 'ユーザ情報';

alter table users
    owner to postgres;

