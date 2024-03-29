$(document).ready(function(){
    var audioPlayed = true;

    function getGameData() {
        return {
            // ui gets time in minutes, api expects nanoseconds so 1000000000
            timeLeft: parseInt($("#input_time").val()) * 60 * 1000000000,
            team1: {
                name: $("#input_team1_name").val()
            },
            team2: {
                name: $("#input_team2_name").val()
            }
        }
    }
    
    function start() {
        var new_game = getGameData()
        new_game.running = true
        $.ajax({
            type: 'post',
            url: "/api/1/create",
            dataType: 'json',
            data: JSON.stringify(new_game)
        });
    }
    
    function create() {
        var new_game = getGameData()
        $.ajax({
            type: 'post',
            url: "/api/1/create",
            dataType: 'json',
            data: JSON.stringify(new_game)
        });
    };

    function scoreTeam(team, scoreChange) {
        var state_change = {
            TeamId: team,
            ScoreChange: scoreChange
        }
        $.ajax({
            type: 'post',
            url: "/api/1/changeState",
            dataType: 'json',
            data: JSON.stringify(state_change)
        });
    }

    function doTimeout() {
        var state_change = {
            ToggleTimeout: true
        }
        $.ajax({
            type: 'post',
            url: "/api/1/changeState",
            dataType: 'json',
            data: JSON.stringify(state_change)
        });
    }

    function status() {
        $.getJSON("/api/1/status", function(json) {
            $('#time').text(json.TimeStr);
            $('#name_team1').text(json.Team1.Name);
            $('#score_team1').text(json.Team1.Goals);
            $('#name_team2').text(json.Team2.Name);
            $('#score_team2').text(json.Team2.Goals);

            // check if a new game has started
            if (json.Running == true) {
                audioPlayed = false
                $('#timeout').text("Pause time")
            } else {
                $('#timeout').text("Start time")
            }

            // play alarm
            if (json.TimeLeft <= 0 && !audioPlayed) {
                playSound()
                audioPlayed = true
            }
        })
    }

    function playSound() {
        var audio = new Audio($('#select_sound').val());
        audio.play();
    }

    function changeSides() {
        var team1 = $('#input_team1_name').val();
        var team2 = $('#input_team2_name').val();
        var scoreTeam1 = $('#score_team1').text();
        var scoreTeam2 = $('#score_team2').text();
        // Convenient function for second half:
        // change sides and create new game with swaped names/goals
        $('#input_team1_name').val(team2);
        $('#input_team2_name').val(team1);
        // we need to switch sides in a single POST to prevent races
        var switchSides = getGameData();
        switchSides["team1"].Goals = parseInt(scoreTeam2);
        switchSides["team2"].Goals = parseInt(scoreTeam1);
        $.ajax({
            type: 'post',
            url: "/api/1/create",
            dataType: 'json',
            data: JSON.stringify(switchSides)
        });
    }

    // the callbacks
    $('#start').click(start);
    $('#create').click(create);
    $('#decTeam1').click(function() {scoreTeam('team1', -1)})
    $('#incTeam1').click(function() {scoreTeam('team1', +1)})
    $('#decTeam2').click(function() {scoreTeam('team2', -1)})
    $('#incTeam2').click(function() {scoreTeam('team2', +1)})
    $('#timeout').click(function() {doTimeout()})
    $('#swap').click(function() {changeSides()})
    $('#testSound').click(playSound);

    // keyboard shortcuts
    $('[data-toggle="tooltip"]').tooltip();
    $(document).keydown(function (event) {
        console.log(event);
        switch (event.keyCode) {
            // left arrow
        case 37:
            if (event.shiftKey == false) {
                scoreTeam('team1', +1);
            } else {
                scoreTeam('team1', -1);
            }
            event.preventDefault()
            break;
            // right arrow
        case 39:
            if (event.shiftKey == false) {
                scoreTeam('team2', +1);
            } else {
                scoreTeam('team2', -1);
            }
            event.preventDefault()
            break;
            // down arrow
        case 40:
            doTimeout();
            event.preventDefault()
            break;
        }
    });
    
    window.setInterval(status, 100);
});
