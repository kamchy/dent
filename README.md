# Running

To run this program, execute:
`go go run main.go  db/sample.db`

# Todo
## Patients
Route /patients:
- button next perhaps at the top?
- if no persons, be ready to post entered data
- if there are persons, always select first one on /patients
- the table should have pagination for browsing; state should remeber current page 
- add search field that loads only persons with given string
- add pane with all visits of selected person (view only)
- person should have  "visits" link that goes to /patients/:pat/visits

## Visits
Route /patients/:pat/visits
 - list of all visits of given patient
 - pagination
 - each line has "show" that displays visit details
 - button "new visit" at the top?
 
