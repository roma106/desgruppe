
let adminSendButton = document.querySelector(".login-send-button");

adminSendButton.addEventListener("click", LoginAdmin)

function LoginAdmin(){
    let data = {
        login: document.querySelector(".admin-input-login").value,
        password: document.querySelector(".admin-input-password").value
    };

    fetch("/auth-des-admin",{
        method: "POST",
        body: JSON.stringify(data),
    }).then((res) =>{
        if(res.ok){
            window.location.href = "/admin/sections";
        }
    });
}