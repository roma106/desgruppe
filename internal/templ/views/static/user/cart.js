

function ConfigureCartPanelToggling() {
    if (!window.location.href.includes("cart")) {
        let cartHeaderButton = document.querySelector(".header-cart-container")
        let cartCloseButton = document.querySelector(".cart-close-button")
        cartHeaderButton.addEventListener("click", () => {
            localStorage.setItem("cartOpened", "true")
            ToggleCart()
        })
        cartCloseButton.addEventListener("click", () => {
            localStorage.setItem("cartOpened", "false")
            ToggleCart()
        })
        
        let cartBg = document.querySelector(".cart-bg");
        cartBg.addEventListener("click", () => {
            localStorage.setItem("cartOpened", "false")
            ToggleCart()
        })
        let cartPanel = document.querySelector(".cart-panel")
        
        if(localStorage.getItem("cartOpened")==="true"){
            ToggleCart()
        }
        function ToggleCart() {
            cartBg.classList.toggle("hidden")
            cartPanel.classList.toggle("cart-panel-opened")
            document.querySelector("body").classList.toggle("overflow-hidden")
        }
        document.querySelector(".go-to-cart-button").addEventListener("click", () => {
            window.location.href = "/cart"
            localStorage.setItem("cartOpened", "false")
        })
    }
}


document.addEventListener("DOMContentLoaded", () => {
    let cartID = getCartIDCookie()
    ConfigureCartPanelToggling()
    ConfigureAddingProductToCart(cartID)
    ConfigureCartProductQty(cartID)
    ConfigureCartProductDeleting(cartID)
})


function ConfigureAddingProductToCart(cartID) {
    if (!window.location.href.includes("product")) return
    let addToCartButtons = document.querySelectorAll(".product-buy-button")
    addToCartButtons.forEach((but) => {
        but.addEventListener("click", () => {
            let qty = but.parentElement.childNodes[0].querySelector(".product-qty").innerHTML;
            AddProductToCart(cartID, but.id, qty)
        })
    })
}

function AddProductToCart(cartID, prID, qID) {
    let data = {
        "id": Number(cartID),
        "productID": Number(prID),
        "quantitiesID": Number(qID)
    }
    if (document.querySelector(".product-color-item-selected")) {
        data.colorID = Number(document.querySelector(".product-color-item-selected").id)
    } else {
        if (!document.querySelector(".product-color-item")){
            data.colorID = null
        }else{
            document.querySelector(".product-select-color-error").classList.remove("hidden");
            return
        }
    }
    
    fetch("/cart/product-add", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(() => {
        localStorage.setItem("cartOpened", "true")
        window.location.reload();
    })
}


function ConfigureCartProductQty(cartID) {
    let cartItemsQtyLess = document.querySelectorAll(".cart-item-qty-less");
    let cartItemsQtyMore = document.querySelectorAll(".cart-item-qty-more");
    for (let i = 0; i < cartItemsQtyLess.length; i++) {
        let node = cartItemsQtyLess[i]
        node.addEventListener("click", () => {
            let qty = node.parentElement.querySelector(".cart-item-qty-cnt")
            if (Number(qty.innerHTML) > 1) {
                EditCartNode(cartID, Number(qty.innerHTML) - 1, qty.id, i, qty)
            }
        })
    }
    for (let i = 0; i < cartItemsQtyMore.length; i++) {
        let node = cartItemsQtyMore[i]
        node.addEventListener("click", () => {
            let qty = node.parentElement.querySelector(".cart-item-qty-cnt")
            EditCartNode(cartID, Number(qty.innerHTML) + 1, qty.id, i, qty)
        })
    }
}

function EditCartNode(cartID, newQty, prID, index) {
    let data = {
        "id": Number(cartID),
        "productID": Number(prID),
        "quantitiesID": newQty,
        "index": index
    }
    console.log(data)
    fetch("/cart/product-edit", {
        method: "PUT",
        body: JSON.stringify(data)
    }).then(() => {
        localStorage.setItem("cartOpened", "true")
        window.location.reload();
    })
}

function ConfigureCartProductDeleting(cartID) {
    let productDeleteButtons = document.querySelectorAll(".cart-item-remove-button")
    for (let i=0;i<productDeleteButtons.length;i++){
        productDeleteButtons[i].addEventListener("click", () => {
            DeleteCartNode(cartID, i)
        })
    }
}

function DeleteCartNode(cartID, index) {
    let data = {
        "id": Number(cartID),
        "index": Number(index),
    }
    fetch("/cart/product-delete", {
        method: "DELETE",
        body: JSON.stringify(data)
    }).then(() => {
        localStorage.setItem("cartOpened", "true")
        window.location.reload();
    })
}

function getCartIDCookie() {
    const parts = `; ${document.cookie}`.split("; cartID=")
    if (parts.length == 2)
        return parts.pop().split(";").shift()
}