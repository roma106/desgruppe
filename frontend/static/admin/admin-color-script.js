ListColors();





let colorRadio = document.querySelectorAll(".add-color-radio");
let pictureRadio = document.querySelectorAll(".add-picture-radio");


colorRadio[0].addEventListener("click", () => {
    pictureRadio[0].checked = false;
    pictureRadio[2].disabled = true;
    colorRadio[2].disabled = false;
    pictureRadio.forEach((radio) => radio.style.color = "gray");
    colorRadio.forEach((radio) => radio.style.color = "black");
})

pictureRadio[0].addEventListener("click", () => {
    colorRadio[0].checked = false;
    colorRadio[2].disabled = true;
    pictureRadio[2].disabled = false;
    pictureRadio.forEach((radio) => radio.style.color = "black");
    colorRadio.forEach((radio) => radio.style.color = "gray");
})






//  ДОБАВЛЕНИЕ  ЦВЕТА
let addColorButton = document.querySelector(".add-form-button");

addColorButton.addEventListener("click", () => {
    if (!Validate()) { return }
    
    let colorPosition = document.querySelector(".add-position").value;
    if (colorPosition === "") {
        fetch("/admin/colors/list")
        .then((res) => res.json())
        .then((list) => AddColor(list.length+1))
    }else{
        AddColor(colorPosition)
    }
})

async function AddColor(colorPosition) {
    let colorColorRadio = document.querySelector(".add-color-radio").checked;
    let colorPictureRadio = document.querySelector(".add-picture-radio").checked;

    let data = {
        name: document.querySelector(".add-name").value,
        code: null,
        picture: null,
        position: Number(colorPosition),
        codeRadio: colorColorRadio,
        pictureRadio: colorPictureRadio
    }
    if (colorColorRadio) {
        data.code = document.querySelector(".add-color").value;
    } else if (colorPictureRadio) {
        data.picture = Array.from(new Uint8Array(await document.querySelector(".add-picture").files[0].arrayBuffer()))
    }
    fetch("/admin/colors/add", {
        method: "POST",
        body: JSON.stringify(data),
    }).then((res) => {
        if (res.ok) {
            CleanFields();
            ListColors();
            alert("Цвет успешно добавлен!");
        } else {
            alert("Произошла ошибка!");
        }
    })
}


function ListColors() {
    fetch("/admin/colors/list", {
        method: "GET",
    }).then((res) => {
        if (res.ok) {
            res = res.json();
            res.then((data) => {
                let table = CreateColorNodes(data);
                ConfigurePages(table);
                ConfigureTable();
            })
        } else {
            alert("Неудалось получить список цветов");
        }
    })
}

function CreateColorNodes(data) {
    let table = [];
    data.forEach((color) => {
        if (color.code != "") {
            let newcolor = document.createElement("div");
            newcolor.className = `items-li items-li-${color.id}`;
            newcolor.innerHTML = (`
                <p class="items-li-elem items-position">${color.position}</p>
                <span class="items-li-elem items-color" style="background-color: #${color.code};"></span>
                <p class="items-li-elem items-name">${color.name}</p>
                <img src="../imgs/icons/admin-edit-icon.png" alt="" class="items-li-edit" onclick="EditColor(${color.id})">
                <img src="../imgs/icons/admin-delete-icon.png" alt="" class="items-li-delete" onclick="DeleteColor(${color.id}, '${color.name}')">`)
            table.push(newcolor);
        } else if (color.code == "") {
            let picturebytes = new Int8Array(color.picture);
            // console.log(picturebytes);
            const blob = new Blob([picturebytes], { type: "image/png" });
            const url = URL.createObjectURL(blob);
            let newcolor = document.createElement("div");
            newcolor.className = `items-li items-li-${color.id}`;
            newcolor.innerHTML = (`
                    <p class="items-li-elem items-position">${color.position}</p>
                    <span class="items-li-elem items-color" style="background-image: url(${url});"></span>
                    <p class="items-li-elem items-name">${color.name}</p>
                    <img src="../imgs/icons/admin-edit-icon.png" alt="" class="items-li-edit" onclick="EditColor(${color.id})">
                    <img src="../imgs/icons/admin-delete-icon.png" alt="" class="items-li-delete" onclick="DeleteColor(${color.id}, '${color.name}')">`)
            table.push(newcolor);
        }
    })
    return table;
}


function EditColor(id) {
    window.open(`/admin/colors/editpage?id=${id}`);
}


// УДАЛЕНИЕ  ЦВЕТА
let deletePopup = document.querySelector(".delete-color-popup");
let closeDeletePopupButton = document.querySelector(".close-delete-popup-button");
let deleteColorPopupButton = document.querySelector(".delete-color-popup-button");


closeDeletePopupButton.addEventListener("click", () => {
    deletePopup.classList.add("display-hidden");
})


function DeleteColor(id, name) {
    deletePopup.classList.remove("display-hidden");
    deleteColorPopupButton.addEventListener("click", () => {
        ConfirmedDelete(id);
    })
    document.querySelector(".delete-color-popup-text").innerHTML = `Вы действительно хотите удалить цвет ${name}?`;
}
function ConfirmedDelete(id) {

    fetch("/admin/colors/delete", {
        method: "DELETE",
        body: JSON.stringify({ id: id })
    }).then((res) => {
        if (res.ok) {
            ListColors();
            deletePopup.classList.add("display-hidden");
        } else {
            alert("Произошла ошибка!");
            deletePopup.classList.add("display-hidden");
        }
    })
    deleteColorPopupButton.removeEventListener("click", DeleteColor);
}


function Validate() {
    let colorName = document.querySelector(".add-name").value;
    let colorColor = document.querySelector(".add-color").value;
    let colorPicture = document.querySelector(".add-picture").value;

    if (colorName.length == 0 || colorColor.length == 0 && colorPicture.length == 0) {
        alert("Заполните все поля!");
        return false
    }
    return true
}

function CleanFields() {
    document.querySelector(".add-name").value = "";
    document.querySelector(".add-color").value = "";
    document.querySelector(".add-picture").value = "";
    document.querySelector(".add-position").value = "";

}