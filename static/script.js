$(document).ready(function(){
    var url_base = "http://localhost:8080"
    var audioPlayed = true

    function create() {
        var new_game = {
            timeleft: parseInt($("#input_time").val()),
            team1: {
                name: $("#input_team1_name").val()
            },
            team2: {
                name: $("#input_team2_name").val()
            },
            running: true
        }
        $.ajax({
            type: 'post',
            url: url_base + "/api/1/create",
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
            url: url_base + "/api/1/changeState",
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
            url: url_base + "/api/1/changeState",
            dataType: 'json',
            data: JSON.stringify(state_change)
        });
    }

    function status() {
        $.getJSON(url_base + "/api/1/status", function(json) {
            $('#time').text(json.TimeStr);
            $('#name_team1').text(json.Team1.Name);
            $('#score_team1').text(json.Team1.Goals);
            $('#name_team2').text(json.Team2.Name);
            $('#score_team2').text(json.Team2.Goals);

            // check if a new game has started
            if (json.Running == true) {
                audioPlayed = false
            }

            // play alarm
            if (json.TimeLeft <= 0 && !audioPlayed) {
                // from http://soundbible.com/1577-Siren-Noise.html
                // licensed public domain
                var audio = new Audio('Siren_Noise-KevanGC-1337458893.wav');
                audio.play();
                audioPlayed = true
            }
        })
    }

    // the callbacks
    $('#create').click(create);
    $('#decTeam1').click(function() {scoreTeam('team1', -1)})
    $('#incTeam1').click(function() {scoreTeam('team1', +1)})
    $('#decTeam2').click(function() {scoreTeam('team2', -1)})
    $('#incTeam2').click(function() {scoreTeam('team2', +1)})
    $('#timeout').click(function() {doTimeout()})

    window.setInterval(status, 100);
});
