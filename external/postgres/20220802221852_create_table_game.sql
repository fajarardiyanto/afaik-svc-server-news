-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Game(
    id uuid not null primary key ,
    code varchar(255) not null unique ,
    ground_id uuid not null ,
    event_id uuid null ,
    name varchar(255) not null ,
    start_date timestamp without time zone ,
    end_date timestamp without time zone ,
    metadata json not null default '{}' ,
    modified_by json not null default '{}' ,
    created_at timestamp without time zone not null ,
    updated_at timestamp without time zone not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Game;
-- +goose StatementEnd
