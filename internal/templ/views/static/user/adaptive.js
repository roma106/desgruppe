let filtersToggleButton = document.querySelector(".shop-toggle-filters");
let filtersPanel = document.querySelector(".shop-filters-container")

let headerBurger = document.querySelector(".header-burger-container");
let headerNav = document.querySelector(".header-nav")

let screenWidth = window.innerWidth;

if (screenWidth <= 630) {
    // header
    headerNav.classList.add("header-nav-hidden")
    headerBurger.addEventListener("click", ()=>{
        headerNav.classList.toggle("header-nav-hidden")
    })


    // filters
    filtersPanel.classList.add("shop-filters-container-hidden");
    filtersToggleButton.addEventListener("click", ()=>{
        filtersPanel.classList.toggle("shop-filters-container-hidden")
    })
}

