document.addEventListener("DOMContentLoaded", ()=>{
    ConfigureTable();
    ConfigureOrdersUI();
})

function ConfigureOrdersUI(){
    let editButtons = document.querySelectorAll(".items-li-order")
    let deleteButtons = document.querySelectorAll(".items-li-delete")
    editButtons.forEach(but =>{
        but.addEventListener("click", () => {
            EditOrder(but.parentElement.id, !but.parentElement.classList.contains("order-not-seen"))
        })
    })
    deleteButtons.forEach(but =>{
        but.addEventListener("click", () => {
            DeleteOrder(but.parentElement.id, but.parentElement.querySelector(".items-date").innerHTML)
        })
    })
}

function EditOrder(id, seen) {
    if (!seen){
        fetch("/admin/orders/seen", {
            method: "PUT",
            body: JSON.stringify({id: Number(id), seen: true})
        })
    }
    window.open('/admin/orders/editpage?id='+id)
}


// УДАЛЕНИЕ  Заказа
let deletePopup = document.querySelector(".delete-color-popup");
let closeDeletePopupButton = document.querySelector(".close-delete-popup-button");
let deletePopupButton = document.querySelector(".delete-color-popup-button");


closeDeletePopupButton.addEventListener("click", () => {
    deletePopup.classList.add("display-hidden");
    deletePopupButton.onclick = () => {}
})


function DeleteOrder(id, date) {
    deletePopup.classList.remove("display-hidden");
    deletePopupButton.onclick = () => {ConfirmedDelete(id)}
    document.querySelector(".delete-color-popup-text").innerHTML = `Вы действительно хотите удалить заказ от ${date}?`;
    
}
function ConfirmedDelete(id) {

    fetch(window.location.href + "/delete", {
        method: "DELETE",
        body: JSON.stringify({ id: Number(id) })
    }).then((res) => {
        if (res.ok) {
            window.location.reload();
        } else {
            alert("Произошла ошибка!");
            deletePopup.classList.add("display-hidden");
        }
    })
    deletePopupButton.onclick = () => {}
}

function ConfigureTable() {
    requestAnimationFrame(() => {
      ["date", "name","email", "phone"].forEach((tablecol) => {
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
  




