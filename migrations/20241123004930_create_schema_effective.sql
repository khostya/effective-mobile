-- +goose Up
-- +goose StatementBegin
create schema if not exists effective;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop schema if exists effective;
-- +goose StatementEnd
