-- +goose Up
-- +goose StatementBegin
create table parkings (
    uuid UUID primary key,
    tariff_uuid UUID not null references tariffs(uuid),
    name varchar(255) not null,
    address varchar(255) not null,
    latitude varchar(255) not null,
    longitude varchar(255) not null,
    total_places integer not null,
    status varchar(50) not null,
    created_at timestamp not null default NOW(),
    updated_at timestamp,
    deleted_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table parkings cascade;
-- +goose StatementEnd
