var url_base = "http://localhost:8080"

$(document).ready(function(){

    var create = function() {
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

    var status = function() {
        $.getJSON(url_base + "/api/1/game", function(json) {
            $('#time').html('<h2>'+json.TimeLeft+'</h2>');
            $('#name_team1').html('<h2>'+json.Team1.Name+'</h2>');
            $('#score_team1').html('<h2>'+json.Team1.Goals+'</h2>');
            $('#name_team2').html('<h2>'+json.Team2.Name+'</h2>');
            $('#score_team2').html('<h2>'+json.Team2.Goals+'</h2>');
        })
    }

    $('#create').click(create);
    window.setInterval(status, 100);
});
