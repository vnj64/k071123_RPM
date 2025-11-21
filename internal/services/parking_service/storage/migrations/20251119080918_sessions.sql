-- +goose Up
-- +goose StatementBegin
create table sessions (
    uuid UUID primary key,
    parking_uuid UUID not null references parkings(uuid),
    car_uuid UUID not null references cars(uuid),
    status varchar(50) not null,
    start_at timestamp not null default NOW(),
    finish_at timestamp,
    cost float
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table sessions cascade;
-- +goose StatementEnd
