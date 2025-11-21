-- +goose Up
-- +goose StatementBegin
create table units (
    uuid UUID primary key,
    status varchar(50) not null,
    network_status varchar(50) not null,
    direction varchar(3) not null, -- in/out
    code varchar(6),
    qr_link text,
    parking_uuid UUID references parkings(uuid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table units cascade;
-- +goose StatementEnd
