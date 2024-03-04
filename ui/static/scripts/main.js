const socket = new WebSocket('ws://localhost:8080/test');

socket.onopen = function(event) {
    console.log('Соединение установлено');
};

socket.onmessage = function(event) {
    console.log('Получено сообщение от сервера:', event.data);
};
function toggleChatbot() {

    let element = document.getElementById('chatbot-body')
    let element2 = document.getElementById('chatbot-footer')

    if (element.style.display === 'none') {
        element.style.display = 'block'
        element2.style.display = 'block'

    } else {
        element.style.display = 'none'
        element2.style.display = 'none'
    }
}


function sendMessage() {
    // Получаем элемент chatbot-body
    let chatbotBody = document.getElementById('chatbot-body');

    // Получаем текст, который пользователь ввел
    let userInput = document.getElementById('user-input').value;

    // Создаем новый элемент для сообщения пользователя
    let userMessageElement = document.createElement('div');

    // Добавляем класс для стилизации
    userMessageElement.classList.add('user-message');

    // Устанавливаем текст сообщения пользователя внутри элемента
    userMessageElement.innerText = "You: " + userInput;

    // Добавляем элемент с сообщением пользователя внутрь chatbot-body
    chatbotBody.appendChild(userMessageElement);
    sendMessageToServer(userInput);

}


function sendMessageToServer(message) {
    let url = '/message';
    let data = {
        message: message
    };
    let options = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    };

    fetch(url, options)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            // Обновляем DOM с полученными данными
            updateChat(data.message);
        })
        .catch(error => {
            console.error('There was a problem with your fetch operation:', error);
        });
}

function updateChat(message) {
    // Получаем элемент, куда будем добавлять сообщение от сервера
    let chatbotBody = document.getElementById('chatbot-body');

    // Создаем новый элемент для сообщения от сервера
    let serverMessageElement = document.createElement('div');
    serverMessageElement.classList.add('server-message');
    serverMessageElement.innerText = "MusicAiBot: " + message;

    // Добавляем сообщение от сервера в чат
    chatbotBody.appendChild(serverMessageElement);
}



