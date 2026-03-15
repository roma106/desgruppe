const productID = new URLSearchParams(window.location.search).get("id");
const productName = new URLSearchParams(window.location.search).get("name");

document.querySelector(".title").innerHTML = `Фотогалерея изделия ${productName}`;
document.querySelector(".delete-color-popup-text").innerHTML = `Вы действительно хотите удалить фото для изделия ${productName}?`;


// MAIN PHOTO

GetMainPhoto();

function GetMainPhoto() {
    fetch("/admin/products/photos/getmain?id=" + productID)
        .then((res) => res.json()).then((data) => {
            const photourl = URL.createObjectURL(new Blob([new Uint8Array(data)]));
            document.querySelector(".photos-main-img").src = photourl;
        })
}



let addMainPhotoButton = document.querySelector(".add-main-photo-button")
addMainPhotoButton.addEventListener("click", AddMainPhoto)

async function AddMainPhoto() {
    let photoInputFile = document.querySelector(".add-photo").files[0]
    if (!photoInputFile) {
        alert("Добавьте файл")
        return
    }

    let data = {
        id: Number(productID),
        photo: Array.from(new Uint8Array(await photoInputFile.arrayBuffer()))
    }
    fetch("/admin/products/photos/editmain", {
        method: "POST",
        body: JSON.stringify(data),
    }).then(res => {
        if (res.ok) {
            alert("Фото успешно обновлено!");
            GetMainPhoto();
        } else {
            alert("Произошла ошибка!");
        }
    })
}




// GALLERY

ListPhotos();


function ListPhotos() {
    fetch("/admin/products/photos/list?id=" + productID)
        .then((res) => res.json()).then((data) => {
            CreateCatalogNodes(data);
        })
}

function CreateCatalogNodes(data) {
    let catalogPhotosContainer = document.querySelector(".photos-items-container")
    catalogPhotosContainer.innerHTML = ""

    data.forEach(galleryPhoto => {
        let newPhoto = document.createElement("div")
        newPhoto.classList.add("photo-item")
        const photourl = URL.createObjectURL(new Blob([new Uint8Array(galleryPhoto.photo)]));
        newPhoto.innerHTML =
            `<img src="${photourl}" alt="" class="photo-item-img" id="${galleryPhoto.id}">
            <div class="photo-item-options">
                <img src="../../imgs/icons/arrow-icon.png" alt="" class="photo-item-position-prev">
                <img src="../../imgs/icons/arrow-icon.png" alt="" class="photo-item-position-next">
                <img src="../../imgs/icons/admin-delete-icon.png" alt="" class="photo-item-delete">
            </div>`;
        newPhoto.querySelector(".photo-item-position-prev").addEventListener("click", () => PhotoPositionPrev(newPhoto))
        newPhoto.querySelector(".photo-item-position-next").addEventListener("click", () => PhotoPositionNext(newPhoto))
        newPhoto.querySelector(".photo-item-delete").addEventListener("click", () => DeletePhoto(galleryPhoto.id))
        catalogPhotosContainer.appendChild(newPhoto)
    })
}

function PhotoPositionNext(photo) {
    let photosContainer = document.querySelector(".photos-items-container")
    if (photo.nextElementSibling.nextElementSibling) {
        photosContainer.insertBefore(photo, photo.nextElementSibling.nextElementSibling)
    } else {
        photosContainer.appendChild(photo)
    }
}
function PhotoPositionPrev(photo) {
    let photosContainer = document.querySelector(".photos-items-container")
    if (photo.previousElementSibling) {
        photosContainer.insertBefore(photo, photo.previousElementSibling)
    }
}



let deletePopup = document.querySelector(".delete-color-popup");
let closeDeletePopupButton = document.querySelector(".close-delete-popup-button");
let deleteColorPopupButton = document.querySelector(".delete-color-popup-button");


closeDeletePopupButton.addEventListener("click", () => {
    deletePopup.classList.add("display-hidden");
})
function DeletePhoto(id){
    deletePopup.classList.remove("display-hidden");
    deleteColorPopupButton.addEventListener("click", () => {
        ConfirmedDelete(id);
    })
}
function ConfirmedDelete(id){
    
    fetch("/admin/products/photos/delete?id=" + id, {
        method: "DELETE",
    }).then((res) => {
        if (res.ok){
            ListPhotos();
            deletePopup.classList.add("display-hidden");
        }else{
            alert("Произошла ошибка!");
            deletePopup.classList.add("display-hidden");
        }
    })
    deleteColorPopupButton.removeEventListener("click", DeletePhoto);
}



let editGalleryButton = document.querySelector(".edit-photo-button")

editGalleryButton.addEventListener("click", EditGallery)

function EditGallery() {

    let catalogPhotos = document.querySelectorAll(".photo-item-img");

    for (let i = 0; i < catalogPhotos.length; i++) {
        let data = {
            id: Number(catalogPhotos[i].id),
            productID: productID,
            position: i,
        };

        fetch("/admin/products/photos/edit", {
            method: "POST",
            body: JSON.stringify(data)
        }).then(res => {
            if (!res.ok) {
                alert("Произошла ошибка!");
            }
        })
    }
    alert("Изменения успешно зафиксированы");
    ListPhotos();
}


let addPhotoButton = document.querySelector(".add-photo-button");

addPhotoButton.addEventListener("click", AddPhoto)
async function AddPhoto(){
    let photoInputFile = document.querySelector(".add-catalog-photo").files[0];
    if (!photoInputFile) {
        alert("Добавьте файл")
        return
    }
    let data = {
        productID: productID,
        photo: Array.from(new Uint8Array(await photoInputFile.arrayBuffer())),
        position: document.querySelectorAll(".photo-item").length,
    }
    fetch("/admin/products/photos/add", {
        method: "POST",
        body: JSON.stringify(data),
    }).then(res => {
        if (res.ok) {
            alert("Фото успешно добавлено!");
            ListPhotos();
        } else {
            alert("Произошла ошибка!");
        }
    })
}