document.addEventListener('DOMContentLoaded', function () {
  const tableBody = document.querySelector('#person-table tbody');
  const form = document.querySelector('#person-form');
  const addButton = document.querySelector('#add-update-btn');
  const resetButton = document.querySelector('#reset-btn');
  const lastManipulatedData = document.querySelector('#last-manipulated-data');

  let persons = JSON.parse(localStorage.getItem('persons')) || [];
  let lastManipulatedPerson = null;
  let updateIndex = null;

  renderTable();
  updateURL();
  showLastManipulatedData();

  form.addEventListener('submit', function (event) {
    event.preventDefault();
    const name = document.querySelector('#name').value;
    const gender = document.querySelector('#gender').value;
    const age = document.querySelector('#age').value;
    const city = document.querySelector('#city').value;

    if (updateIndex === null) {
      addPerson(name, gender, age, city);
      lastManipulatedPerson = { name, gender, age, city };
    } else {
      updatePerson(updateIndex, name, gender, age, city);
      lastManipulatedPerson = { name, gender, age, city };
    }

    form.reset();
    addButton.textContent = 'Add';
    updateIndex = null;
    updateURL();
    showLastManipulatedData();
  });

  resetButton.addEventListener('click', function () {
    form.reset();
    addButton.textContent = 'Add';
    updateIndex = null;
  });

  function renderTable() {
    tableBody.innerHTML = '';
    persons.forEach(function (person, index) {
      const row = createTableRow(person, index);
      tableBody.appendChild(row);
    });
    saveToLocalStorage();
  }

  function createTableRow(person, index) {
    const row = document.createElement('tr');
    row.innerHTML = `
      <td>${person.name}</td>
      <td>${person.gender}</td>
      <td>${person.age}</td>
      <td>${person.city}</td>
      <td><a href="#" class="update" data-index="${index}">Update</a> / <a href="#" class="remove" data-index="${index}">Remove</a></td>
    `;
    return row;
  }

  function addPerson(name, gender, age, city) {
    persons.push({ name, gender, age, city });
    renderTable();
  }

  function updatePerson(index, name, gender, age, city) {
    persons[index] = { name, gender, age, city };
    renderTable();
  }

  function saveToLocalStorage() {
    localStorage.setItem('persons', JSON.stringify(persons));
  }

  function updateURL() {
    const params = [];
    persons.forEach(function (person) {
      params.push(`name=${encodeURIComponent(person.name)}`);
      params.push(`gender=${encodeURIComponent(person.gender)}`);
      params.push(`age=${encodeURIComponent(person.age)}`);
      params.push(`city=${encodeURIComponent(person.city)}`);
    });
    const queryString = params.join('&');
    const newUrl = window.location.origin + window.location.pathname + '?' + queryString;
    history.pushState({}, '', newUrl);
  }

  function showLastManipulatedData() {
    if (lastManipulatedPerson) {
      lastManipulatedData.textContent = `Last Manipulated Data: Name: ${lastManipulatedPerson.name}, Gender: ${lastManipulatedPerson.gender}, Age: ${lastManipulatedPerson.age}, City: ${lastManipulatedPerson.city}`;
    } else {
      lastManipulatedData.textContent = 'No data manipulated yet.';
    }
  }

  tableBody.addEventListener('click', function (event) {
    event.preventDefault();
    const target = event.target;
    if (target.classList.contains('update')) {
      const index = target.getAttribute('data-index');
      const person = persons[index];
      updateIndex = index;
      fillFormWithData(person);
      addButton.textContent = 'Update';
    } else if (target.classList.contains('remove')) {
      const index = target.getAttribute('data-index');
      persons.splice(index, 1);
      renderTable();
      updateURL();
      showLastManipulatedData();
    }
  });

  function fillFormWithData(person) {
    document.querySelector('#name').value = person.name;
    document.querySelector('#gender').value = person.gender;
    document.querySelector('#age').value = person.age;
    document.querySelector('#city').value = person.city;
  }
});
