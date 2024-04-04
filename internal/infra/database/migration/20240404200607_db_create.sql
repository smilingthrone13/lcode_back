-- +goose Up
-- +goose StatementBegin
create table "user"
(
    id            uuid    default gen_random_uuid() not null
        constraint user_pk
            primary key,
    first_name    varchar(50)                       not null,
    last_name     varchar(50)                       not null,
    username      varchar(50)                       not null,
    password_hash text                              not null,
    is_admin      boolean default false             not null
);

alter table "user"
    owner to postgres;

create table user_progress
(
    user_id uuid not null
        constraint user_progress_user_id_fk
            references "user",
    task_id uuid not null,
    status  text not null
);

alter table user_progress
    owner to postgres;

create unique index user_progress_user_id_task_id_uindex
    on user_progress (user_id, task_id);

create table task
(
    id          uuid default gen_random_uuid() not null
        constraint task_pk
            primary key,
    number      bigserial,
    name        text                           not null,
    description text default ''::text          not null,
    difficulty  text                           not null,
    category    text                           not null
);

alter table task
    owner to postgres;

CREATE OR REPLACE FUNCTION update_task_number()
    RETURNS TRIGGER AS $$
DECLARE
max_number BIGINT;
    deleted_row RECORD; -- объявляем переменную для хранения удаленной строки
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Ищем максимальный номер для задачи с таким же task_id
SELECT COALESCE(MAX(number), 0) + 1 INTO max_number
FROM task;

-- Устанавливаем номер
NEW.number := max_number;
    ELSIF TG_OP = 'DELETE' THEN
        -- Получаем номер удаляемой записи
SELECT number INTO deleted_row
FROM task
WHERE id = OLD.id;

-- Уменьшаем номера последующих записей
UPDATE task
SET number = number - 1
WHERE number > deleted_row;
END IF;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

create trigger update_task_number_trigger
    after insert or delete
on task
    for each row
execute procedure update_task_number();

create table task_template
(
    id       uuid default gen_random_uuid() not null
        constraint task_template_pk
            primary key,
    task_id  uuid                           not null
        constraint task_template_task_id_fk
            references task,
    language text                           not null,
    template text                           not null,
    wrapper  text                           not null
);

alter table task_template
    owner to postgres;

create table submission
(
    user_id          uuid             not null
        constraint submission_user_id_fk
            references "user",
    task_template_id uuid             not null
        constraint submission_task_template_id_fk
            references task_template,
    solution         text             not null,
    status           text             not null,
    runtime          bigint           not null,
    memory           double precision not null
);

alter table submission
    owner to postgres;

create table test_case
(
    id      uuid default gen_random_uuid() not null
        constraint test_case_pk
            primary key,
    task_id uuid                           not null
        constraint test_case_task_id_fk
            references task,
    number  bigserial,
    input   json                           not null,
    output  json                           not null
);

alter table test_case
    owner to postgres;

CREATE OR REPLACE FUNCTION update_test_case_number()
RETURNS TRIGGER AS $$
DECLARE
max_number BIGINT;
    deleted_row RECORD; -- объявляем переменную для хранения удаленной строки
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Ищем максимальный номер для задачи с таким же task_id
SELECT COALESCE(MAX(number), 0) + 1 INTO max_number
FROM test_case
WHERE task_id = NEW.task_id;

-- Устанавливаем номер
NEW.number := max_number;
    ELSIF TG_OP = 'DELETE' THEN
        -- Получаем номер удаляемой записи
SELECT number INTO deleted_row
FROM test_case
WHERE id = OLD.id;

-- Уменьшаем номера последующих записей
UPDATE test_case
SET number = number - 1
WHERE task_id = OLD.task_id
  AND number > deleted_row;
END IF;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

create trigger update_test_case_number_trigger
    after insert or delete
on test_case
    for each row
execute procedure update_test_case_number();


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
