new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        username: null, // Our username
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            self.chatContent += '<div class="chip">'
                    + msg.username
                + '</div>'
                + emojione.toImage(msg.message) + '<br/>' // Parse emojis
                + '<i>' + msg.created_at+ '</i>' +'<br/>';
            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
        });
    },

    methods: {
        send: function () {
            var now = new Date();
            now = now.toString();
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        username: this.username,
                        message: $('<p>').html(this.newMsg).text(), // Strip out html
                        created_at: now
                    }
                ));
                this.newMsg = ''; // Reset newMsg
            }
        },
    }
});
