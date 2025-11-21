-- +goose Up
-- +goose StatementBegin
create table car_settings (
                              uuid UUID primary key,
                              vin varchar(255) not null default ''
);

create table cars (
    uuid UUID primary key,
    user_uuid UUID not null,
    gos_number varchar(8) not null unique,
    is_active bool not null default true,
    settings_uuid UUID references car_settings(uuid),
    created_at timestamp not null default NOW(),
    updated_at timestamp,
    deleted_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table cars cascade;
drop table car_settings cascade;
-- +goose StatementEnd
