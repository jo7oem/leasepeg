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

create table devices
(
    id         uuid                     default uuidv7() not null
        constraint devices_pk
            primary key,
    user_id    uuid                                      not null
        constraint devices_users_id_fk
            references users,
    name       text                                      not null,
    memo       text                                      not null,
    created_at timestamp with time zone default now()    not null
);

comment on table devices is '管理対象';

alter table devices
    owner to postgres;

create table records
(
    id          uuid                     default uuidv7() not null
        constraint records_pk
            primary key,
    device_id   uuid                                      not null
        constraint records_devices_id_fk
            references devices,
    type        text                                      not null,
    ip_address  inet,
    reported_at timestamp with time zone default now()    not null,
    constraint records_uk_device_id_type
        unique (device_id, type)
);

alter table records
    owner to postgres;

create table credentials
(
    user_id     uuid                                   not null
        constraint credentials_users_id_fk
            references users,
    provider    text                                   not null,
    secret_hash text,
    subject     text                                   not null,
    created_at  timestamp with time zone default now() not null,
    updated_at  timestamp with time zone default now() not null,
    constraint credentials_pk
        primary key (user_id, provider),
    constraint credentials_uq_provider_subject
        unique (provider, subject),
    constraint credentials_local_secret_chk
        check (provider <> 'local' or secret_hash is not null)
);

comment on table credentials is 'ユーザ認証情報';

alter table credentials
    owner to postgres;