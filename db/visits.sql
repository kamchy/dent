with 
vn(id, dt, pid, t) as (
  select
        visit.id as id,
        visit.vdatetime as dt, 
        visit.patient_id as pid, 
        note.text as t
  from visit, note 
  where visit.note_id = note.id),
pn(name, surname, pid, bdate, t) as (
  select 
        patient.name as name, 
        patient.surname as surname, 
        patient.id as pid, 
        patient.birthdate as bdate,
        note.text  as t
  from patient, note 
  where patient.note_id = note.id)
select vn.id, vn.dt, vn.t, pn.pid, pn.name, pn.surname, pn.bdate, pn.t from vn, pn where vn.pid = pn.pid;

