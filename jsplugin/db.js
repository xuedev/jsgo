top.db = "插件db初始化";
top.db_exe = function(key,sql){
    top.dbInit(key,top.get_db_url(key));
    var pr = top.dbQuery(key,sql);
    return pr;
}
top.db_page = function(key,sql,curPage,pageSize){
    top.dbInit(key,top.get_db_url(key));
    var csql = "select count(1) as count from ("+sql+") t";
    var pr = top.dbQuery(key,csql);
    
    var count = 0;
    if(pr.code == 200){
        var cj = eval("("+pr.data+")"); 
        count = cj;
      
    	sql = sql+" limit "+(curPage-1)*pageSize+","+pageSize;
        var pr = top.dbQuery(key,sql);
        if(pr.code != 200){
        	return pr;
        }else{
        	var data = {};
            data.total = count[0].count-0;
            data.rows = eval("("+pr.data+")");
            return data;
        }
        
    }else
      return pr;
}

top.get_db_url = function(key){
    if(!top.db_urls){
        top.db_urls = {}
    }

    if(top.db_urls[key]){
        return top.db_urls[key];
    }else{//初始化db url字符串
        var sql = squel.select()
            .from("x_data", "d")
            .where("d.type = 'database'")
			.where("d.code = '"+key+"'")
			.toString();
        //{"sql":"update sys_user set nick_name='xuegx' "}
        var url = "root:gxdb@tcp(127.0.0.1:3306)/gxapi?charset=utf8&parseTime=True&loc=Local&timeout=1000s";
        var pr = top.dbInit("db",url)
        var pr = top.dbQuery("db",sql);
        if(pr.code == 200){
            var arr = eval("("+pr.data+")");
            if(arr.length>0){
                var d = arr[0];
                var url = d.c1+":"+d.c2+"@tcp("+d.c0+")/gxapi?charset=utf8&parseTime=True&loc=Local&timeout=1000s";
                top.db_urls[key] = url;
                return url;
            }
            return pr;
        }
        return pr;
    }
}

