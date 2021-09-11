
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
  whole boolean not null default false,
  color string not null,
  svgpath string not null default "M0,0h100v100h-100Z"
);

create table if not exists change (
 id integer primary key,
 visit_id int not null,
 state_id int not null,
 tooth_num string not null,
 tooth_side int not null,
 time int not null default CURRENT_TIMESTAMP,
 foreign key(visit_id) references visit(id),
 foreign key(state_id) references state(id)
);
 
-- dane

insert into state (id, name, whole, color) 
  values (0, "zdrowy", true, "#fff");
insert into state (id, name, whole, color) 
  values (1, "ubytki niepróchnicowe", false, "CadetBlue");
insert into state (id, name, whole, color) 
  values (2, "wypełnienie stałe", false, "Coral");
insert into state (id, name, whole, color) 
  values (3, "wypełnienie czasowe", false, "Peru");
insert into state (id, name, whole, color, svgpath) 
  values (4, "korona", true, "#abc", "M10 30 l 20 30 l 20 -30 l 20 30 l20 -30 v 60 h -80 Z
");
insert into state (id, name, whole, color, svgpath) 
  values (5, "przęsło  mostu", true, "#abc", "m 21 45 c 23 -19 42 -19 65 0 l 0 11 c -22 -16 -43 -16 -65 -1 z m -1 25 c 23 -19 42 -19 65 0 l 0 11 c -22 -16 -43 -16 -65 -1 z
");
insert into state (id, name, whole, color, svgpath) 
  values (6, "proteza", true, "#abc", "M 30 30 L 20 50 L 30 70 L 45 70 L 45 95 L 60 95 L 60 70 L 75 70 L 85 50 L 75 30 L 65 50 L 40 50 Z
");
insert into state (id, name, whole, color, svgpath) 
  values (7, "implant", true, "#abc", "M 45 0 h 10 v 45 h 45 v 10 h -45 v 45 h -10 v -45 h -45 v -10 h 45 Z
");
insert into state (id, name, whole, color, svgpath) 
  values (8, "brak zęba", true, "#abc", "M 10 45 h 80 v 10 h -80 Z
");
insert into state (id, name, whole, color, svgpath) 
  values (9, "do usunięcia", true, "#abc", "M 20 20 L 25 15 L 50 40 L 75 15 L 80 20 L 55 45 L 80 70 L 75 75 L 50 50 L 25 75 L 20 70 L 45 45 Z");
insert into state (id, name, whole, color) 
  values (10, "próchnica", false, "#235");
