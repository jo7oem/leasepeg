create table credentials
(
    user_id     uuid                                   not null,
    provider    text                                   not null,
    secret_hash text,
    subject     text                                   not null,
    created_at  timestamp with time zone default now() not null,
    updated_at  timestamp with time zone default now() not null,
    constraint credentials_pk
        primary key (user_id, provider)
);

comment on table credentials is 'ユーザ認証情報';

alter table credentials
    owner to postgres;

