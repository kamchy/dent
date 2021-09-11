![gabinet](gabinet.png)
Tiny webapp (local, no cloud) for dentists that stores data in SQLite.
Written as excercise when learning Go.

# Warnings

**Warning**
Web inerface is only available in Polish language.

**Warning2**
Don't use in production.

![usuwanie pacjenta](usuwanie_pacjenta.png)
# Running

To run this program, execute:
`go run main.go patient_routes.go visit_routes.go routes.go aa.db`
where `aa.db` is your datbase (SQLite) file.

# Todo
## Patients
Route /patients:
- the table should have pagination for browsing; state should remeber current page 

## Visits
 - pagination

## Dental notation block for patient and single visit
 - view of [teeth](https://en.wikipedia.org/wiki/Human_tooth) like dentists need and use([PL: system oznaczania zębów](https://pl.wikipedia.org/wiki/Systemy_oznaczania_z%C4%99b%C3%B3w]), [EN: dental notation ](https://en.wikipedia.org/wiki/Dental_notation)

## Table
Selection - maybe clickable whole row?

## Shapes
Those were designed and can be viewed in https://yqnn.github.io/svg-path-editor/
Korona:
M10 30 l 20 30 l 20 -30 l 20 30 l20 -30 v 60 h -80 Z
Most:
m 21 45 c 23 -19 42 -19 65 0 l 0 11 c -22 -16 -43 -16 -65 -1 z m -1 25 c 23 -19 42 -19 65 0 l 0 11 c -22 -16 -43 -16 -65 -1 z
Implant (plus sign):
M 45 0 h 10 v 45 h 45 v 10 h -45 v 45 h -10 v -45 h -45 v -10 h 45 Z
No toot (minus)
M 10 45 h 80 v 10 h -80 Z
Proteza (key shape):
M 30 30 L 20 50 L 30 70 L 45 70 L 45 95 L 60 95 L 60 70 L 75 70 L 85 50 L 75 30 L 65 50 L 40 50 Z
For removal:
M 92 93 L 121 64 L 197 161 L 276 79 L 301 104 L 205 181 L 311 252 L 281 282 L 195 200 L 113 283 L 83 253 L 177 179 Z
Healthy:
M 0 0 L 99 0 L 99 99 L 0 100 Z M 30 35 L 70 35 L 70 60 L 30 60 L 30 35 Z M 0 0 L 30 35 M 99 0 L 70 35 M 99 99 L 70 60 M 30 60 L 0 100

