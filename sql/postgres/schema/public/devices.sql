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

