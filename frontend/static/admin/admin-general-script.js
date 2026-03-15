ConfigureTable();


function ConfigureTable() {
    requestAnimationFrame(() => {
        ["position", "photo", "name", "producer", "designer", "size", "price", "discount"].forEach((tablecol) => {
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



// SEARCH
let searchInput = document.querySelector(".search-input");
let searchButton = document.querySelector(".search-button");
let searchResultText = document.querySelector(".search-result-text");

searchInput.addEventListener("input", () => {
    if (searchInput.value.length > 0) {
        searchButton.classList.remove("search-button-disabled");
    } else {
        searchButton.classList.add("search-button-disabled");
    }
})

searchButton.addEventListener("click", () => {
    if (searchInput.value.length > 0) {
        Search();
    } else {
        return
    }
})

function Search() {
    fetch(`${window.location.href}/search?q=${searchInput.value}`,
        {
            method: "GET",
        }).then((res) => {
            if (res.ok) {
                res.json().then((items) => {
                    searchResultText.innerHTML =
                        `По запросу "${searchInput.value}" найдено ${items.length} результатов.`;
                    searchResultText.classList.remove("display-hidden");
                    let table = CreateColorNodes(items);
                    ConfigurePages(table);
                    ConfigureTable();
                })
            }
        })
}


// PAGES

let pagesContainer = document.querySelector(".pages-container");
let itemsTable = document.querySelector(".items-table");

function ConfigurePages(items) {
    pagesContainer.innerHTML = "";
    const itemsOnPage = 50;
    let pagesRequired = Math.floor(items.length / itemsOnPage) + 1;
    for (let i = 0; i < pagesRequired; i++) {
        let page = document.createElement("p");
        page.classList.add("page-p");
        page.innerHTML = i + 1;
        pagesContainer.appendChild(page);
        let pageItems = [];
        for (let j = 0; j < itemsOnPage; j++) {
            pageItems.push(items[i * itemsOnPage + j]);
        }
        if (i == 0) {
            itemsTable.innerHTML = "";
            pageItems.forEach((item) => {
                if (item != undefined) {
                    itemsTable.appendChild(item);
                }
            })
        }
        page.addEventListener("click", () => {
            itemsTable.innerHTML = "";
            pageItems.forEach((item) => {
                if (item != undefined) {
                    itemsTable.appendChild(item);
                }
            })
            ConfigureTable();
            for (const pageelem of pagesContainer.children) {
                pageelem.classList.remove("selected-page");
            }
            page.classList.add("selected-page");
        })
    }
    pagesContainer.children[0].classList.add("selected-page");
}
