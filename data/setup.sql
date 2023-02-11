drop table if exists participant_expense;
drop table if exists expense;
drop table if exists participant;
drop table if exists billsplit;
drop table if exists user;



create table participant (
                             id         serial primary key,
                             email_name      varchar(255) not null,
                             billSplit_id    integer references billSplit(id),
                             user_id		integer references user(id),
                             created_at timestamp not null,
                             CONSTRAINT U_Participant UNIQUE (name, billSplit_id)
);

create table billsplit (
                       id         serial primary key,
                       name       varchar(255),
                       created_at timestamp not null
);

create table user (
                             id         serial primary key,
                             name      varchar(64),
                             email      varchar(255) not null,
                             created_at timestamp not null,
);

create table participant_expense (
                        id         serial primary key,
                        participant_id    integer references participant(id),
                        expense_id    integer references expense(id)
);

create table expense (
                        id         serial primary key,
                        name       varchar(255) not null,
                        amount     float8 not null,
                        billSplit_id    integer references billsplit(id),
                        participant_id    integer references participant(id),
                        created_at timestamp not null
);