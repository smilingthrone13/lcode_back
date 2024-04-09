-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION pg_trgm;

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

INSERT INTO "user" ("first_name", "last_name", "username", "password_hash", "is_admin")
VALUES ('Admin', 'Admin', 'admin', '$2a$14$NRd0YacLcLfK6.yOmUUpXeGzzgGsWWYOaXkXZg3DK.9GqF0GEZ/Rq', true);

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
            primary key                   default gen_random_uuid(),
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
    memory           bigint           not null,
    stdout           text,
    stderr           text
);

create unique index solution_result_solution_id_test_case_id_uindex
    on solution_result (solution_id, test_case_id);

CREATE FUNCTION update_solution_metrics()
    RETURNS TRIGGER AS
$$
DECLARE
    trigger_row     RECORD;
    max_runtime_row RECORD;
BEGIN
    IF TG_OP = 'INSERT' THEN
        trigger_row = NEW;
    ELSIF TG_OP = 'DELETE' THEN
        trigger_row = OLD;
    end if;

    SELECT memory, runtime
    INTO max_runtime_row
    FROM solution_result
    WHERE solution_id = trigger_row.solution_id
    ORDER BY runtime DESC
    LIMIT 1;

    UPDATE solution
    SET runtime = COALESCE(max_runtime_row.runtime, 0),
        memory  = COALESCE(max_runtime_row.memory, 0)
    WHERE id = trigger_row.solution_id;

    RETURN trigger_row;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_solution_metrics_trigger
    AFTER INSERT OR DELETE
    ON solution_result
    FOR EACH ROW
EXECUTE FUNCTION update_solution_metrics();

create table article
(
    id         uuid       default gen_random_uuid()            not null
        constraint articles_pk
            primary key,
    author_id  uuid                                            not null
        constraint articles_user_id_fk
            references "user",
    created_at timestamp  default timezone('utc'::text, now()) not null,
    title      varchar(150)                                    not null
        constraint articles_pk_2
            unique,
    content    text       default 'No text yet...'             not null,
    categories text array default '{}'                         not null
);

INSERT INTO article (author_id, title, content, categories)
VALUES ((SELECT id FROM "user" WHERE username = 'admin'),
        'Practice Article',
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus.
 Suspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed, dolor.
 Cras elementum ultrices diam. Maecenas ligula massa, varius a, semper sagittis, dapibus gravida, tellus.
 Nulla vitae elit. Nulla facilisi. Ut fringilla. Suspendisse eu ligula. Etiam porta sem.',
        '{"Practice"}');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table article;

drop table task_template;

drop table solution_result;

drop table solution;

drop table "user";

drop table test_case;

drop table task;

drop function update_task_number() cascade;

drop function update_test_case_number() cascade;

drop function update_solution_metrics() cascade;

drop extension pg_trgm cascade;
-- +goose StatementEnd
