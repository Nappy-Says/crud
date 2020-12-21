create table managers (
    id bigserial primary key,
    name text not null,
    login text not null unique,
    password text not null,
    salary integer not null check(salary >0 ),
    plan integer not null check(plan >0 ),
    boss_id bigint not null,
    department text,
    active boolean not null default true,
    created timestamp not null default current_timestamp 

);

create table if not exists customers (
    id bigserial primary key,
    name text not null,
    phone text not null unique,
    password text not null,
    active boolean not null default true,
    created timestamp default current_timestamp
);