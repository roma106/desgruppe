let filterButtonCategories = document.querySelector(".filter-button-categories");
let filterCategoriesList = document.querySelector(".shop-filter-sections");

filterButtonCategories.addEventListener("click", () => {
    filterCategoriesList.classList.toggle("filter-categories-hidden");
    document.querySelector(".categories-filter-img").classList.toggle("shop-filter-img-opened");
})

let filterButtonDesigners = document.querySelector(".filter-button-designers");
let filterDesignersList = document.querySelector(".shop-filter-designers");

filterButtonDesigners.addEventListener("click", () => {
    filterDesignersList.classList.toggle("filter-designers-hidden");
    document.querySelector(".designers-filter-img").classList.toggle("shop-filter-img-opened");
})



let filterButtonFabrics = document.querySelector(".filter-button-fabrics");
let filterFabricsList = document.querySelector(".shop-filter-producers");

filterButtonFabrics.addEventListener("click", () => {
    filterFabricsList.classList.toggle("filter-fabrics-hidden");
    document.querySelector(".fabrics-filter-img").classList.toggle("shop-filter-img-opened");
})


let filterButtonAvailable = document.querySelector(".filter-button-available");
let filterAvailableList = document.querySelector(".shop-filter-available");

filterButtonAvailable.addEventListener("click", () => {
    filterAvailableList.classList.toggle("filter-available-hidden");
    document.querySelector(".available-filter-img").classList.toggle("shop-filter-img-opened");
})

//  PRICE FILTER

let filterButtonPrice = document.querySelector(".filter-button-price");
let filterPriceList = document.querySelector(".shop-filter-price");

filterButtonPrice.addEventListener("click", () => {
    filterPriceList.classList.toggle("filter-price-hidden");
    document.querySelector(".price-filter-img").classList.toggle("shop-filter-img-opened");
})

let priceRange1 = document.querySelector(".price-range-1");
let priceRange2 = document.querySelector(".price-range-2");
let priceLabel1 = document.querySelector(".price-label-1");
let priceLabel2 = document.querySelector(".price-label-2");

const debouncedListProducts = debounce(() => {
    ListProductsByPrice();
}, 200);

priceRange1.addEventListener("input", () => {
    if (Number(priceRange1.value) > Number(priceRange2.value)) {
        priceRange1.value = priceRange2.value - 1;
    }
    priceLabel1.innerHTML = priceRange1.value + " " + String(priceLabel1.innerHTML).at(-1);
    debouncedListProducts()
})
priceRange2.addEventListener("input", () => {
    if (Number(priceRange2.value) < Number(priceRange1.value)) {
        priceRange2.value = Number(priceRange1.value) + 1;
    }
    priceLabel2.innerHTML = priceRange2.value + " " + String(priceLabel1.innerHTML).at(-1);
    debouncedListProducts()
})

function ListProductsByPrice() {
    let shopItems = document.querySelectorAll(".shop-item");
    let minprice = Number(priceLabel1.innerHTML.split(" ")[0]);
    let maxprice = Number(priceLabel2.innerHTML.split(" ")[0]);
    shopItems.forEach(item => {
        item.style.display = "flex";
        let itemprice = Number(item.querySelector(".shop-item-price").innerHTML.replace(" ", ""))
        if (itemprice > maxprice || itemprice < minprice) {
            item.style.display = "none";
        }
    })
}



// SORT

let sortSelect = document.querySelector(".shop-sort-select")

sortSelect.addEventListener("change", () => {
    ListProductsWithFilter(localStorage.getItem("type"))
})




// Переход с главной страницы от кнопки "Смотреть все" начличие
let urlParams = new URLSearchParams(document.location.search)
let availableParam = Boolean(urlParams.get("available"))


function ShowOnlyAvailable(){
    document.querySelector(".shop-filter-available-checkbox").ariaChecked = true;
    document.querySelector(".shop-filter-available-checkbox-bg").classList.remove("hidden");
    ListProductsWithFilter(localStorage.getItem("type"))
}




// GETTING PROUCTS

