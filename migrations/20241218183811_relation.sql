-- +goose Up
-- +goose StatementBegin
create table if not exists users (
    id uuid primary key,
    email text not null
);

insert into users (id, email) values ('e05fa11d-eec3-4fba-b223-d6516800a047', 'test1@exampl.com');
insert into users (id, email) values ('3966749e-45d4-460d-8e59-34235672f03b', 'test2@exampl.com');
insert into users (id, email) values ('b3f7c269-1e35-4139-b882-2ec0b6629f7e', 'test3@exampl.com');
insert into users (id, email) values ('70d77738-2f5f-447e-8fa3-c36b238d9301', 'test4@exampl.com');

alter table refresh_tokens
    add constraint fk_user_id foreign key (user_id) references users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
