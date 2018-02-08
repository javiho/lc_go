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
	Day TimeUnit = iota + 1
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