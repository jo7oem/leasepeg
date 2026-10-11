create table users
(
    id           uuid                                                   not null,
    user_name    varchar(32)                                            not null,
    created_at   timestamp with time zone default now()                 not null,
    display_name varchar(32)              default ''::character varying not null,
    updated_at   timestamp with time zone default now()                 not null,
    primary key (id),
    constraint users_username_key
        unique (user_name)
);

alter table users
    owner to postgres;

