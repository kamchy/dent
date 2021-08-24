
create table if not exists patient (
  id integer primary key,
  name string not null,
  surname string not null,
  birthdate date not null,
  note_id integer,
  foreign key(note_id) references note(id) 
);

create table if not exists note (
  id integer primary key,
  text string not null
);

create table if not exists visit (
  id integer primary key,
  vdatetime datetime not null default CURRENT_TIMESTAMP,
  patient_id integer not null,
  note_id integer,
  foreign key(patient_id) references patient(id),
  foreign key(note_id) references note(id)
);
