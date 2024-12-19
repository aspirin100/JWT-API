-- +goose Up
-- +goose StatementBegin
create table if not exists refresh_tokens
(
    pair_id       uuid primary key,
    user_id       uuid        not null,
    refresh_token bytea       not null,
    created_at    timestamptz not null default now(),
    used          bool        not null default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists refresh_tokens;
-- +goose StatementEnd
