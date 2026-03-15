
ListSections();







//  ДОБАВЛЕНИЕ  РАЗДЕЛА
let addSectionButton = document.querySelector(".add-form-button");

addSectionButton.addEventListener("click", () => {
    if (!Validate()) { return }
    let sectionName = document.querySelector(".add-name").value;
    let sectionType = document.querySelector(".add-type").value;


    let data = {
        name: sectionName,
        type: sectionType,
    }
    fetch(window.location.href + "/add", {
        method: "POST",
        body: JSON.stringify(data),
    }).then((res) => {
        if (res.ok) {
            CleanFields();
            ListSections();
            alert("Элемент успешно добавлен!");
        } else {
            alert("Произошла ошибка!");
        }
    })
})


function ListSections() {
    fetch(window.location.href + "/list", {
        method: "GET",
    }).then((res) => {
        if (res.ok) {
            res = res.json();
            res.then((data) => {
                CreateSectionNodes(data);
            })
        } else {
            alert("Неудалось получить список элементов");
        }
    })
}

function CreateSectionNodes(data) {
    let furnitureTable = document.querySelectorAll(".items-table")[0];
    furnitureTable.innerHTML="";
    let lightTable = document.querySelectorAll(".items-table")[1];
    lightTable.innerHTML="";
    let interiorTable = document.querySelectorAll(".items-table")[2];
    interiorTable.innerHTML="";
    data.forEach((sec) => {
        let newsection = document.createElement("div");
        newsection.className = `items-li items-li-${sec.id}`;
        newsection.innerHTML = (`
        <p class="items-li-elem items-name">${sec.name}</p>
        <img src="../imgs/icons/admin-delete-icon.png" alt="" class="items-li-delete" onclick="DeleteSection(${sec.id}, '${sec.name}')">`)
        
        if (sec.type === "furniture"){
            furnitureTable.appendChild(newsection)
        }else if (sec.type === "light"){
            lightTable.appendChild(newsection)
        }else if (sec.type === "interior"){
            interiorTable.appendChild(newsection)
        }
    })
}



// УДАЛЕНИЕ  ДИЗАЙНЕРА
let deletePopup = document.querySelector(".delete-color-popup");
let closeDeletePopupButton = document.querySelector(".close-delete-popup-button");
let deletePopupButton = document.querySelector(".delete-color-popup-button");


closeDeletePopupButton.addEventListener("click", () => {
    deletePopup.classList.add("display-hidden");
    deletePopupButton.onclick = () => {}
})


function DeleteSection(id, name) {
    deletePopup.classList.remove("display-hidden");
    deletePopupButton.onclick = () => {ConfirmedDelete(id)}
    document.querySelector(".delete-color-popup-text").innerHTML = `Вы действительно хотите удалить раздел ${name}?`;
    
}
function ConfirmedDelete(id) {

    fetch(window.location.href + "/delete", {
        method: "DELETE",
        body: JSON.stringify({ id: id })
    }).then((res) => {
        if (res.ok) {
            ListSections();
            deletePopup.classList.add("display-hidden");
        } else {
            alert("Произошла ошибка!");
            deletePopup.classList.add("display-hidden");
        }
    })
    deletePopupButton.onclick = () => {}
}


function Validate() {
    let desName = document.querySelector(".add-name").value;
    let desType = document.querySelector(".add-type").value;

    if (desName.length == 0 || desType.length == 0) {
        alert("Заполните все поля!");
        return false
    }
    return true
}

function CleanFields() {
    document.querySelector(".add-name").value = "";

}








