package amqpcast

const indexTemplate = `
<!DOCTYPE html>
<html>

<head>
<title>amqpcast</title>
<meta charset="utf-8">
<script src="http://code.jquery.com/jquery-2.0.3.min.js"></script>
<style>
#messages {
  font-family: monospace;
}
</style>
</head>

<body>

<h1>amqpcast</h1>

<div id="messages"></div>

<script>
$(function() {
    var ws = new WebSocket("ws://" + window.location.host + "/ws")

    ws.onopen = function() {
        console.log("websocket open");
    };

    ws.onmessage = function(e) {
        $('#messages').prepend($("<p/>").append(new Date() + ' ' + e.data));
    };

    ws.onclose = function(e) {
        console.log("closed");
        console.log(e);
    };

    ws.onerror = function(e) {
        console.log("error");
        console.log(e);
    }
});
</script>

</body>
</html>
`
