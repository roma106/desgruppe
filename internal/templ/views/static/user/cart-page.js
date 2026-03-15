document.addEventListener("DOMContentLoaded", ()=>{
    let cartID = getCartIDCookie();
    ConfigureCartPageNodes(cartID);
})

function ConfigureCartPageNodes(cartID) {
    let cartItemsQtyLess = document.querySelectorAll(".item-elem-qty-less");
    let cartItemsQtyMore = document.querySelectorAll(".item-elem-qty-more");
    for (let i = 0; i < cartItemsQtyLess.length; i++) {
        let node = cartItemsQtyLess[i]
        node.addEventListener("click", () => {
            let qty = node.parentElement.querySelector(".item-elem-qty-cnt")
            if (Number(qty.innerHTML) > 1) {
                EditCartNode(cartID, Number(qty.innerHTML) - 1, qty.id, i, qty)
            }
        })
    }
    for (let i = 0; i < cartItemsQtyMore.length; i++) {
        let node = cartItemsQtyMore[i]
        node.addEventListener("click", () => {
            let qty = node.parentElement.querySelector(".item-elem-qty-cnt")
            EditCartNode(cartID, Number(qty.innerHTML) + 1, qty.id, i, qty)
        })
    }
    let productDeleteButtons = document.querySelectorAll(".item-elem-remove-button")
    for (let i=0;i<productDeleteButtons.length;i++){
        productDeleteButtons[i].addEventListener("click", () => {
            DeleteCartNode(cartID, i)
        })
    }
}




// ОФОРМАЛЕНИЕ ЗАКАЗА

let sendButton = document.querySelector(".total-send-button");
let cartEndBg = document.querySelector(".cart-end-bg")
let cartEndPanel = document.querySelector(".cart-end-panel")
sendButton.addEventListener("click", SendOrder)

function SendOrder(){
    if(!Validate()) return 

    let data= {
        name: document.querySelector(".total-name").value,
        phone: document.querySelector(".total-phone").value,
        email: document.querySelector(".total-email").value,
        comment: document.querySelector(".total-comment").value,
        cartId: getCartIDCookie(),
    }
    fetch("/orders/add", {
        method: "POST",
        body: JSON.stringify(data)
    }).then(res => {
        if (res.ok){
            ToggleEndCartPanel()
        }
    })
}

function Validate(){
    let name = document.querySelector(".total-name").value;
    let phone = document.querySelector(".total-phone").value;
    let email = document.querySelector(".total-email").value;
    if (name.length<1 || phone.length<1  || email.length<1){
        document.querySelector(".total-wrong-data").classList.remove("hidden")
        setTimeout(()=>{
            document.querySelector(".total-wrong-data").classList.add("hidden")
        }, 6000)
        return false
    }
    return true
}

function ToggleEndCartPanel(){
    cartEndBg.classList.toggle("hidden")
    cartEndPanel.classList.toggle("cart-end-panel-hidden")
}

document.querySelector(".cart-end-close-button").addEventListener("click", ()=>window.location.reload())