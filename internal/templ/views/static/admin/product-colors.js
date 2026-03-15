
let selectedColorsBorder = "solid 2px #0A6CCE";
function ColorBorderChange(inputID) {
    let colorNode = document.querySelector(".add-form-colors-container-all").querySelector(`.add-form-color-checkbox#color${inputID}`)
    let selectedColorsContainer = document.querySelector(".add-form-selected-colors")
    if (colorNode.checked) {
        selectedColorsContainer.appendChild(colorNode.parentElement.cloneNode(true))
        colorNode.parentElement.style.border = selectedColorsBorder;
    } else {
        colorNode.parentElement.style.border = "solid 1.5px gray";
        selectedColorsContainer.childNodes.forEach((col) => {
            if (colorNode.id === col.childNodes[0].id) {
                col.remove();
            }
        })
    }

}


function GetCheckedColors() {
    let colors = [];
    let checkboxes = document.querySelectorAll(".add-form-color-checkbox:checked");
    checkboxes.forEach(checkbox => {
        colors.push(Number(checkbox.id.slice(5)));
    })
    return Array.from(new Set(colors));
}


let selectTypeInput = document.querySelector(".add-type")
selectTypeInput.addEventListener("change", () => {
    ListSections(selectTypeInput.value)
})

function ListSections(type) {
    fetch("/admin/sections/list", {
        method: "GET",
    }).then((res) => {
        if (res.ok) {
            res = res.json();
            res.then((data) => {
                let dataWithType = []
                data.forEach(elem => {
                    if (elem.type === type) {
                        dataWithType.push(elem)
                    }
                })
                let select = document.querySelector(".add-section")
                CreateSelectInputNodes(select, dataWithType);
            })
        } else {
            alert("Неудалось получить список разделов");
        }
    })
}

function CreateSelectInputNodes(input, data) {
    input.innerHTML = "";
    data.forEach((node) => {
        let newOpt = document.createElement("option");
        newOpt.className = `add-form-opt`;
        newOpt.value = node.id;
        newOpt.innerHTML = node.name;
        input.appendChild(newOpt);
    })
}


// SALE CHECKBOX 

const saleCheckbox = document.querySelector(".add-sale-checkbox");
const saleContainer = document.querySelector(".add-form-sale-block");
saleCheckbox.addEventListener("change", () => {
    saleContainer.classList.toggle("hidden")
})




// SEARCH COLOR
let searchColorInput = document.querySelector(".search-color-input");
let searchColorButton = document.querySelector(".search-color-button");
let searchColorResultText = document.querySelector(".search-color-result-text");

searchColorInput.addEventListener("input", () => {
    if (searchColorInput.value.length > 0) {
        searchColorButton.classList.remove("search-color-button-disabled");
    } else {
        searchColorButton.classList.add("search-color-button-disabled");
    }
})

searchColorButton.addEventListener("click", () => {
    if (searchColorInput.value.length > 0) {
        SearchColor();
    } else {
        return
    }
})

function SearchColor() {
    fetch(`/admin/products/searchcolor?q=${searchColorInput.value}`,
        {
            method: "GET",
        }).then((res) => {
            console.log(res)
            if (res.ok) {
                return res.json()
            }
        }).then((items) => {
            console.log(items);
            searchColorResultText.innerHTML =
                `По запросу "${searchColorInput.value}" найдено ${items.length} результатов.`;
            searchColorResultText.classList.remove("display-hidden");

            let colorsContainer = document.querySelector(".add-form-colors-container-all")
            colorsContainer.childNodes.forEach((colorNode) => {
                colorNode.style.border = "solid 1.5px gray";
                items.forEach((searchedcolor) => {
                    if (String(searchedcolor.id) === colorNode.childNodes[0].id) {
                        colorsContainer.insertBefore(colorNode, colorsContainer.firstChild)
                        colorNode.style.border = "solid 2px green"
                    }
                })
            })
        })
}

let selectedColorsContainer = document.querySelector(".add-form-selected-colors");
const selectedColorsObserver = new MutationObserver((mutationsList) => {
    for (const mutation of mutationsList) {
        selectedColorsContainer.childNodes.forEach((selCol) => {
            if (selCol.querySelector("input")) {
                selCol.querySelector("input").style.display = "none";
            }
            selCol.style.border = "solid 1.5px gray";
        })
    }
})
selectedColorsObserver.observe(selectedColorsContainer, { childList: true })



// PRICE + SALE CONVERTATION

let salePercent = document.querySelector(".add-sale")
let priceWithSale = document.querySelector(".add-price-with-sale")
let price = document.querySelector(".add-price")

salePercent.addEventListener("input", () => {
    priceWithSale.value = Math.round(Number(price.value) - (Number(price.value) * (Number(salePercent.value) / 100)))
})
price.addEventListener("input", () => {
    priceWithSale.value = Math.round(Number(price.value) - (Number(price.value) * (Number(salePercent.value) / 100)))
})
priceWithSale.addEventListener("input", () => {
    salePercent.value = Math.round(100 - Number(priceWithSale.value) / Number(price.value) * 100)
})