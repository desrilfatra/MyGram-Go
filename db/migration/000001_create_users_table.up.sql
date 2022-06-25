begin;

create table if not exists users(
Id SERIAL primary key,
Username varchar(100) unique not null,
Email varchar(50) unique not null,
Password varchar(100)not null,
Age int not null,
CreatedAt date not null,
UpdatedAt date not null
);


commit;