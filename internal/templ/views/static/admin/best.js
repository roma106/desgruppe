document.addEventListener("DOMContentLoaded", () => {
	ConfigureUI()
	ConfigurePages(document.querySelectorAll(".items-container")[1].querySelectorAll(".items-li"))
	ConfigureTable();
})

function ConfigureUI() {
	let makeAvButtons = document.querySelectorAll(".items-li-make-sale");
	makeAvButtons.forEach(but => {
		but.addEventListener("click", ()=>{
			EditBest(but.id, true)
		})
	})
	let deleteButtons = document.querySelectorAll(".items-li-delete");
	deleteButtons.forEach(but => {
		but.addEventListener("click", ()=>{
			EditBest(but.id, false)
		})
	})
}


function EditBest(id, newStatus) {
	fetch("/admin/best/edit", {
		method: "PUT",
		body: JSON.stringify({ id: Number(id), best: Boolean(newStatus) })
	}).then(()=>window.location.reload())
}


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
