const input = document.querySelector('#textarea')
const messages = document.querySelector('#messages')
const username = document.querySelector('#username')
const send = document.querySelector('#send')

const url = "ws://" + window.location.host + "/ws";
const ws = new WebSocket(url);

ws.onmessage = (msg) => {
    insertMessage(JSON.parse(msg.data));
};

send.onclick = () => {
    if(username.value == "") {
        alert("Enter a username.");
        return;
    }
    else if(input.value == "") {
        alert("Message body can't be empty.");
        return;
    }

    const message = {
        username: username.value,
        message: input.value,
    }

    ws.send(JSON.stringify(message));
    input.value = "";
};

insertMessage = (messageObj) => {
    const message = document.createElement('div');

    message.setAttribute('class', 'chat-message')
    console.log("name: " + messageObj.username + " content: " + messageObj.message)
    message.textContent = `${messageObj.username}: ${messageObj.message}`

    messages.appendChild(message)
    messages.insert(message, messages.firstChild)
}