function ListProductsWithFilter(type) {
    let sectionIds = [];
    let producerIds = [];
    let designerIds = [];
    let available = document.querySelector(".shop-filter-available-checkbox").ariaChecked;
    let minprice = priceRange1.value;
    let maxprice = priceRange2.value;
    let sort = sortSelect.value;

    minprice /= Number(localStorage.getItem("exchangeRate"))
    maxprice /= Number(localStorage.getItem("exchangeRate"))

    document.querySelectorAll(".shop-filter-section-checkbox").forEach((item) => {
        if (item.ariaChecked === "true") {
            sectionIds.push(Number(item.id));
        }
    })
    document.querySelectorAll(".shop-filter-producer-checkbox").forEach((item) => {
        if (item.ariaChecked === "true") {
            producerIds.push(Number(item.id));
        }
    })
    document.querySelectorAll(".shop-filter-designer-checkbox").forEach((item) => {
        if (item.ariaChecked === "true") {
            designerIds.push(Number(item.id));
        }
    })

    fetch(`/products/listbyfilter?sections=${sectionIds.join(",")}&producers=${producerIds.join(",")}&designers=${designerIds.join(",")}&available=${available}`)
        .then((res) => res.json()).then((data) => {

            data = data.filter(product => product.type === localStorage.getItem("type"))

            let minprice = Math.min(...data
                .map((item) => item.type === type ? item.price : undefined)
                .filter(price => price !== undefined)
            );

            let maxprice = Math.max(...data
                .map(item => item.type === type ? item.price : undefined)
                .filter(price => price !== undefined)
            );
            if (sort === "increasing") data.sort((a, b) => a.price - b.price)
            else if (sort === "decreasing") data.sort((a, b) => b.price - a.price)
            else if (sort === "abc") data.sort((a, b) => {
                if (a.name < b.name) {
                    return -1;
                }
            })
            console.log(data)
            CreateProductNodes(data, localStorage.getItem("exchangeRate"), true)

            ConfigurePriceRange(minprice, maxprice)
        })
}

function ListProducts(type) {
    fetch(`/products/list1`)
        .then((res) => res.json()).then((data) => {

            data = data.filter(product => product.type === type)
            let minprice = Math.min(...data
                .map((item) => item.type === type ? item.price : undefined)
                .filter(price => price !== undefined)
            );

            let maxprice = Math.max(...data
                .map(item => item.type === type ? item.price : undefined)
                .filter(price => price !== undefined)
            );
            // data.sort((a, b) => a.name.localeCompare(b.name));

            CreateProductNodes(data, localStorage.getItem("exchangeRate"), true)
            SecondListProducts(type, minprice, maxprice)
        })
}

function SecondListProducts(type, minprice2, maxprice2) {
    fetch(`/products/list2`)
        .then((res) => res.json()).then((data) => {

            data = data.filter(product => product.type === type)
            let minprice = Math.min(...data
                .map((item) => item.type === type ? item.price : undefined)
                .filter(price => price !== undefined)
            );

            let maxprice = Math.max(...data
                .map(item => item.type === type ? item.price : undefined)
                .filter(price => price !== undefined)
            );
            // data.sort((a, b) => a.name.localeCompare(b.name));

            CreateProductNodes(data, localStorage.getItem("exchangeRate"), false)

            ConfigurePriceRange(Math.min(minprice, minprice2), Math.max(maxprice, maxprice2))
        })
}

