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
    // URL, на который будет отправлено сообщение
    let url = '/message';

    // Данные для отправки на сервер (в данном случае, просто текст сообщения)
    let data = {
        message: message
    };

    // Опции запроса
    let options = {
        method: 'POST', // метод запроса
        headers: {
            'Content-Type': 'application/json' // тип данных, отправляемых на сервер (JSON)
        },
        body: JSON.stringify(data) // преобразуем объект в JSON-строку
    };

    // Отправка запроса на сервер с использованием Fetch API
    fetch(url, options)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json(); // разбираем ответ сервера в формате JSON
        })
        .then(data => {
            // Обработка ответа сервера (если требуется)
            console.log('Server response:', data);
        })
        .catch(error => {
            // Обработка ошибок при отправке запроса
            console.error('There was a problem with your fetch operation:', error);
        });
}


