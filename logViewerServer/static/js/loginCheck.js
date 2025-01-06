function loginCheck(){
    // 去除传入表单中的所有空格，判断传入值是否为空，为空则不允许提交表单
    if (loginPage.username.value.replace(/\s/gi,'') === ""){
        alert("账号信息为空！！！");
        return false
    }
    if (loginPage.password.value.replace(/\s/gi,'') === ""){
        alert("密码信息为空！！！");
        return false
    }
    const start = loginPage.name.value.length - "@myemal.com".length;
    const arr = loginPage.name.value.substr(start, "@myemal.com".length);
    if(arr !== "@myemal.com"){
        alert("请使用myemal.com邮箱登陆！！！")
        return false;
    }
    return true;
}