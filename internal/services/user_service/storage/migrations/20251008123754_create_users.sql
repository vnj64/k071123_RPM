-- +goose Up
-- +goose StatementBegin
create table users (
    uuid UUID primary key,
    first_name varchar(255),
    second_name varchar(255),
    birth_date timestamp,
    status varchar(100) not null,
    phone_number varchar(255) unique,
    email varchar(255) not null unique,
    role varchar(50) not null,
    created_at timestamp not null default NOW(),
    updated_at timestamp,
    deleted_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users cascade;
-- +goose StatementEnd
