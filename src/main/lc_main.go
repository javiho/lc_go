package main

import (
	"net/http"
	"log"
	"html/template"
	//"os"
	"fmt"
	"time"
	"github.com/kjk/betterguid"
)

/*
Aikavälien alku in inklusiivinen ja loppu on eksklusiivinen. Niiden ei siis koskaan pidä olla sama.
 */

 /*
 uusi rakenne:
 pää lol

  */



type NoteBox struct{
	Note *Note
	Id string // Not the same thing as Note.Id
}

type TimeBox struct {
	Start time.Time
	End time.Time
	//Notes []Note
	NoteBoxes []NoteBox
	Id int
}

type LcPageVariables struct{
	TimeBoxes []TimeBox
	Notes []*Note
	LifeStart string
	LifeEnd string
	ResolutionUnit string
	AllResolutionUnits []string
	TrueTime time.Time
	Errors []string
}

var ResolutionUnit TimeUnit
var TheLife *Life
var AllCategories []*Category
//var NoteBoxes []NoteBox

func main(){
	id := betterguid.New()
	fmt.Println("id: ", id)
	initializeData()
	fmt.Println("data initialized")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./src/main/static"))))
	http.HandleFunc("/favicon.ico", SendNothing)
	http.HandleFunc("/", SendDefaultView)
	http.HandleFunc("/note_changed", ChangeAndSendCalendar)
	http.HandleFunc("/note_added", AddNoteAndSendCalendar)
	http.HandleFunc("/lc_options_changed", ChangeLcOptionsAndSendCalendar)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initializeData() {
	ResolutionUnit = Week

	catEgory := &Category{"Cat1"}
	dogCategory := &Category{"Dog1"}
	AllCategories = []*Category{catEgory, dogCategory}
	notes := []*Note{
		&Note{"Note number 1", time.Date(2017, time.February, 15, 0,0,0,0,time.UTC),
		time.Date(2017, time.April, 1, 0,0,0,0,time.UTC), []*Category{catEgory}, "hcNote1"},
		&Note{"Note number 2", time.Date(2018, time.November, 1, 0,0,0,0,time.UTC),
			time.Date(2018, time.November, 15, 0,0,0,0,time.UTC), []*Category{dogCategory}, "hcNote2"},
	}
	TheLife = &Life{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
	notes}
}

func SendCalendar(w http.ResponseWriter, r *http.Request, errorsForClient []string){
	//TODO: Miksi tätä funktiota kutsutaan kahdesti joka requestilla?
	fmt.Println("SendCalendar called")
	timeBoxes := createTimeBoxes(TheLife, ResolutionUnit)
	lcPageVariables := LcPageVariables{timeBoxes,TheLife.Notes, TheLife.Start.Format(yyMMddLayout),
		TheLife.End.Format(yyMMddLayout),getStringFromTimeUnit(ResolutionUnit), timeUnitStrings,
		time.Now(), errorsForClient}
	fmt.Println(errorsForClient)
	//fmt.Println(getStringFromTimeUnit(ResolutionUnit))
	/*for _, timeBox := range timeBoxes{
		for _, noteBox := range timeBox.NoteBoxes{
			NoteBoxes = append(NoteBoxes, noteBox)
		}
	}*/

	templ, err := template.ParseFiles("src/main/main_view.html")
	if err != nil{
		log.Print("template parsing error: ", err)
	}
	err = templ.Execute(w, lcPageVariables)
	if err != nil{
		log.Print("template executing error: ", err)
	}
}

func SendDefaultView(w http.ResponseWriter, r *http.Request){
	fmt.Println("SendDefaultView called")
	SendCalendar(w, r, []string{})
}

func ChangeAndSendCalendar(w http.ResponseWriter, r *http.Request){
	//TODO: VALIDOINTI
	r.ParseForm()
	fmt.Println(r.Form)
	noteId := r.Form["note-id"][0]
	note := TheLife.getNoteById(noteId)

	_, isActionSave := r.Form["save-submit"]
	_, isActionDelete := r.Form["delete-submit"]
	if isActionSave{
		// TODO: voisi olla omassa functiossa
		fmt.Println("Trying to save note")
		noteText := r.Form["note-text"][0]
		//startDate := r.Form["note-start"][0]
		//endDate := r.Form["note-end"][0]

		dates, lesUnparsebles, errorMessages := parseStartAndEndDates(r.Form["note-start"][0], r.Form["note-end"][0])
		if lesUnparsebles{
			log.Println("parse error")
			SendCalendar(w, r, errorMessages)
			return
		}
		startDate := dates[0]
		endDate := dates[1]
		/*
		var errorsForClient []string
		startDate, err := time.Parse(yyMMddLayout, start)
		if err != nil{
			log.Panic("erroneous start date")
		}
		endDate, err := time.Parse(yyMMddLayout, end)
		if err != nil{
			log.Panic("erroneous end date")
		}
		*/
		if endDate.Before(startDate) || endDate.Equal(startDate){
			//TODO: tästä pitäisi tämän sijaan ilmoittaa käyttäjälle, jtota hän voi korjata arvot.
			log.Println("erroneous date values")
			SendCalendar(w, r, []string{"End date not after start date."})
			return
		}

		note.Text = noteText
		note.Start = startDate
		note.End = endDate
	} else if isActionDelete{
		fmt.Println("Trying to delete note")
		TheLife.deleteNote(note)
	} else{
		log.Panic("erroneous submit handling")
	}

	SendCalendar(w, r, []string{})
}

func AddNoteAndSendCalendar(w http.ResponseWriter, r *http.Request){
	//TODO:
	//TODO: VALIDOINTI
	r.ParseForm()
	fmt.Println(r.Form)
	noteText := r.Form["note-text"][0]
	//start := r.Form["note-start"][0]
	//end := r.Form["note-end"][0]
	dates, lesUnparsebles, errorMessages := parseStartAndEndDates(r.Form["note-start"][0], r.Form["note-end"][0])
	if lesUnparsebles{
		log.Println("parse error")
		SendCalendar(w, r, errorMessages)
		return
	}
	startDate := dates[0]
	endDate := dates[1]

	/*var errorsForClient []string
	startDate, err1 := time.Parse(yyMMddLayout, start)
	_, parseErrorOccurred1 := err1.(*time.ParseError)
	if parseErrorOccurred1{
		errorsForClient = append(errorsForClient, "Erroneous start date.")
	}
	endDate, err2 := time.Parse(yyMMddLayout, end)
	_, parseErrorOccurred2 := err2.(*time.ParseError)
	if parseErrorOccurred2{
		errorsForClient = append(errorsForClient, "Erroneous end date.")
	}
	if !parseErrorOccurred1 || !parseErrorOccurred2{
		log.Println("parse error")
		SendCalendar(w, r, errorsForClient)
		return
	}
	if err1 != nil || err2 != nil{
		log.Panic("Unexpected error")
	}*/

	if endDate.Before(startDate) || endDate.Equal(startDate){
		log.Println("erroneous date values")
		SendCalendar(w, r, []string{"End date not after start date."})
		return
	}

	note := Note{noteText, startDate, endDate, []*Category{}, betterguid.New()}
	TheLife.addNote(&note)
	fmt.Println("note added with id: ", note.Id)
	SendCalendar(w, r, []string{})
}

func ChangeLcOptionsAndSendCalendar(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fmt.Println(r.Form)
	resolutionUnitString := r.Form["resolution-unit"][0]
	//lifeStart, err1 := time.Parse(yyMMddLayout, r.Form["life-start"][0])
	//_, parseErrorOccurred1 := err1.(*time.ParseError) // TODO: jos * ottaa pois, niin editorin mukaan virhe. Miksi?
	//var errorsForClient []string
	// TODO: miten virheidenkäsittely pitäisi tehdä tässä?
	dates, lesUnparsables, errorMessages := parseStartAndEndDates(r.Form["life-start"][0], r.Form["life-end"][0])
	if lesUnparsables {
		log.Println("parse error")
		SendCalendar(w, r, errorMessages)
		return
	}
	lifeStart := dates[0]
	lifeEnd := dates[1]
	/*
	if parseErrorOccurred1{
		errorsForClient = append(errorsForClient, "Erroneous start date value.")
	}
	if !parseErrorOccurred1 && err1 != nil{
		log.Panic(err1)
	}

	lifeEnd, err2 := time.Parse(yyMMddLayout, r.Form["life-end"][0])
	_, parseErrorOccurred2 := err2.(*time.ParseError) // TODO: jos * ottaa pois, niin editorin mukaan virhe. Miksi?
	if parseErrorOccurred2{
		errorsForClient = append(errorsForClient, "Erroneous end date value.")
	}
	if !parseErrorOccurred2 && err2 != nil{
		log.Panic(err2)
	}

	if parseErrorOccurred1 || parseErrorOccurred2{
		log.Println("erroneous date values: parse error")
		SendCalendar(w, r, errorsForClient)
		return
	}
	*/
	if lifeEnd.Before(lifeStart) || lifeEnd.Equal(lifeStart){
		//TODO: tästä pitäisi tämän sijaan ilmoittaa käyttäjälle, jtota hän voi korjata arvot.
		log.Println("erroneous date values: erroneous chronology")
		SendCalendar(w, r, []string{"Start date must be before end date."})
		return
	}

	resolutionUnit := timeUnitFromString[resolutionUnitString] // TODO: entä jos on virheellinen stringi?
	//fmt.Println("new resolution unit:", resolutionUnit)
	ResolutionUnit = resolutionUnit
	TheLife.Start = lifeStart
	TheLife.End = lifeEnd
	SendCalendar(w, r, []string{})
}

func SendNothing(w http.ResponseWriter, r *http.Request){
	return
}

// ROUTING ENDS HERE //

func createTimeBoxes(life *Life, resolutionUnit TimeUnit) []TimeBox{
	/*
	Otetaan elämän alku ja loppukohdat ja lasketaan näytettävät alku ja loppukohdat

	Lähdetään liikkeelle elämän ensimmäisestä päivästä
	otetaan sen resoluutioaikaväli ja lisätään se listaan
	liikutaan eteenpäin resoluutioyksikön verran ja toistetaan edellinen
	tätä jatketaan, kunnes elämän loppupäivän ssältävä resoluutioaikaväli on lisätty listaan
	 */
 	if life.Start.After(life.End) || life.Start.Equal(life.End){
 		log.Panic("Life start not before life end.")
	}
	var timeBoxes []TimeBox
	counter := 0
 	adjustedLifeStart := getFirstDateOfTimeUnit(life.Start, resolutionUnit)
 	for t := adjustedLifeStart; true; t = addTimeUnit(t, resolutionUnit){
		tPlusResolutionUnit := addTimeUnit(t, resolutionUnit)
 		newTimeBox := TimeBox{t,tPlusResolutionUnit,getNoteBoxesByInterval(life, t, tPlusResolutionUnit), counter}
		timeBoxes = append(timeBoxes, newTimeBox)
		// Life end time is exclusive.
 		if tPlusResolutionUnit.After(life.End) || tPlusResolutionUnit.Equal(life.End){
 			break
		}
		counter += 1
	}
	return timeBoxes
}

func getFirstDateOfTimeUnit(date time.Time, timeUnit TimeUnit) time.Time{
	y, m, d := date.Date()
	if timeUnit == Month{
		return time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	}
	if timeUnit == Year{
		return time.Date(y, time.January, 1,0, 0, 0, 0, time.UTC)
	}
	if timeUnit == Week{
		return getFirstDateOfWeek(date)
	}
	if timeUnit == Day{
		return time.Date(y, m, d, 0,0,0,0, time.UTC)
	}
	log.Panic("Time unit doesn't exist.")
	return date
}

func getFirstDateOfWeek(dateInWeek time.Time) time.Time{
	mondayCandidate := dateInWeek
	for mondayCandidate.Weekday() != time.Monday {
		mondayCandidate = mondayCandidate.AddDate(0, 0, -1)
	}
	return mondayCandidate
}

func addTimeUnit(time time.Time, timeUnit TimeUnit) time.Time{
	if timeUnit == Day{
		return time.AddDate(0, 0, 1)
	}
	if timeUnit == Month{
		return time.AddDate(0, 1, 0)
	}
	if timeUnit == Year{
		return time.AddDate(1, 0, 0)
	}
	if timeUnit == Week{
		return time.AddDate(0, 0, 7)
	}
	log.Panic("Time unit doesn't exist.")
	return time
}

func createNoteBoxes(notes []*Note) []NoteBox{
	var noteBoxes []NoteBox
	for _, n := range notes{
		noteBox := NoteBox{n, betterguid.New()}
		noteBoxes = append(noteBoxes, noteBox)
	}
	return noteBoxes
}

func getNoteBoxesByInterval(life *Life, start time.Time, end time.Time) []NoteBox{
	notes := getNotesByInterval(life, start, end)
	noteBoxes := createNoteBoxes(notes)
	return noteBoxes
}

func (tb TimeBox) StartAsString() string{
	/*
	Returns start as string which can be inserted to HTML date input element.
	 */
	return tb.Start.Format(yyMMddLayout)
}

func (tb TimeBox) EndAsString() string{
	return tb.End.Format(yyMMddLayout)
}

func parseStartAndEndDates(d1string string, d2string string) ([]time.Time, bool, []string){
	/*
	Returns parsed dates, if parse errors occurred, the error messages for client.
	 */
	d1, err1 := time.Parse(yyMMddLayout, d1string)
	d2, err2 := time.Parse(yyMMddLayout, d2string)
	var errorMessages []string
	_, parseErrorOccurred1 := err1.(*time.ParseError) // TODO: jos * ottaa pois, niin editorin mukaan virhe. Miksi?
	_, parseErrorOccurred2 := err2.(*time.ParseError) // TODO: jos * ottaa pois, niin editorin mukaan virhe. Miksi?
	if parseErrorOccurred1{
		errorMessages = append(errorMessages, "Parse error of start date")
	}
	if parseErrorOccurred2{
		errorMessages = append(errorMessages, "Parse error of end date")
	}
	if err1 != nil && !parseErrorOccurred1{
		log.Panic("bug")
	}
	if err2 != nil && !parseErrorOccurred2{
		log.Panic("bug")
	}
	parseErrorsOccurred := parseErrorOccurred1 || parseErrorOccurred2
	return []time.Time{d1, d2}, parseErrorsOccurred, errorMessages
}