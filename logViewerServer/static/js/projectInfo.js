function initData()
{
    var xmlhttp;
    if (window.XMLHttpRequest)
    {
        // IE7+, Firefox, Chrome, Opera, Safari 浏览器执行代码
        xmlhttp=new XMLHttpRequest();
    }
    else
    {
        // IE6, IE5 浏览器执行代码
        xmlhttp=new ActiveXObject("Microsoft.XMLHTTP");
    }

    // 请求后端接口
    xmlhttp.open("GET","/logviewer/log/pi",true);
    xmlhttp.send();
    xmlhttp.onreadystatechange=function()
    {
        if (xmlhttp.readyState ===4 && xmlhttp.status===200)
        {
            // 请求完成且请求成功时，将数据暂存在localStorage
            localStorage.setItem("projectData",xmlhttp.responseText)
        }
    }
}

initData()

// 然后将数据从localStorage取出，转换成json并返回
function Data(){
    const obj = JSON.parse(localStorage.getItem("projectData"));
    return obj["projectList"]
}

function getProjectList(){
    let projectHtml = "";

    console.log("项目清单：")
    console.log(Data())

    for(const p in Data()){
        projectHtml += '<option value=' + p + '>';
    }
    document.getElementById("project").innerHTML = projectHtml;
}

// 当鼠标离开项目选择框时，写入选择的项目名称到localStorage
function wProject(project){
    localStorage.setItem("selectProject",project)
}

// 依据localStorage中的项目信息获取已注册的环境信息
function getEnv(){
    const project =localStorage.getItem("selectProject")
    const envList = Data()[project]
    const projectEnv = []
    let envHtml = "";
    for(const e in envList){
        envHtml += '<option value=' + e + '>';
        projectEnv.push(e)
    }

    console.log("环境清单")
    console.log(projectEnv)

    localStorage.setItem("projectEnv",projectEnv.toString())
    document.getElementById("env").innerHTML = envHtml;
}

function wEnv(env){
    localStorage.setItem("selectEnv",env)
}


// 依据localStorage中的项目信息获取已注册的环境信息
function getSvc(){
    const project =localStorage.getItem("selectProject")
    const env =localStorage.getItem("selectEnv")
    const svcList = Data()[project][env]

    console.log("服务清单")
    console.log(svcList)

    let svcHtml = ""
    for (let i = 0; i < svcList.length ;i++) {
        svcHtml += '<option value=' + svcList[i] + '>';
    }
    document.getElementById("service").innerHTML = svcHtml;
}

function wSvc(svc){
    localStorage.setItem("selectSvc",svc)
}
