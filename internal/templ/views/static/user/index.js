document.querySelector(".about-button").addEventListener("click", () => {
    window.location.href = "/about"
})


// Кнопка "Смотреть все" наличие 

let availableButton = document.querySelector(".main-available-button")

availableButton.addEventListener("click", ()=>{
    // window.location.href = "/furniturefilter?sections=&producers=&designers=&available=true"
    window.location.href = "/available"
})



// BANNER

let currentIndex = 0;
const bannersContainer = document.querySelector('.main-banner-container');
const banners = document.querySelectorAll('.main-banner-elem');
const totalBanners = banners.length;

document.querySelector('.banner-left').addEventListener('click', () => {
    currentIndex = (currentIndex > 0) ? currentIndex - 1 : totalBanners - 1;
    updateBannerDisplay();
});

document.querySelector('.banner-right').addEventListener('click', () => {
    currentIndex = (currentIndex < totalBanners - 1) ? currentIndex + 1 : 0;
    updateBannerDisplay();
});
function updateBannerDisplay() {
    let width = document.querySelector(".main-banner-carousel").offsetWidth;
    const offset = currentIndex * (width * (-1)); // Высчитываем смещение в зависимости от индекса
    bannersContainer.style.transform = `translateX(${offset}px)`;
}

document.addEventListener("DOMContentLoaded", ()=>{
    updateBannerDisplay();
    ConfigureAvailable("available")
    ConfigureAvailable("sale")
})


function ConfigureAvailable(type) {
    let carouselWidth = document.querySelector(`.main-${type}-carousel`).offsetWidth;
    let container = document.querySelector(`.main-${type}-container`);
    let prevBtn = document.querySelector(`.${type}-left`);
    let nextBtn = document.querySelector(`.${type}-right`);

    let offset = 0;
    if (!document.querySelector(`.product-${type}-item`)){return}
    let elemwidth = document.querySelector(`.product-${type}-item`).offsetWidth;
    let lastOffset = elemwidth * (container.children.length - (carouselWidth / elemwidth));


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


