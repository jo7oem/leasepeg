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

