document.getElementById('start').onclick = function() {
    var username = document.getElementById('username').value;
    var channel = document.getElementById('channel').value;

    // Connect to the WebSocket server
    var socket = new WebSocket("ws://localhost:8101/ws");

    // Listen for messages from the server
    socket.onmessage = function(event) {
        var data = JSON.parse(event.data);
        console.log(data)
        var cursor = document.getElementById(data.userName);
        if (cursor) {
            // Remove the old cursor if it exists
            document.body.removeChild(cursor);
        }
        // Create a new cursor if it doesn't exist
        cursor = document.createElement('div');
        cursor.id = data.userName;
        cursor.style.position = 'absolute';
        cursor.style.width = '20px';
        cursor.style.height = '20px';
        cursor.style.background = data.cursorColor;
        document.body.appendChild(cursor);
        cursor.style.left = data.x + 'px';
        cursor.style.top = data.y + 'px';

        // Create a new span element for the username
        var usernameSpan = document.createElement('span');
        usernameSpan.style.position = 'absolute';
        usernameSpan.style.color = 'black';
        usernameSpan.style.fontSize = '20px';
        usernameSpan.textContent = data.userName;

        // Append the username span to the cursor div
        cursor.appendChild(usernameSpan);   
    };

    // Send cursor location data to the server
    document.onmousemove = function(event) {
        var data = JSON.stringify({
            userName: username,
            channel: channel,
            x: event.clientX,
            y: event.clientY
        });
        socket.send(data);
        console.log("sent", data)
    };
}