"use strict";
var timeBoxPastClass = "tb-past";
var timeBoxFutureClass = "tb-future";
var jsErrorInfoIconClass = '.js-error-info-icon';

var editNoteForm;
var noteSelectedTextInput;
var lcOptionsForm;

const MIN_YEAR = 1905;
const MAX_YEAR = 2195;

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

function toDate(dateStr) {
    //console.log("dateStr:");
    //console.log(dateStr);
    //const [day, month, year] = dateStr.split("-"); // TODO: Mikä tämä syntaksi on?
    //return new Date(year, month - 1, day);
    var dateComponents = dateStr.split("-");
    var date = new Date(dateComponents[0], dateComponents[1] - 1, dateComponents[2]);
    //console.log(date);
    return date;
}

function validateLifeOptionsForm(){
    console.log("form input changed!")
    var lifeStartInput = $('#life-start-input');
    var lifeEndInput = $('#life-end-input');
    lifeStartInput[0].setCustomValidity("");
    lifeEndInput[0].setCustomValidity("");
    //console.log("end check validity:", lifeEndInput[0].checkValidity());
    // lifeEndInput.attr('value') returns the original (not selected by user) value. TODO: MIKSI?
    var lifeStartDate = toDate(lifeStartInput[0].value);
    var lifeEndDate = toDate(lifeEndInput[0].value);
    var mess = `Minimum year is ${MIN_YEAR} and maximum year is ${MAX_YEAR}.`;
    if(lifeStartDate.getUTCFullYear() < MIN_YEAR || lifeStartDate.getUTCFullYear() > MAX_YEAR){
        //console.log(`Minimum year is ${MIN_YEAR} and maximum year is ${MAX_YEAR}.`);
        lifeStartInput[0].setCustomValidity(mess);
    }
    if(lifeEndDate.getUTCFullYear() > MAX_YEAR || lifeEndDate.getUTCFullYear() < MIN_YEAR){
        lifeEndInput[0].setCustomValidity(mess);
    }
    if(lifeStartDate >= lifeEndDate){
        mess = "Birth date must be before estimated death date.";
        console.log("Birth date must be before estimated death date.");
        lifeStartInput[0].setCustomValidity(mess);
        lifeEndInput[0].setCustomValidity(mess);
    }
    if(lifeStartInput[0].checkValidity() && lifeEndInput[0].checkValidity()){
        //console.log("submitting life option form!");
        lcOptionsForm.closest("form")[0].submit();
    }else{
        var closestInfoIcon = lifeStartInput.next(jsErrorInfoIconClass);
        if(!lifeStartInput[0].checkValidity()){
            closestInfoIcon.css("visibility", "visible");
            closestInfoIcon.attr("title", lifeStartInput[0].validationMessage);
        }else{
            closestInfoIcon.css("visibility", "hidden");
        }
        closestInfoIcon = lifeEndInput.next(jsErrorInfoIconClass);
        if(!lifeEndInput[0].checkValidity()){
            closestInfoIcon.attr("title", lifeEndInput[0].validationMessage);
            closestInfoIcon.css("visibility", "visible");
        }else{
            closestInfoIcon.css("visibility", "hidden");
        }
    }
    // TODO: kommentti:
    // No niin elikkä html:n validaatiokuplat näkyvät vain jos submitataab submit napun kautta (ei submit() function)
    // ks. https://stackoverflow.com/questions/16707743/html5-validation-when-the-input-type-is-not-submit/31741546
    // Eli voi tehdä piilotetu nsubmit nappulan tai säätää omat validaatio messagen displayaamiset.
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

    $(jsErrorInfoIconClass).css("visibility", "hidden");

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

    lcOptionsForm = $("#lc-options-form");
    lcOptionsForm.change(validateLifeOptionsForm);

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

    // Set validation things
    var lifeStartInput = $('#life-start-input')[0];
    var lifeEndInput = $('#life-end-input')[0];

}

console.log("executing script");
$(initialize);