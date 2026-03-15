
// GETTING  DESIGNER
window.onload = () => {
    if (window.location.href.includes("designer")) {
        GetDesigner("designer");
    } else if (window.location.href.includes("producer")) {
        GetDesigner("producer");
    }
}
function GetDesigner(type) {
    fetch(`/${type}/get/` + window.location.href.split("/").at(-1))
        .then((res) => res.json())
        .then((designer) => {
            console.log(designer)
            let photourl = URL.createObjectURL(new Blob([new Uint8Array(designer.picture)]));
            document.querySelector(`.${type}-info-img`).src = photourl;
            GetDesignerProducts(designer.id, type);
        })
}

function GetDesignerProducts(id, type) {
    if (localStorage.getItem("exchangeRate")) {
        fetch(`/products/listbyfilter?${type}s=` + id)
            .then((res) => res.json())
            .then((data) => {
                CreateProductNodes(data, localStorage.getItem("exchangeRate"));
            })
    } else {
        fetch("/ex-rate").then((res) => res.json()).then((ex) => {
            localStorage.setItem("exchangeRate", ex.exchangeRate);
        }).then(()=>{
            fetch(`/products/listbyfilter?${type}s=` + id)
                .then((res) => res.json())
                .then((data) => {
                    CreateProductNodes(data, localStorage.getItem("exchangeRate"));
                })
        })
    }
}

function CreateProductNodes(data, exRate) {
    let itemsContainer = document.querySelector(".designer-items-block");
    itemsContainer.innerHTML = "";
    ex=Number(exRate)
    data.forEach((product) => {
        let item = document.createElement("div");
        item.classList.add("shop-item");
        let photourl = URL.createObjectURL(new Blob([new Uint8Array(product.photo)]));
        let pricelabel = "₽";
        item.innerHTML =
            `<a href="/product/${product.slug}" class="shop-item-link">
            <img src="${photourl}" alt="" class="shop-item-img">
            <div class="shop-item-info-container">
                <p class="shop-item-name">${product.name}</p>
                <div class="shop-item-price-container">
                    <p class="shop-item-price">${formatPrice((product.price*ex).toFixed(2))}</p>
                    <span class="shop-item-price-label">${pricelabel}</span>
                </div>
            </div>
            </a>`;

        if (product.onSale) {
            let productPriceContainer = item.querySelector(".shop-item-price-container");
            productPriceContainer.children[0].classList.add("shop-item-line-through");
            productPriceContainer.children[1].classList.add("shop-item-line-through");
            productPriceContainer.innerHTML += `
                    <p class="shop-item-sale-price">${formatPrice((product.price*ex * ((100 - product.sale) / 100)).toFixed(2))}</p>
                    <span class="shop-item-sale-price-label">${pricelabel}</span>
                `;
        }

        itemsContainer.appendChild(item);
    })
}



// READ MORE
let readMoreButton = document.querySelector(".designer-read-more");
readMoreButton.addEventListener("click", () => {
    document.querySelector(".designer-info-text-container").classList.add("designer-info-container-visible");
    document.querySelector(".designer-info-container").classList.add("designer-info-container-visible");
    readMoreButton.style.display = "none";
    document.querySelector(".designer-info-gradient").style.display = "none";
})



function formatPrice(numberString) {
    const numericString = numberString.replace(/[^\d.]/g, '');
    let [integerPart, decimalPart] = numericString.split('.');
    integerPart = integerPart.replace(/\B(?=(\d{3})+(?!\d))/g, ' ');
    return decimalPart ? `${integerPart}.${decimalPart}` : integerPart;
}