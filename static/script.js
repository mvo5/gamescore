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
            url: url_base + "/api/1/game",
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
            type: 'put',
            url: url_base + "/api/1/game",
            dataType: 'json',
            data: JSON.stringify(state_change)
        });
    }

    function doTimeout() {
        var state_change = {
            ToggleTimeout: true
        }
        $.ajax({
            type: 'put',
            url: url_base + "/api/1/game",
            dataType: 'json',
            data: JSON.stringify(state_change)
        });
    }

    function status() {
        $.getJSON(url_base + "/api/1/game", function(json) {
            $('#time').html('<h2>'+json.TimeStr+'</h2>');
            $('#name_team1').html('<h2>'+json.Team1.Name+'</h2>');
            $('#score_team1').html('<h2>'+json.Team1.Goals+'</h2>');
            $('#name_team2').html('<h2>'+json.Team2.Name+'</h2>');
            $('#score_team2').html('<h2>'+json.Team2.Goals+'</h2>');

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
