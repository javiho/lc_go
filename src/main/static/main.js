"use strict";
var timeBoxPastClass = "tb-past";
var timeBoxFutureClass = "tb-future";

var editNoteForm;
var noteSelectedTextInput;

function utcDateFromString(dateString){
    //Pre-condition: dateString must be yyyy-MM-dd.
    return new Date(dateString + 'T00:00:00.000Z'); // lol
}

function setPastFutureClasses(){
    var timeBoxes = $('.js-time-box');
    //console.log(timeBoxes);
    timeBoxes.each(function(){
        var tb = $(this);
        var tbEndDate = utcDateFromString(tb.data('end'));
        //Date-objecteja voi vertailla toimivasti, vaikka ne ovat eri aikavyöhykkeellä.
        var now = new Date();
        // If end date is 18.2. and now is 18.2. 13:00, isPast == true.
        // This is ok, because end date is not included in the time span.
        var isPast = now > tbEndDate;
        //console.log("now: " + now + " end: " + tbEndDate + " past: " + isPast);
        var pastFutureStripe = tb.find('.js-tb-past-future-stripe');
        if(isPast){
            pastFutureStripe.addClass(timeBoxPastClass);
            pastFutureStripe.removeClass(timeBoxFutureClass);
            /*if(tb.hasClass(timeBoxFutureClass)){
                tb.removeClass(timeBoxFutureClass);
            }*/
        }else{
            pastFutureStripe.addClass(timeBoxFutureClass);
            pastFutureStripe.removeClass(timeBoxPastClass);
            /*if(tb.hasClass(timeBoxPastClass)){
                tb.removeClass(timeBoxPastClass);
            }*/
        }
    })
}

function editNoteButtonClicked(e){
    console.log("edit note button click'd!")
    var clickedButton = $(e.target);
    var noteBoxDiv = clickedButton.closest("div").filter(".js-note-box");
    console.assert(noteBoxDiv !== undefined, "noteBoxDiv undefined");
    var start = noteBoxDiv.data("note-start");
    var end = noteBoxDiv.data("note-end");
    var text = noteBoxDiv.data("note-text");
    var id = noteBoxDiv.data("note-id");
    console.assert([start, end, text, id].every(x => x !== undefined), [start, end, text, id]);
    $('#note-selected-text').attr('value', text);
    $('#note-selected-start').attr('value', start);
    $('#note-selected-end').attr('value', end);
    $('#selected-note-id').attr('value', id);
}

function noteBoxClicked(e){
    console.log("time box click'd!");
    //OLTAVA SAMANTYYPPINEN FUNKTIO KUIN editNoteButtonClicked ja KUTSUTTAVA SITÄÖ

}

function timeBoxClicked(e){
    console.log("time box click'd!")
    var timeBox = $(e.target);
    var noteBoxes = timeBox.find('div.js-note-box');
    var noteIds = [];
    //TODO: () => noteIds.push($(this).data.noteId) EI TOIMI, MUTTA MIKSI? MITEN THIS TOIMII?
    noteBoxes.each(function() {
        noteIds.push($(this).data('note-id'));
    });
    var allNoteReps = $('.js-note-rep');
    allNoteReps.each(function(){
        var self = $(this); // TODO: miten this toimii?
        var noteId = self.data('note-id');
        if(noteIds.includes(noteId)){
            self.show(400);
        }else{
            self.hide();
        }
    });
    //TODO: uusi note -lomakkeen päivät asetettava
    var tbStart = timeBox.data('start');
    var tbEnd = timeBox.data('end');
    $('#new-note-start').attr('value', tbStart);
    $('#new-note-end').attr('value', tbEnd);
}

function noteRepClicked(e){
    console.log("note-rep click'd!");

    var noteRep = $(e.target);
    console.assert(noteRep !== undefined, "noteRep undefined");
    var start = noteRep.data("note-start");
    var end = noteRep.data("note-end");
    var text = noteRep.data("note-text");
    var id = noteRep.data("note-id");
    console.assert([start, end, text, id].every(x => x !== undefined), [start, end, text, id]);
    $('#note-selected-text').attr('value', text);
    $('#note-selected-start').attr('value', start);
    $('#note-selected-end').attr('value', end);
    $('#selected-note-id').attr('value', id);
}

function initialize(){
    console.log("intialize call'd!");
    editNoteForm = $('#edit-note-form');

    // Set default selected life calendar option.
    var defaultSelectedOptionString = $('#resolution-unit-select').data('default-selected');
    console.log($('#resolution-unit-select').data("default-selected"));
    console.assert(defaultSelectedOptionString !== undefined);
    var defaultSelectedOption = $('#resolution-unit-select').find(`option[value="${defaultSelectedOptionString}"]`);
    defaultSelectedOption.attr('selected', true);

    $(document).click(function(e){
        var eventTargetJQuery = $(e.target)
        if(eventTargetJQuery.is(".js-note-box")){
            editNoteButtonClicked(e);
        }
        if(eventTargetJQuery.is(".js-time-box")){
            timeBoxClicked(e);
        }
        if(eventTargetJQuery.is(".js-note-rep")){
            noteRepClicked(e);
        }
    });

    var lcOptionsForm = $("#lc-options-form")
    //var resolutionSelect = $("#resolution-unit-select");
    lcOptionsForm.change(function() {
        console.log("form input changed!")
        //TODO: ensin pitää validoida - vai pitääkö?
        // TODO: requider huolehtii, että päivämäärän arvo on validi (mutta vain jos se on ylipäätään laitettu)
        // TODO: lisäksi voisi validoida että loppu on alun jälkeen
        // TODO: mutta pitäisi ainakin antaa mahdollisuus kirjoittaa loppuun vuosi!
        lcOptionsForm.closest("form")[0].submit();
    });

    setPastFutureClasses();

    // Display error messages only if there are any.
    var serverErrorMessages = $('.js-server-error-message');
    if(serverErrorMessages.length > 0){
        var allMessages = [];
        var errorMessagesDiv = $('#server-error-messages-div');
        errorMessagesDiv.find('p').each(function(){
            allMessages.push($(this).text());
        });
        alert(allMessages.join('\n'));
    }
}

console.log("executing script");
$(initialize);