
function ListColors(callback) {
    fetch("/admin/colors/list", {
        method: "GET",
    }).then((res) => {
        if (res.ok) {
            res = res.json();
            res.then((data) => {
                CreateColorItems(data)
            })
        } else {
            alert("Неудалось получить список цветов");
        }
    }).then(() => {
        if (callback && typeof callback === 'function') {
            callback();
        }
    })
}

let selectedColorsBorder = "solid 2px #0A6CCE";
function CreateColorItems(data) {
    let colorsContainer = document.querySelector(".add-form-colors-container-all");
    let table = CreateColorInputNodes(data);
    colorsContainer.innerHTML = "";
    table.forEach((color) => {
        colorsContainer.appendChild(color);
    })
    let checkboxes = document.querySelectorAll(".add-form-color-checkbox");
    let selectedColorsContainer = document.querySelector(".add-form-selected-colors")
    for (let i = 0; i < checkboxes.length; i++) {
        checkboxes[i].addEventListener("change", () => {
            if (checkboxes[i].checked) {
                selectedColorsContainer.appendChild(table[i].cloneNode(true))
                checkboxes[i].parentElement.style.border = selectedColorsBorder;
            } else {
                checkboxes[i].parentElement.style.border = "solid 1.5px gray";
                selectedColorsContainer.childNodes.forEach((col) => {
                    if (checkboxes[i].id === col.childNodes[0].id) {
                        col.remove();
                    }
                })
            }

        })

    }

}

function CreateColorInputNodes(data) {
    let table = [];
    data.forEach((color) => {
        let newcolor = document.createElement("div");
        if (color.code != "") {
            newcolor.className = `add-form-color`;
            newcolor.innerHTML = (`<input type="checkbox" class="add-form-color-checkbox" id="${color.id}">
                <span class="add-form-color-span"
                    style="background-color: #${color.code};"></span>
                <p class="add-form-color-p">${color.name}</p>`)
            table.push(newcolor);
        } else if (color.code == "") {
            let picturebytes = new Int8Array(color.picture);
            // console.log(picturebytes);
            const blob = new Blob([picturebytes], { type: "image/png" });
            const url = URL.createObjectURL(blob);
            newcolor.className = `add-form-color`;
            newcolor.innerHTML = (`<input type="checkbox" class="add-form-color-checkbox" id="${color.id}">
                    <span class="add-form-color-span"
                        style="background-image: url(${url});"></span>
                    <p class="add-form-color-p">${color.name}</p>`)
            table.push(newcolor);
        }
    })
    return table;
}

function GetCheckedColors() {
    let colors = [];
    let checkboxes = document.querySelectorAll(".add-form-color-checkbox:checked");
    checkboxes.forEach(checkbox => {
        colors.push(Number(checkbox.id));
    })
    return Array.from(new Set(colors));
}





function ListDesigners(callback) {
    fetch("/admin/designers/list", {
        method: "GET",
    }).then((res) => {
        if (res.ok) {
            res = res.json();
            res.then((data) => {
                let select = document.querySelector(".add-designer")
                CreateSelectInputNodes(select, data);
            })
        } else {
            alert("Неудалось получить список дизайнеров");
        }
    }).then(() => {
        if (callback && typeof callback === 'function') {
            callback();
        }
    })
}

function ListProducers(callback) {
    fetch("/admin/producers/list", {
        method: "GET",
    }).then((res) => {
        if (res.ok) {
            res = res.json();
            res.then((data) => {
                let select = document.querySelector(".add-producer")
                CreateSelectInputNodes(select, data);
            })
        } else {
            alert("Неудалось получить список фабрик");
        }
    }).then(() => {
        if (callback && typeof callback === 'function') {
            callback();
        }
    })
}

let selectTypeInput = document.querySelector(".add-type")
selectTypeInput.addEventListener("change", () => {
    ListSections(selectTypeInput.value)
})

function ListSections(type, callback) {
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
    }).then(() => {
        if (callback && typeof callback === 'function') {
            callback();
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