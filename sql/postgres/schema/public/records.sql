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

