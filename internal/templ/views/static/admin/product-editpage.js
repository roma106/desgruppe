const productId = new URLSearchParams(window.location.search).get("id");
// const productType = new URLSearchParams(window.location.search).get("type");


const editButton = document.querySelector(".add-form-button");
editButton.addEventListener("click", () => {

    let data = {
        id: Number(productId),
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
        description: document.querySelector(".add-description").value,
        colors: GetCheckedColors(),
        position: Number(document.querySelector(".add-position").value),
        slug: document.querySelector(".add-slug").value,
    }
    if (data.onSale) {
        data.sale = Number(document.querySelector(".add-sale").value)
    }
    fetch("/admin/products/edit", {
        method: "PUT",
        body: JSON.stringify(data),
    }).then((res) => {
        if (res.ok) {
            alert("Изделие успешно изменено!");
        } else {
            if (res.status==501){
                alert("Этот url уже используется другим изделием")
            }else{
                alert("Произошла ошибка!");
            }
        }
    })
})


