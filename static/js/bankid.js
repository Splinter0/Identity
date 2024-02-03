let sameDeviceButton = document.getElementById("sameDevice");
let otherDeviceButton = document.getElementById("otherDevice");
let cancelButton = document.getElementById("cancel");
let extendButton = document.getElementById("extend");
let manualLink = document.getElementById("manualLink");
let qrCodeElement = document.getElementById("qrCode");
function setMessage(message) {
    document.getElementById("statusMessage").textContent = message;
}
function renderQrCode(data) {
    console.log(data);
    qrCodeElement.innerHTML = "";
    var qrcode = new QRCode(qrCodeElement, {
        text: data,
        width: 180,
        height: 180,
        colorDark: "#000000",
        colorLight: "#ffffff",
    });
}
function collect() {
    count = 0;
    var timer = setInterval(
        async () => {
            await fetch("/status")
                .then((r) => r.json())
                .then((json) => {
                    count++;
                    console.log(json);
                    if (json.status === "failed") {
                        setMessage(json.message);
                        clearTimeout(timer);
                        cancelButton.style.display = "none";
                        extendButton.style.display = "none";
                        qrCodeElement.style.display = "none";
                        manualLink.style.display = "none";
                    } else if (json.status === "complete")  {
                        setMessage("");
                        clearTimeout(timer);
                        document.getElementById("userData").textContent = json.data.user.name;
                    } else {
                        message = json.message;
                        if (json.data && json.data.qrData) {
                            renderQrCode(json.data.qrData);
                            if (count >= 25) {
                                extendButton.style.display = "block";
                                message += "\nQR code is about to expire, extend if you wish to continue";
                            }
                        }
                        setMessage(message);
                    }
                })
        },
        1000
    );
    cancelButton.style.display = "block";
    cancelButton.addEventListener("click", () => {
        clearTimeout(timer);
        cancel();
    });
    extendButton.addEventListener("click", () => {
        clearTimeout(timer);
        extend();
    });
}
function start(same) {
    fetch(
        `/start?same=${same}`, 
        {
            method: "POST",
            headers: {
                "X-BankID-CSRF": "1",
            },
        },
    )
        .then((r) => r.json())
        .then((json) => {
            sameDeviceButton.style.display = "none";
            otherDeviceButton.style.display = "none";
            if (json.qrCodeData) {
                renderQrCode(json.qrCodeData);
            } else {
                window.location = json.launchUrl;
                manualLink.style.display = "block";
                manualLink.href = json.launchUrl;
            }
            collect();
        });
}
function cancel() {
    fetch(
        `/cancel`, 
        {
            method: "POST",
            headers: {
                "X-BankID-CSRF": "1",
            },
        },
    ).then(() => {
        cancelButton.style.display = "none";
        extendButton.style.display = "none";
        qrCodeElement.style.display = "none";
        manualLink.style.display = "none";
        setMessage("Authentication was cancelled, refresh the page to try again");
    });
}
function extend() {
    cancelButton.style.display = "none";
    extendButton.style.display = "none";
    start(0);
}
sameDeviceButton.addEventListener("click", () => start(1));
otherDeviceButton.addEventListener("click", () => start(0));