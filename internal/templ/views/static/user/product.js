

document.addEventListener("DOMContentLoaded", () => {
    ConfigureColors()
    ConfigureRecommendations()
    ConfigurePhotos()
})



function ConfigureColors() {
    let colorsContainer = document.querySelector(".product-colors-container");
    colorsContainer.childNodes.forEach(color => {
        color.addEventListener("click", () => {
            let colorItems = document.querySelectorAll(".product-color-item");
            colorItems.forEach((item) => {
                item.classList.remove("product-color-item-selected");
            })
            color.classList.add("product-color-item-selected");
        })
    })
}

function ConfigurePhotos() {
    let sliderWidth = 67;
    if (window.innerWidth <= 880){sliderWidth=56}
    if (window.innerWidth <= 450){sliderWidth=45}
    let slider = document.querySelector(".product-slider-all");
    // slider.style.marginLeft = "0px"
    let photosContainer = document.querySelector(".product-catalog-container");
    let photos = photosContainer.childNodes;
    let mphotoNode = document.querySelectorAll(".product-catalog-item")[0]
    mphotoNode.addEventListener("click", () => slider.style.marginLeft = "0vh")

    for (let i = 0; i < photos.length; i++) {
        let photoNode = document.querySelectorAll(".product-catalog-item")[i]
        photoNode.addEventListener("click", () => {
            slider.style.marginLeft = "-" + (Number(sliderWidth) * i) + "vh";
        })
    }

    let sliderLeft = document.querySelector(".product-photos-left-arrow");
    let sliderRight = document.querySelector(".product-photos-right-arrow");
    sliderLeft.addEventListener("click", () => {
        let sliderMargin = String(slider.style.marginLeft).slice(0, -2);
        if (Number(sliderMargin) >= 0) {
            return
        }
        slider.style.marginLeft = (Number(sliderMargin) + sliderWidth) + "vh";
    })
    sliderRight.addEventListener("click", () => {
        let sliderMargin = String(slider.style.marginLeft).slice(0, -2);
        if (("-" + (sliderWidth * (photos.length-1))) === sliderMargin) {
            return
        }
        slider.style.marginLeft = (Number(sliderMargin) - sliderWidth) + "vh";
    })
}



// PRODUCT QUANTITY

let productQtyLess = document.querySelector(".product-qty-less");
let productQtyMore = document.querySelector(".product-qty-more");
let productQty = document.querySelector(".product-qty");

productQtyLess.addEventListener("click", () => {
    let qty = Number(productQty.innerHTML);
    if (qty > 1) {
        productQty.innerHTML = qty - 1;
    }
})

productQtyMore.addEventListener("click", () => {
    let qty = Number(productQty.innerHTML);
    productQty.innerHTML = qty + 1;
})



// NAVIGATION

let navPanel = document.querySelector(".product-nav-panel");
let navIndicator = document.querySelector(".product-nav-indicator");
let navItems = navPanel.querySelectorAll(".product-nav-item");
let navInfos = document.querySelectorAll(".product-nav-info");

for (let i = 0; i < navItems.length; i++) {
    let item = navItems[i];
    item.addEventListener("click", () => {
        let offset = (navPanel.getBoundingClientRect().left - item.getBoundingClientRect().left);
        navIndicator.style.marginLeft = `${-offset}px`;
        navIndicator.style.width = `${item.getBoundingClientRect().width}px`;
        navItems.forEach((itemn) => {
            itemn.classList.remove("product-nav-item-active");
        })
        item.classList.add("product-nav-item-active");

        navInfos.forEach((info) => {
            info.classList.remove("product-nav-info-active");
        })
        navInfos[i].classList.add("product-nav-info-active");
    })
}




// recommendations

function ConfigureRecommendations() {
    let carouselWidth = document.querySelector(".product-recommendations-carousel").offsetWidth;
    let container = document.querySelector('.product-recommendations-container');
    let prevBtn = document.querySelector('.rec-left');
    let nextBtn = document.querySelector('.rec-right');
    let offset = 0;
    let elemwidth = document.querySelector(".product-rec-item").offsetWidth;
    let lastOffset = elemwidth * (container.children.length -1 - (carouselWidth / elemwidth));


    prevBtn.addEventListener("click", () => {
        if (offset <= 0) {
            return
        }
        offset -= elemwidth;
        container.style.marginLeft = -offset + "px";
    })

    nextBtn.addEventListener("click", () => {
        if (offset >= lastOffset) {
            return
        }
        offset += elemwidth;
        container.style.marginLeft = -offset + "px";
    })
}


function GoToProducerPage(slug) {
    window.location.href = `/producer/${slug}`
}