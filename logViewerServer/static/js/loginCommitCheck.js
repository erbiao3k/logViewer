function loginCommitCheck(){
    // 去除传入表单中的所有空格，判断传入值是否为空，为空则不允许提交表单
    if (logCommit.project.value.replace(/\s/gi,'') === ""){
        alert("请填写：项目名称");
        return false
    }
    if (logCommit.env.value.replace(/\s/gi,'') === ""){
        alert("请填写：环境类型");
        return false
    }
    if (logCommit.service.value.replace(/\s/gi,'') === ""){
        alert("请填写：服务名称");
        return false
    }
    if (logCommit.date.value.replace(/\s/gi,'') === ""){
        alert("请填写：日志时间");
        return false
    }

    let dateTime
    let yy = new Date().getFullYear()
    let mm = new Date().getMonth() + 1
    let dd = new Date().getDate()
    if (Number(mm) < 10){
        mm = "0" + mm
    }
    if (Number(dd) < 10){
        dd = "0" + dd
    }
    dateTime =yy + "-" + mm  + "-" + dd

    if (Number(logCommit.date.value.replace(/\-/g,'')) - Number(dateTime.replace(/\-/g,'')) > 0){
        alert(logCommit.date.value + "为未来日期")
        return false
    }

    if (new Date(dateTime) - new Date(logCommit.date.value) > 2592000000){
        alert("仅能拿取一个月内的日志")
        return false
    }
    return true;
}