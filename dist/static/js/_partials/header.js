"use strict";

// onValidate is that verify the inputs
var onValidate = function onValidate(username, password) {
    "use strict";

    var usernameValidator = /^.[a-zA-Z0-9_].{1,16}$/,
        // 保障用户名必须伟大小写字母以及数字下划线组成，长度控制在８~16
    passwordValidator = /^(?=.*[A-Z])(?=.*[0-9])(?=.*[a-z]).{7,32}$/; // 保障密码至少一个大写字母，一个数字，一个小写字母，密码长度在7~32之间

    if (usernameValidator.test(username) !== true) {
        return {
            msg: "failed, illegal username, username regulation shows above:" + "\n 1. only accept upper case / lower case / digit / underscore character" + "\n 2. the username's length between 8 to 16",
            result: usernameValidator.test(username)
        };
    }
    if (passwordValidator.test(password) !== true) {
        return {
            msg: "failed, illegal password, password regulation shows above:" + "\n 1. only permitted upper case / lower case / digit " + "\n 2. at least contain one upper case / lower case / digit / underscore character" + "\n 3. the password's length between 7 to 32",
            result: passwordValidator.test(password)
        };
    }
    return {
        msg: "success",
        result: true
    };
};

// login is that login logical action
$(".login.form").submit(function () {
    "use strict";

    var username = $("#login-username").val(),
        password = $("#login-password").val();
    var result = onValidate(username, password);

    if (result.result === true) {
        $.ajax({
            url: "/user/login",
            dataType: "json",
            contentType: "application/json;charset=utf-8",
            type: "POST",
            data: JSON.stringify({
                username: username,
                password: password
            }),
            success: function success(data, status) {
                alert("success, data: " + data + ", status: " + status);
            }
        }).fail(function (response) {
            alert("failed: data: " + response.responseText + ", status: " + response.statusText);
        });
    } else {
        alert(result.msg + ", value: " + password);
    }
    return false;
});

// register is action that to redirect to register page
$("#register-click").click(function () {
    "use strict";

    $.ajax({
        url: "/user/register",
        type: "GET",
        success: function success() {
            window.location.replace("/user/register");
        }
    }).always(function () {
        window.location.replace("/user/register");
    });
});
//# sourceMappingURL=header.js.map