
create table if not exists patient (
  id integer primary key,
  name string not null,
  surname string not null,
  birthdate date not null,
  note_id integer,
  deleted boolean default false,
  foreign key(note_id) references note(id) 
);

create table if not exists note (
  id integer primary key,
  deleted boolean default false,
  text string not null
);

create table if not exists visit (
  id integer primary key,
  vdatetime datetime not null default CURRENT_TIMESTAMP,
  patient_id integer not null,
  note_id integer,
  deleted boolean default false,
  foreign key(patient_id) references patient(id),
  foreign key(note_id) references note(id)
);


create table if not exists state (
  id integer primary key,
  name string not null,
  whole boolean not null default false
);

create table if not exists change (
 id integer primary key,
 visit_id int not null,
 state_id int not null,
 tooth_num string,
 tooth_side int,
 foreign key(visit_id) references visit(id),
 foreign key(state_id) references state(id)
);
 
-- dane

insert into state (id, name, whole) 
  values (0, "próchnica", false);
insert into state (id, name, whole) 
  values (1, "ubytki niepróchnicowe", false);
insert into state (id, name, whole) 
  values (2, "wypełnienie stałe", false);
insert into state (id, name, whole) 
  values (3, "wypełnienie czasowe", false);
insert into state (id, name, whole) 
  values (4, "korona", true);
insert into state (id, name, whole) 
  values (5, "przęsło  mostu", true);
insert into state (id, name, whole) 
  values (6, "proteza", true);
insert into state (id, name, whole) 
  values (7, "implant", true);
insert into state (id, name, whole) 
  values (8, "brak zęba", true);
insert into state (id, name, whole) 
  values (9, "do usunięcia", true);
