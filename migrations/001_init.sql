
-- [person] * ----------- * [company]
--
-- [person] 1 ----------- * [company_person] * ----------- 1 [company]

create table person (
    id bigserial primary key,
    first_name text,
    last_name text,
    dob date
);

create table company (
    id bigserial primary key,
    name text
);

create table company_person (
    company_id bigint references company(id),
    person_id bigint references person(id),
    primary key (company_id, person_id)
);

