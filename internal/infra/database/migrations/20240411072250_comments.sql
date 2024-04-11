-- +goose Up
-- +goose StatementBegin
create table comment
(
    id           uuid      default gen_random_uuid()            not null
        constraint comment_pk
            primary key,
    parent_id    uuid
        constraint comment_comment_id_fk
            references comment
            on delete cascade,
    entity_id    uuid                                           not null
        constraint comment_article_id_fk
            references article
            on delete cascade
        constraint comment_task_id_fk
            references task
            on delete cascade,
    author_id    uuid                                           not null
        constraint comment_user_id_fk
            references "user"
            on delete cascade,
    comment_text text      default ''                           not null,
    created_at   timestamp default timezone('utc'::text, now()) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table comment;
-- +goose StatementEnd
