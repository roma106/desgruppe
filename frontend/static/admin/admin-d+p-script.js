
ListDesigners();







//  ДОБАВЛЕНИЕ  ДИЗАЙНЕРА
let addDesignerButton = document.querySelector(".add-form-button");

addDesignerButton.addEventListener("click", () => {
    if (!Validate()) { return }
    let desName = document.querySelector(".add-name").value;
    let desDescription = document.querySelector(".add-description").value;
    let desPicture = document.querySelector(".add-picture").files[0];
    let desPosition = document.querySelector(".add-position").value;
    if (desPosition === ""){
        desPosition = document.querySelector(".items-table").children.length+1
    }


    const reader = new FileReader();
    reader.readAsArrayBuffer(desPicture);
    reader.onload = () => {
        let data = {
            name: desName,
            description: desDescription,
            picture: Array.from(new Uint8Array(reader.result)),
            position: Number(desPosition),
        }
        fetch(window.location.href + "/add", {
            method: "POST",
            body: JSON.stringify(data),
        }).then((res) => {
            if (res.ok) {
                CleanFields();
                ListDesigners();
                alert("Элемент успешно добавлен!");
            } else {
                alert("Произошла ошибка!");
            }
        })
    }
})


function ListDesigners() {
    fetch(window.location.href + "/list", {
        method: "GET",
    }).then((res) => {
        if (res.ok) {
            res = res.json();
            res.then((data) => {
                let table = CreateDesignerNodes(data);
                if (window.location.href.includes("designers")) {
                    ConfigurePages(table);
                    ConfigureTable();
                } else if (window.location.href.includes("producers")) {
                    ConfigureProducersTable(table);
                }
            })
        } else {
            alert("Неудалось получить список элементов");
        }
    })
}

function CreateDesignerNodes(data) {
    let table = [];
    data.forEach((des) => {
        let newdes = document.createElement("div");
        newdes.className = `items-li items-li-${des.id}`;
        newdes.innerHTML = (`
                    <p class="items-li-elem items-position">${des.position}</p>
                    <p class="items-li-elem items-name">${des.name}</p>
                    <img src="../imgs/icons/admin-edit-icon.png" alt="" class="items-li-edit" onclick="EditDesigner(${des.id})">
                    <img src="../imgs/icons/admin-delete-icon.png" alt="" class="items-li-delete" onclick="DeleteDesigner(${des.id}, '${des.name}')">`)
        table.push(newdes);
    })
    return table;
}



function EditDesigner(id) {
    window.open(window.location.href + `/editpage?id=${id}`);
}


// УДАЛЕНИЕ  ДИЗАЙНЕРА
let deletePopup = document.querySelector(".delete-color-popup");
let closeDeletePopupButton = document.querySelector(".close-delete-popup-button");
let deletePopupButton = document.querySelector(".delete-color-popup-button");


closeDeletePopupButton.addEventListener("click", () => {
    deletePopup.classList.add("display-hidden");
})


function DeleteDesigner(id, name) {
    deletePopup.classList.remove("display-hidden");
    deletePopupButton.addEventListener("click", () => {
        ConfirmedDelete(id);
    })
    if (window.location.href.includes("designers")) {
        document.querySelector(".delete-color-popup-text").innerHTML = `Вы действительно хотите удалить дизайнера ${name}?`;
    } else if (window.location.href.includes("producers")) {
        document.querySelector(".delete-color-popup-text").innerHTML = `Вы действительно хотите удалить производителя ${name}?`;
    }
}
function ConfirmedDelete(id) {

    fetch(window.location.href + "/delete", {
        method: "DELETE",
        body: JSON.stringify({ id: id })
    }).then((res) => {
        if (res.ok) {
            ListDesigners();
            deletePopup.classList.add("display-hidden");
        } else {
            alert("Произошла ошибка!");
            deletePopup.classList.add("display-hidden");
        }
    })
    deletePopupButton.removeEventListener("click", DeleteDesigner);
}


function Validate() {
    let desName = document.querySelector(".add-name").value;
    let desDescription = document.querySelector(".add-description").value;
    let desPicture = document.querySelector(".add-picture").value;

    if (desName.length == 0 || desDescription.length == 0 || desPicture.length == 0) {
        alert("Заполните все поля!");
        return false
    }
    return true
}

function CleanFields() {
    document.querySelector(".add-name").value = "";
    document.querySelector(".add-description").value = "";
    document.querySelector(".add-picture").value = "";
    document.querySelector(".add-position").value = "";

}









function ConfigureProducersTable(table) {

    let tableContainer = document.querySelector(".items-table");

    tableContainer.innerHTML = "";

    table.forEach((item) => {
        tableContainer.appendChild(item);
    });


    ["position", "name"].forEach((tablecol) => {
        let col = document.querySelectorAll(".items-" + tablecol);
        for (let i = 0; i < col.length; i++) {

            if (i < col.length - 1) {
                if (col[i].offsetWidth >= col[i + 1].offsetWidth) {
                    col[i + 1].style.width = col[i].offsetWidth - 16 + "px";
                } else {
                    col[i].style.width = col[i + 1].offsetWidth - 16 + "px";
                }
            }
        }
    })
}