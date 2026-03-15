document.addEventListener("DOMContentLoaded", ()=>{
    ConfigureSettingsUI()
})

function ConfigureSettingsUI(){
    let changeExRateButton = document.querySelector(".change-ex-rate");
    
    changeExRateButton.addEventListener("click", () => {
        let exRate = document.querySelector(".ex-rate-input").value;
        if (exRate.includes(",")) { exRate = exRate.replace(",", "."); }
        let data = {
            exchangeRate: Number(exRate)
        }
        EditSettings(data)
    })
    let changeEmailButton = document.querySelector(".change-email");
    
    changeEmailButton.addEventListener("click", () => {
        let data = {
            email: document.querySelector(".email-input").value
        }
        EditSettings(data)
    })
}

function EditSettings(data){
    fetch("/admin/ex-rate", {
        method: "PUT",
        body: JSON.stringify(data),
    }).then((res) => {
        if (res.ok) {
            alert("Настройки успешно изменены");
        }else{
            alert("Неудалось изменить настройки");
        }
    })
}




