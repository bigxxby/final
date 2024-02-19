// Функция для открытия модального окна
function openModal() {
    var modal = document.getElementById("myModal");
    modal.style.display = "block";
}

// Функция для закрытия модального окна
function closeModal() {
    var modal = document.getElementById("myModal");
    modal.style.display = "none";
}

// Закрыть модальное окно при клике вне его области
window.onclick = function(event) {
    var modal = document.getElementById("myModal");
    if (event.target == modal) {
        modal.style.display = "none";
    }
}
