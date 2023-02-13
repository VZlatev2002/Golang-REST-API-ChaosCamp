drop table if exists participant_expense;
drop table if exists expense;
drop table if exists participant;
drop table if exists billsplit;
drop table if exists user;


create table billsplit (
                       id         INT UNSIGNED AUTO_INCREMENT,
                       primary key(id),
                       name       varchar(255),
                       created_at timestamp not null
);

create table user (
							 id         INT UNSIGNED AUTO_INCREMENT,
							 key(id),
                             name      varchar(64),
                             email      varchar(255) not null,
                             created_at timestamp not null
);

create table participant (
							 id         INT UNSIGNED AUTO_INCREMENT,
							 primary key(id),
                             email_name      varchar(255) not null,
                             billSplit_id INT UNSIGNED,
                             user_id INT UNSIGNED,
                             foreign key (billSplit_id) references billsplit(id),
                             foreign key (user_id) references user(id),
                             created_at timestamp not null
);

create table expense (
						id         INT UNSIGNED AUTO_INCREMENT,
                        primary key(id),
                        name       varchar(255) not null,
                        amount     float8 not null,
                        billSplit_id INT UNSIGNED,
                        participant_id INT UNSIGNED,
                        foreign key (billSplit_id) references billsplit(id),
                        foreign key (participant_id) references participant(id),
                        created_at timestamp not null
);

create table participant_expense (
						id         INT UNSIGNED AUTO_INCREMENT,
                        primary key(id),
                        participant_id INT UNSIGNED ,
                        expense_id INT UNSIGNED ,
                        foreign key (participant_id) references participant(id),
                        foreign key (expense_id) references expense(id)
);

