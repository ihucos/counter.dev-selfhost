$(document).ready(function () {
    let active = window.location.href.endsWith("?sign-up") ? 2 : 1;
    $(".tabs").tabslet({
        active: active,
    });
});

simpleForm("#sign-in form[action='/login']", "/dashboard.html");

document.addEventListener("userloaded", () => {
    window.location.href = "dashboard.html";
});
