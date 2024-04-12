-- +goose Up
-- +goose StatementBegin
create table article_comment
(
    id           uuid      default gen_random_uuid()            not null
        constraint article_comment_pk
            primary key,
    parent_id    uuid
        constraint article_comment_article_comment_id_fk
            references article_comment
            on delete cascade,
    entity_id    uuid                                           not null
        constraint article_comment_article_id_fk
            references article
            on delete cascade,
    author_id    uuid                                           not null
        constraint article_comment_user_id_fk
            references "user"
            on delete cascade,
    comment_text text      default ''                           not null,
    created_at   timestamp default timezone('utc'::text, now()) not null
);

create table task_comment
(
    id           uuid      default gen_random_uuid()            not null
        constraint task_comment_pk
            primary key,
    parent_id    uuid
        constraint task_comment_task_comment_id_fk
            references task_comment
            on delete cascade,
    entity_id    uuid                                           not null
        constraint task_comment_task_id_fk
            references task
            on delete cascade,
    author_id    uuid                                           not null
        constraint task_comment_user_id_fk
            references "user"
            on delete cascade,
    comment_text text      default ''                           not null,
    created_at   timestamp default timezone('utc'::text, now()) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table article_comment;

drop table task_comment;
-- +goose StatementEnd
