package main

import (
	"net/http"
	"log"
	"html/template"
	//"os"
	"fmt"
	"time"
	"github.com/kjk/betterguid"
	"strconv"
)

/*
Aikavälien alku in inklusiivinen ja loppu on eksklusiivinen. Niiden ei siis koskaan pidä olla sama.
 */

type Note struct{
	Text string
	Start time.Time
	End time.Time
	Id int
}

type NoteBox struct{
	Note *Note
	Id string // Not the same thing as Note.Id
}

type TimeUnit int

//NÄMÄ VOITANEEN KORVAT time.Month ym. tyypeillä, jos viikko kuuluu niihin
const(
	Day TimeUnit = iota
	Week
	Month
	Year
)

const yyMMddLayout = "2006-01-02"

type Life struct{
	Start time.Time
	End time.Time
	Notes []*Note
}

type TimeBox struct {
	Start time.Time
	End time.Time
	//Notes []Note
	NoteBoxes []NoteBox
	Id int
}

var ResolutionUnit TimeUnit
var TheLife *Life
//var NoteBoxes []NoteBox

func main(){
	id := betterguid.New()
	fmt.Println("id: ", id)
	initializeData()
	fmt.Println("data initialized")
	http.HandleFunc("/", SendCalendar)
	http.HandleFunc("/note_changed", ChangeAndSendCalendar)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initializeData() {
	ResolutionUnit = Month

	notes := []*Note{
		&Note{"Note number 1", time.Date(2017, time.February, 15, 0,0,0,0,time.UTC),
		time.Date(2017, time.April, 1, 0,0,0,0,time.UTC), 1},
		&Note{"Note number 2", time.Date(2018, time.November, 1, 0,0,0,0,time.UTC),
			time.Date(2018, time.November, 15, 0,0,0,0,time.UTC), 2},
	}
	TheLife = &Life{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
	notes}
}

func SendCalendar(w http.ResponseWriter, r *http.Request){
	timeBoxes := createTimeBoxes(TheLife, ResolutionUnit)

	/*for _, timeBox := range timeBoxes{
		for _, noteBox := range timeBox.NoteBoxes{
			NoteBoxes = append(NoteBoxes, noteBox)
		}
	}*/

	template, err := template.ParseFiles("src/main/main_view.html")
	if err != nil{
		log.Print("template parsing error: ", err)
	}
	err = template.Execute(w, timeBoxes)
	if err != nil{
		log.Print("template executing error: ", err)
	}
}

func ChangeAndSendCalendar(w http.ResponseWriter, r *http.Request){
	//TODO: VALIDOINTI
	r.ParseForm()
	fmt.Println(r.Form)
	noteId := r.Form["note-id"][0]
	noteText := r.Form["note-text"][0]
	start := r.Form["note-start"][0]
	end := r.Form["note-end"][0]

	noteIdString, err := strconv.Atoi(noteId)
	if err != nil{
		log.Panic("erroneous note id")
	}
	startDate, err := time.Parse(yyMMddLayout, start)
	if err != nil{
		log.Panic("erroneous start date")
	}
	endDate, err := time.Parse(yyMMddLayout, end)
	if err != nil{
		log.Panic("erroneous end date")
	}

	if endDate.Before(startDate) || endDate.Equal(startDate){
		//TODO: tästä pitäisi tämän sijaan ilmoittaa käyttäjälle, jtota hän voi korjata arvot.
		log.Panic("erroneous date values")
	}

	note := TheLife.getNoteById(noteIdString)
	note.Text = noteText
	note.Start = startDate
	note.End = endDate

	SendCalendar(w, r)
}

func createTimeBoxes(life *Life, resolutionUnit TimeUnit) []TimeBox{
	/*
	Otetaan elämän alku ja loppukohdat ja lasketaan näytettävät alku ja loppukohdat

	Lähdetään liikkeelle elämän ensimmäisestä päivästä
	otetaan sen resoluutioaikaväli ja lisätään se listaan
	liikutaan eteenpäin resoluutioyksikön verran ja toistetaan edellinen
	tätä jatketaan, kunnes elämän loppupäivän ssältävä resoluutioaikaväli on lisätty listaan
	 */
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

func getNotesByInterval(life *Life, start time.Time, end time.Time) []*Note{
	/*
	Pre-condition start < end (ENTÄ JOS ON SAMA?). life.start < life.end
	 */
	/*
	käy läpi kaikki notet lifessä
		jos ne > is ja ns < ie
	 */
 	var notesInInterval []*Note
 	for _, n := range life.Notes{
 		if n.End.After(start) && n.Start.Before(end){
			notesInInterval = append(notesInInterval, n)
		}
	}
	//fmt.Println("getting notes by interval ", start, " to ", end, ". Notes returned: ", len(notesInInterval))
	return notesInInterval
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

func (l Life) getNoteById(id int) *Note{
	// PITÄISI MIELUUMMIN Notes-attribuutin olla map kuin käyttää tämmöistä luuppia. ID:iden uniikkiuskin olisi taattu.
	for _, n := range l.Notes{
		if n.Id == id{
			return n
		}
	}
	log.Panic("Note not found of id ", id)
	return &Note{}
}

func (l *Life) replaceNote(newNote *Note) {
	for i, n := range l.Notes{
		if n.Id == newNote.Id{
			l.Notes[i] = newNote
			break
		}
	}
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

func (n Note) StartAsString() string{
	/*
	Returns start as string which can be inserted to HTML date input element.
	 */
	return n.Start.Format(yyMMddLayout)
}

func (n Note) EndAsString() string{
	return n.End.Format(yyMMddLayout)
}