function CreateProductNodes(data, exRate, clearCont) {
    let itemsContainer = document.querySelector(".shop-items-all");
    if (clearCont) {
        itemsContainer.innerHTML = "";
    }
    document.querySelector(".shop-no-result").classList.add("hidden")
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
                    <p class="shop-item-price">${formatPrice(String((product.price * exRate).toFixed(2)))}</p>
                    <span class="shop-item-price-label">${pricelabel}</span>
                </div>
            </div></a>`;
        if (product.onSale) {
            let productPriceContainer = item.querySelector(".shop-item-price-container");
            productPriceContainer.children[0].classList.add("shop-item-line-through");
            productPriceContainer.children[1].classList.add("shop-item-line-through");
            productPriceContainer.innerHTML += `
                <p class="shop-item-sale-price">${formatPrice(String((product.price * exRate * ((100 - product.sale) / 100)).toFixed(2)))}</p>
                <span class="shop-item-sale-price-label">${pricelabel}</span>
            `;
        }

        itemsContainer.appendChild(item);
    })
    if (itemsContainer.children.length == 0) {
        document.querySelector(".shop-no-result").classList.remove("hidden")
    }
}

function ConfigurePriceRange(minprice, maxprice) {
    let ex = localStorage.getItem("exchangeRate")
    if (minprice == Infinity || maxprice == Infinity) {
        minprice = 0;
        maxprice = 0;
    }
    priceLabel1.innerHTML = (minprice * ex) + " ₽";
    priceLabel2.innerHTML = (maxprice * ex) + " ₽";
    priceRange1.min = minprice * ex;
    priceRange2.min = minprice * ex;
    priceRange1.max = maxprice * ex;
    priceRange2.max = maxprice * ex;
    priceRange1.value = minprice * ex;
    priceRange2.value = maxprice * ex;
}



// EX RATE

let exRateText = document.querySelector(".shop-ex-rate");

window.onload = () => {
    GetExRate();
    if (window.location.href.includes("furniture")) {
        localStorage.setItem("type", "furniture");
    } else if (window.location.href.includes("light")) {
        localStorage.setItem("type", "light");
    } else if (window.location.href.includes("interior")) {
        localStorage.setItem("type", "interior");
    }
    GetSections();
    GetProducersAndDesigners();

    if (availableParam){
        ShowOnlyAvailable()
    }else{
        ListProducts(localStorage.getItem("type"))
    }
    
}


function GetExRate() {
    fetch("/ex-rate")
        .then((res) => res.json())
        .then((ex) => {
            localStorage.setItem("exchangeRate", ex.exchangeRate);
            //exRateText.innerHTML = `Курс: 1 € = ${ex.exchangeRate} ₽`;
        })
}




// FILTERS



function GetSections() {
    fetch("/sections/list")
        .then((res) => res.json()).then((data) => {
            let dataAll = [];
            data.forEach((item) => {
                if (item.type === localStorage.getItem("type")) dataAll.push(item);
            })
            CreateFilterNodes(dataAll, "section");
        })
}
function GetProducersAndDesigners() {
    fetch("/producers/list")
        .then((res) => res.json()).then((data) => {
            CreateFilterNodes(data, "producer");
        })
    fetch("/designers/list")
        .then((res) => res.json()).then((data) => {
            CreateFilterNodes(data, "designer");
        })
}

function CreateFilterNodes(data, filtertype) {
    let catgoriesContainer = document.querySelector(`.shop-filter-${filtertype}s`);
    catgoriesContainer.innerHTML = "";
    data.forEach((item) => {
        catgoriesContainer.innerHTML +=
            `<li class="shop-filter-${filtertype}-li">
                <div class="shop-filter-${filtertype}-checkbox" aria-checked="false" id="${item.id}">
                    <div class="shop-filter-section-checkbox-bg hidden"></div>
                </div>
                <p class="shop-filter-section-label">${item.name}</p>
            </li>`;
    })
    for (let i = 0; i < catgoriesContainer.children.length; i++) {
        let item = catgoriesContainer.children[i];
        item.addEventListener("click", () => {
            let checkbox = item.querySelector(`.shop-filter-${filtertype}-checkbox`);
            checkbox.ariaChecked = Boolean(checkbox.ariaChecked !== "true");
            checkbox.querySelector(".shop-filter-section-checkbox-bg").classList.toggle("hidden");
            priceRange1.value = 0;
            priceRange2.value = 10000;
            if (exIsRub) ListProductsWithFilter("RUB", localStorage.getItem("type"));
            else ListProductsWithFilter("EUR", localStorage.getItem("type"));
        })
    }
}
ConfigureAvailableFilter()
function ConfigureAvailableFilter() {
    let checkbox = document.querySelector(`.shop-filter-available-checkbox`);
    document.querySelector(".shop-filter-available-li").addEventListener("click", () => {
        checkbox.ariaChecked = Boolean(checkbox.ariaChecked !== "true");
        checkbox.querySelector(".shop-filter-available-checkbox-bg").classList.toggle("hidden");

        ListProductsWithFilter("RUB", localStorage.getItem("type"));
    })
}


function debounce(func, delay) {
    let timeoutId;
    return function (...args) {
        if (timeoutId) {
            clearTimeout(timeoutId);
        }
        timeoutId = setTimeout(() => {
            func.apply(this, args);
        }, delay);
    };
}



function formatPrice(numberString) {
    const numericString = numberString.replace(/[^\d.]/g, '');
    let [integerPart, decimalPart] = numericString.split('.');
    integerPart = integerPart.replace(/\B(?=(\d{3})+(?!\d))/g, ' ');
    return decimalPart ? `${integerPart}.${decimalPart}` : integerPart;
}