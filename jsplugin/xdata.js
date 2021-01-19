top.xdata = "插件xdata初始化";
top.xdata_query = function (type,code,vals){
    var sql = squel.select()
			.from("x_data", "d")
            .where("d.type = '"+type+"'")
            .where("d.code = '"+code+"'");
    for(var i=0;i<vals.length;i++){
        sql.where("d.c"+i+" = '"+vals[i]+"'")
    }
    sql = sql.toString();

    var pr = top.db_exe("db",sql)
	if(pr.code == 200){
    	var d = eval("("+pr.data+")");
        return d;
    }
    return [];

}

top.xdata_add = function (type,code,vals){
    var sql = squel.insert()
                .into("x_data")
                .set("type",type)
                .set("code",code)
                .set("created_at",new Date().Format("yyyy-MM-dd hh:mm:ss"))
    for(var i=0;i<vals.length;i++){
        sql.set("c"+i,vals[i])
    }
    sql = sql.toString();

    var pr = top.db_exe("db",sql)
	if(pr.code == 200){
    	var d = eval("("+pr.data+")");
        return d;
    }
    return [];

}

top.xdata_update = function (type,code,vals){
    var sql = squel.update()
                .table("x_data")
                .set("updated_at",new Date().Format("yyyy-MM-dd hh:mm:ss"))
    for(var i=0;i<vals.length;i++){
        sql.set("c"+i,vals[i])
    }
    sql.where("type = '"+type+"'")
    .where("code = '"+code+"'");

    sql = sql.toString();

    var pr = top.db_exe("db",sql)
	return pr;

}
