var app = {
    "appid":"wxc5c96be3cb9b4652",
    "secret":"0926be326cb0ff53d39c52114831c1bc"
}

top.getToken = function (app){
    var get = false;
    if(!top.wechat){
        top.wechat = {}
    }
    if(!top.wechat[app.appid]){
        top.wechat[app.appid] = {}
    }
    if(!top.wechat[app.appid].expire){
        get = true;
    }
    
    if(top.wechat[app.appid].expire <= new Date().getTime()){
        get = true;
    }
        
    if(get){
        var pr = callp("wechat","1.0","AccessToken",app);
        if(pr.code == 200){
            top.wechat[app.appid].expire = new Date().getTime()+7000*1000;
            top.wechat[app.appid].token = pr.data;
            return pr.data;
        }
        return "";
    }else{
        return top.wechat[app.appid].token;
    }
}

top.followList = function (token){
    var json = {};
    json.token = token;
    var pr =callp("wechat","1.0","GetFollowList",json);
    if(pr.code == 200){
        var d = JSON.parse(pr.data)
        var arr = [];
        //return JSON.stringify(d.data);
        for(var i=0;i<d.count;i++){
            arr.push(d.data.openid[i])
        }
        return arr;
    }
    return "[]";
}