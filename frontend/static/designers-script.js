

// GETTING PROUCTS

window.onload = () => {
    if(window.location.href.includes("designer")){
        ListDesigners("designer");
    }else if(window.location.href.includes("producer")){
        ListDesigners("producer");
    }
}

function ListDesigners(type) {
    fetch(`/${type}s/list`)
    .then((res) => res.json()).then((data) => {
        CreateDesignerNodes(data, type);
    })
}

function CreateDesignerNodes(data, type) {
    let itemsContainer = document.querySelector(`.${type}s-list`);
    itemsContainer.innerHTML = "";
    data.forEach((des) => {
        let item = document.createElement("li");
        item.classList.add(`${type}s-list-elem`);
        let photourl = URL.createObjectURL(new Blob([new Uint8Array(des.picture)]));
        
        item.innerHTML = 
            `<a href="/${type}/${des.slug}" class="designers-elem-a">
                <img src="${photourl}" alt="" class="${type}s-elem-img">
                <p class="designers-elem-name">${des.name}</p>
            </a>`;
        if(type==="producer"){
            item.innerHTML = 
                `<a href="/${type}/${des.slug}" class="designers-elem-a">
                    <img src="${photourl}" alt="" class="${type}s-elem-img">
                </a>`;

        }

        itemsContainer.appendChild(item);
    })
}


