-- +goose Up
-- +goose StatementBegin
create table "user"
(
    id            uuid    default gen_random_uuid() not null
        constraint user_pk
            primary key,
    first_name    varchar(50)                       not null,
    last_name     varchar(50)                       not null,
    username      varchar(50)                       not null unique,
    password_hash text                              not null,
    is_admin      boolean default false             not null
);

create table task
(
    id            uuid             default gen_random_uuid() not null
        constraint task_pk
            primary key,
    number        bigint                                     not null,
    name          text                                       not null,
    description   text             default ''::text          not null,
    difficulty    text                                       not null,
    category      text                                       not null,
    runtime_limit double precision default 5.0               not null,
    memory_limit  bigint           default 128000            not null
);

create unique index task_name_uindex
    on task (name);

CREATE FUNCTION update_task_number() RETURNS TRIGGER
    LANGUAGE plpgsql
AS
$$
BEGIN
    IF TG_OP = 'INSERT' THEN
        SELECT (COALESCE(MAX(number), 0) + 1) INTO NEW.number FROM task;
        RETURN NEW;

    ELSIF TG_OP = 'DELETE' THEN
        UPDATE task
        SET number = number - 1
        WHERE number > (SELECT number FROM task WHERE id = OLD.id);
        RETURN OLD;

    END IF;
END;
$$;

create trigger update_task_number_trigger
    before insert or delete
    on task
    for each row
execute procedure update_task_number();

create table task_template
(
    id          uuid default gen_random_uuid() not null
        constraint task_template_pk
            primary key,
    task_id     uuid                           not null
        constraint task_template_task_id_fk
            references task
            on delete cascade,
    language_id integer                        not null,
    template    text                           not null,
    wrapper     text                           not null
);

create unique index task_template__index
    on task_template (task_id, language_id);

create table solution
(
    id          uuid             not null
        constraint solution_pk
            primary key default gen_random_uuid(),
    task_id     uuid             not null
        constraint solution_task_id_fk
            references task
            on delete cascade,
    language_id integer,
    user_id     uuid             not null
        constraint solution_user_id_fk
            references "user"
            on delete cascade,
    code        text             not null,
    status      text             not null,
    runtime     double precision not null default 0,
    memory      bigint           not null default 0
);

create table test_case
(
    id      uuid default gen_random_uuid() not null
        constraint test_case_pk
            primary key,
    task_id uuid                           not null
        constraint test_case_task_id_fk
            references task
            on delete cascade,
    number  bigint                         not null,
    input   text                           not null,
    output  text                           not null
);

CREATE FUNCTION update_test_case_number() RETURNS TRIGGER
    LANGUAGE plpgsql
AS
$$
BEGIN
    IF TG_OP = 'INSERT' THEN
        SELECT (COALESCE(MAX(number), 0) + 1)
        INTO NEW.number
        FROM test_case
        WHERE task_id = NEW.task_id;
        RETURN NEW;

    ELSIF TG_OP = 'DELETE' THEN
        UPDATE test_case
        SET number = number - 1
        WHERE task_id = OLD.task_id
          AND number > (SELECT number FROM test_case WHERE id = OLD.id);
        RETURN OLD;

    END IF;
END;
$$;

create trigger update_test_case_number_trigger
    before insert or delete
    on test_case
    for each row
execute procedure update_test_case_number();

create table solution_result
(
    solution_id      uuid             not null
        constraint solution_result_solution_id_fk
            references solution
            on delete cascade,
    test_case_id     uuid             not null
        constraint solution_result_test_case_id_fk
            references test_case
            on delete cascade,
    submission_token uuid             not null,
    status           integer          not null,
    runtime          double precision not null,
    memory           bigint           not null
);

create unique index solution_result_solution_id_test_case_id_uindex
    on solution_result (solution_id, test_case_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table user_progress;

drop table task_template;

drop table solution_result;

drop table solution;

drop table "user";

drop table test_case;

drop table task;

drop function update_task_number();

drop function update_test_case_number();
-- +goose StatementEnd
