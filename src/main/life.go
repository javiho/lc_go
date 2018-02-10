package main

import (
	"time"
	"log"
)

type Note struct{
	Text string
	Start time.Time
	End time.Time
	Id string
}

type TimeUnit int

// There doesn't appear to be anything like these constants in the time package.
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

/*
// TODO: Miksi compiler sallii tämän? Koska kyseiset constantit ovat inttejä?
var timeUnitNamesArr = []string{
	Day: "Day",
	Week: "Week",
	Month: "Month",
	Year: "Year",
}*/

// TODO: Miten voi tehdä niin, että kovakoodaa time unitit ja niiden stringit vain kerran?
var timeUnitStrings = []string{"Day", "Week", "Month", "Year"}
var timeUnitFromString = map[string]TimeUnit{
	"Day": Day,
	"Week": Week,
	"Month": Month,
	"Year": Year,
}

/*func getAllTimeUnitStrings() []string{
	tuStrings := []string{}
	for k, _ := range timeUnitFromString{
		tuStrings = append(tuStrings, k)
	}
	return tuStrings
}*/

func getStringFromTimeUnit(tu TimeUnit) string{
	for k, v := range timeUnitFromString{
		if tu == v{
			return k
		}
	}
	log.Panic("Erroneous parameter: ", tu)
	return ""
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

func (l Life) getNoteById(id string) *Note{
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

func (l *Life) addNote(newNote *Note){
	l.Notes = append(l.Notes, newNote)
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