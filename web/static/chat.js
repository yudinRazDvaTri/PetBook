new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        to_id: 0, // Our username
    },

    created: function() {
        var self = this;
        to_id = document.getElementById("UserID").textContent;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws/' + to_id);
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data)
            self.chatContent += 
            //'<button style= "border-style:none; background:white;" class="material-icons">highlight_off</button>'
                    '<div class="chip">'
                    + msg.username +' '
                    +'<i>' + msg.created_at + '</i>'
                    + '</div>'
                    + msg.message + '<br/>';

            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight - element.clientHeight; // Auto scroll to the bottom
        });
    },

    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        username: this.username,
                        message: $('<p>').html(this.newMsg).text() // Strip out html
                        //created_at: now
                    }
                ));
                this.newMsg = ''; // Reset newMsg
            }
        },
        searchByDate() {

            cancelSearchBtnClass = document.getElementById("cancelSearchBtn").className;
            var hiddenWordIndex = cancelSearchBtnClass.lastIndexOf(" ");
            document.getElementById("cancelSearchBtn").className = cancelSearchBtnClass.substring(0, hiddenWordIndex);
            document.getElementById("sendMessageBtn").className += " hide"
            document.getElementById("messageInputbox").className += " hide"

            var calendar = M.Datepicker.getInstance(document.getElementById('calendar'));
            var date = calendar.toString();
            if(date != ""){
                this.chatContent='';
                var xhr = new XMLHttpRequest();
                xhr.open('GET', 'http://' + window.location.host + '/chats/' + to_id + '/search/' + date, false);
                xhr.send();
                if (xhr.status != 200) {
                    alert( xhr.status + ': ' + xhr.statusText );
                    } 
                    else {
                        var messages = JSON.parse(xhr.responseText)
                        for(var i = 0; i < messages.length; i++){
                            var msg = messages[i]
                            this.chatContent += '<div class="chip">'
                                    + msg.username +' '
                                    +'<i>' + msg.created_at + '</i>'
                                    + '</div>'
                                    + msg.message + '<br/>';
                            var element = document.getElementById('chat-messages');
                            element.scrollTop = element.scrollHeight - element.clientHeight;
                            }
                 }
          } 
          else {
            document.location.reload(true);
          }
    },
    cancelSearch(){
        document.location.reload(true);

        cancelSearchBtnClass = document.getElementById("sendMessageBtn").className;
        var hiddenWordIndex = cancelSearchBtnClass.lastIndexOf(" ");
        document.getElementById("sendMessageBtn").className = cancelSearchBtnClass.substring(0, hiddenWordIndex);

        inputMessageClass = document.getElementById("messageInputbox").className;
        var hiddenWordIndex = cancelSearchBtnClass.lastIndexOf(" ");
        document.getElementById("messageInputbox").className = cancelSearchBtnClass.substring(0, hiddenWordIndex);

        document.getElementById("cancelSearchBtn").className += " hide";
    },
    }
});
