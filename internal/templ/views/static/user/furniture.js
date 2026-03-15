let filterButtonCategories = document.querySelector(".filter-button-sections");
let filterCategoriesList = document.querySelector(".shop-filter-sections");

filterButtonCategories.addEventListener("click", () => {
    filterCategoriesList.classList.toggle("filter-sections-hidden");
    document.querySelector(".sections-filter-img").classList.toggle("shop-filter-img-opened");
})

let filterButtonDesigners = document.querySelector(".filter-button-designers");
let filterDesignersList = document.querySelector(".shop-filter-designers");

filterButtonDesigners.addEventListener("click", () => {
    filterDesignersList.classList.toggle("filter-designers-hidden");
    document.querySelector(".designers-filter-img").classList.toggle("shop-filter-img-opened");
})



let filterButtonproducers = document.querySelector(".filter-button-producers");
let filterproducersList = document.querySelector(".shop-filter-producers");

filterButtonproducers.addEventListener("click", () => {
    filterproducersList.classList.toggle("filter-producers-hidden");
    document.querySelector(".producers-filter-img").classList.toggle("shop-filter-img-opened");
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



// GETTING PROUCTS

function ListProductsWithFilter(type) {
    let sectionIds = [];
    let producerIds = [];
    let designerIds = [];
    let available = document.querySelector(".shop-filter-available-checkbox").ariaChecked;
    let sort = sortSelect.value;

    document.querySelectorAll(".shop-filter-section-checkbox").forEach((item) => {
        if (item.ariaChecked==="true") {
            sectionIds.push(Number(item.id));
        }
    })
    document.querySelectorAll(".shop-filter-producer-checkbox").forEach((item) => {
        if (item.ariaChecked==="true") {
            producerIds.push(Number(item.id));
        }
    })
    document.querySelectorAll(".shop-filter-designer-checkbox").forEach((item) => {
        if (item.ariaChecked==="true") {
            designerIds.push(Number(item.id));
        }
    })
    if (sectionIds.length==0 && designerIds.length==0 && producerIds.length==0 && available==="false" && sort === "abc"){
        window.location.href=`/${type}`
        return
    }
    if (available==="false") available=""
    
    let query = `?sections=${sectionIds.join(",")}&producers=${producerIds.join(",")}&designers=${designerIds.join(",")}&available=${available}&sort=${sort}`
    // console.log(sort)
    window.location.href=`/${type}filter${query}`
}

// FILTERS


document.addEventListener("DOMContentLoaded", ()=>{
    ConfigureFilterNodes("section")
    ConfigureFilterNodes("designer")
    ConfigureFilterNodes("producer")
    ConfigureAvailableFilter()
})

let urlParams = new URLSearchParams(document.location.search)

function ConfigureFilterNodes(filtertype) {
    let catgoriesContainer = document.querySelector(`.shop-filter-${filtertype}s`);
    let filterParams = urlParams.get(`${filtertype}s`)
    
    for (let i = 0; i < catgoriesContainer.children.length; i++) {
        let item = catgoriesContainer.children[i];
        let checkbox = item.querySelector(`.shop-filter-${filtertype}-checkbox`);
        if (filterParams!="" && filterParams!=null){
            let filterIds = Array.from(filterParams.split(","))
            filterIds.forEach((filterId)=>{
                if (filterId===checkbox.id){
                    checkbox.ariaChecked = true;
                    checkbox.querySelector(`.shop-filter-${filtertype}-checkbox-bg`).classList.remove("hidden");
                }
            })
        }
        item.addEventListener("click", () => {
            checkbox.ariaChecked = Boolean(checkbox.ariaChecked !== "true");
            checkbox.querySelector(`.shop-filter-${filtertype}-checkbox-bg`).classList.toggle("hidden");
            ListProductsWithFilter(localStorage.getItem("type"));
        })
    }
}

function ConfigureAvailableFilter() {
    let checkbox = document.querySelector(`.shop-filter-available-checkbox`);
    let availableParam = urlParams.get("available")
    if (availableParam != "" && availableParam != null){
        checkbox.ariaChecked = true;
        checkbox.querySelector(".shop-filter-available-checkbox-bg").classList.remove("hidden");
    }
    document.querySelector(".shop-filter-available-li").addEventListener("click", () => {
        checkbox.ariaChecked = Boolean(checkbox.ariaChecked !== "true");
        checkbox.querySelector(".shop-filter-available-checkbox-bg").classList.toggle("hidden");
        ListProductsWithFilter(localStorage.getItem("type"));
    })
}




// Переход с главной страницы от кнопки "Смотреть все" начличие
// let urlParams = new URLSearchParams(document.location.search)
// let availableParam = Boolean(urlParams.get("available"))


// function ShowOnlyAvailable(){
//     document.querySelector(".shop-filter-available-checkbox").ariaChecked = true;
//     document.querySelector(".shop-filter-available-checkbox-bg").classList.remove("hidden");
//     ListProductsWithFilter(localStorage.getItem("type"))
// }



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
