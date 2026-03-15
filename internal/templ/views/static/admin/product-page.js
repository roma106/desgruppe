
document.addEventListener("DOMContentLoaded", ()=>{
    ConfigureTable();
    ConfigurePages(document.querySelectorAll(".items-li"))
})


function ConfigureTable() {
    requestAnimationFrame(() => {
        ["position", "photo", "name", "producer", "designer", "size", "price", "discount"].forEach((tablecol) => {
            const cols = document.querySelectorAll(`.items-${tablecol}`);
            if (cols.length > 0) {
                let maxWidth = 0;
                cols.forEach(col => {
                    col.style.width = ''; // Сброс ширины
                    maxWidth = Math.max(maxWidth, col.offsetWidth);
                });
                cols.forEach(col => col.style.width = `${maxWidth}px`);
            }
        });
    });
}



// PAGES

let pagesContainer = document.querySelector(".pages-container");
let itemsTable = document.querySelector(".items-table");

function ConfigurePages(items) {
    pagesContainer.innerHTML = "";
    const itemsOnPage = 50;
    let pagesRequired = Math.floor(items.length / itemsOnPage) + 1;
    for (let i = 0; i < pagesRequired; i++) {
        let page = document.createElement("p");
        page.classList.add("page-p");
        page.innerHTML = i + 1;
        pagesContainer.appendChild(page);
        let pageItems = [];
        for (let j = 0; j < itemsOnPage; j++) {
            pageItems.push(items[i * itemsOnPage + j]);
        }
        if (i == 0) {
            itemsTable.innerHTML = "";
            pageItems.forEach((item) => {
                if (item != undefined) {
                    itemsTable.appendChild(item);
                }
            })
        }
        page.addEventListener("click", () => {
            itemsTable.innerHTML = "";
            pageItems.forEach((item) => {
                if (item != undefined) {
                    itemsTable.appendChild(item);
                }
            })
            ConfigureTable();
            for (const pageelem of pagesContainer.children) {
                pageelem.classList.remove("selected-page");
            }
            page.classList.add("selected-page");
        })
    }
    pagesContainer.children[0].classList.add("selected-page");
}



function OpenEditPage(ID, type){
    window.open(`/admin/products/editpage?id=${ID}&type=${type}`)
}

function OpenPhotosPage(id, name){
    window.open(`/admin/products/photos?id=${id}&name=${name}`);
}


// DELETING PRODUCT
let deletePopup = document.querySelector(".delete-color-popup");
let closeDeletePopupButton = document.querySelector(".close-delete-popup-button");
let deleteColorPopupButton = document.querySelector(".delete-color-popup-button");


closeDeletePopupButton.addEventListener("click", () => {
    deletePopup.classList.add("display-hidden");
})


function DeleteProduct(id, name){
    deletePopup.classList.remove("display-hidden");
    deleteColorPopupButton.addEventListener("click", () => {
        ConfirmedDelete(id);
    })
    document.querySelector(".delete-color-popup-text").innerHTML = `Вы действительно хотите удалить изделие ${name}?`;
}
function ConfirmedDelete(id){
    
    fetch("/admin/products/delete", {
        method: "DELETE", 
        body: JSON.stringify({ID: Number(id)})
    }).then((res) => {
        if (res.ok){
            DeleteAllPhotos(id);
            document.querySelector(`.items-li-${id}`).remove()
            deletePopup.classList.add("display-hidden");
        }else{
            alert("Произошла ошибка!");
            deletePopup.classList.add("display-hidden");
        }
    })
    deleteColorPopupButton.removeEventListener("click", DeleteProduct);
}

function DeleteAllPhotos(id){
    fetch("/admin/products/photos/deleteall?id=" + id, {
        method: "DELETE",
    }).then((res) => {
        if (!res.ok){
            alert("Произошла ошибка!");
        }
    })
}







// ADDING PRODUCT

let addProductButton = document.querySelector(".add-form-button")

addProductButton.addEventListener("click", AddProduct)

async function AddProduct() {
    if (document.querySelector(".add-photo").files.length == 0) {
        alert("Добавьте фото продукта")
        return
    }

    
    let data = {
        name: document.querySelector(".add-name").value,
        type: document.querySelector(".add-type").value,
        sectionId: Number(document.querySelector(".add-section").value),
        producerId: Number(document.querySelector(".add-producer").value),
        designerId: Number(document.querySelector(".add-designer").value),
        size: document.querySelector(".add-size").value,
        freeForm: document.querySelector(".add-free-form").value,
        available: document.querySelector(".add-available").checked,
        price: Number(document.querySelector(".add-price").value),
        onSale: document.querySelector(".add-sale-checkbox").checked,
        sale: 0,
        photo: Array.from(new Uint8Array(await document.querySelector(".add-photo").files[0].arrayBuffer())),
        // catalogPhotos: catalogPhotosBytes,
        description: document.querySelector(".add-description").value,
        colors: GetCheckedColors(),
        position:  Number(document.querySelector(".add-position").value),
    }
    if (data.onSale == true) {
        data.sale = Number(document.querySelector(".add-sale").value);
    }
    console.log(document.querySelector(".add-position").value)
    if (document.querySelector(".add-position").value === ""){
        data.position = 0;
    }
    // return

    fetch("/admin/products/add",
        {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        }
    ).then(res => {
        if (res.ok) {
            alert("Продукт успешно добавлен!");
            CleanFields();
            localStorage.clear();
            return res.text();
        } else {
            alert("Произошла ошибка!");
        }
    }).then(res => AddPhotosToGallery(res)).finally(()=>window.location.reload())

}

async function AddPhotosToGallery(productID){
    // Обработка фотографий
    let catalogPhotos = document.querySelector(".add-catalog-photo").files;


    for (let i = 0; i < catalogPhotos.length; i++) {
        let catalogPhotoBytes = new Uint8Array(await document.querySelector(".add-catalog-photo").files[i].arrayBuffer())
        await fetch("/admin/products/photos/add", {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                productID: productID,
                photo: Array.from(catalogPhotoBytes),
                Position: i,
            })
        })
    }
}

function CleanFields() {
    document.querySelector(".add-name").value = "";
    document.querySelector(".add-size").value = "";
    document.querySelector(".add-free-form").value = "";
    document.querySelector(".add-available").checked = false;
    document.querySelector(".add-price").value = "";
    document.querySelector(".add-sale-checkbox").checked = false;
    document.querySelector(".add-sale").value = "";
    document.querySelector(".add-photo").value = "";
    document.querySelector(".add-description").value = "";
    document.querySelector(".add-position").value = "";
    document.querySelectorAll(".add-form-color-checkbox").forEach((checkbox) => {
        checkbox.checked = false;
    });
    document.querySelector(".add-form-sale-block").classList.add("hidden");
    document.querySelector(".add-form-selected-colors").innerHTML="";
}
