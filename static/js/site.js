jQuery.extend({
    alert: function (dom, result) {
        if (result.status != "success") {
            html = '<div class="alert alert-danger">' +
                ' <a class="close" data-dismiss="alert">×</a> ' +
                '<strong>'+result.status+'!</strong>  ' + result.message
            '</div>'
            $("#"+dom).html(html)
            //  setTimeout(function(){
            //     $("#"+dom).slideUp()
            //     $("#"+dom).html("")
            //     $("#"+dom).slideDown()
            // },5000)
        }else{
            console.log("alert success")
            html = '<div class="alert alert-success">' +
                ' <a class="close" data-dismiss="alert">×</a> ' +
                '<strong>'+result.status+'!</strong>  ' + result.message
            '</div>'
            $("#"+dom).html(html)
             setTimeout(function(){
                $("#"+dom).slideUp()
                $("#"+dom).html("")
                $("#"+dom).slideDown()
            },5000)
        }
    },
    alertFail: function (dom, message) {
            html = '<div class="alert alert-danger">' +
                ' <a class="close" data-dismiss="alert">×</a> ' +
                '<strong>Fail!</strong>  ' + message
            '</div>'
            $("#"+dom).html(html)
            // setTimeout(function(){
            //     $("#"+dom).slideUp()
                // $("#"+dom).html("")
                // $("#"+dom).slideDown()
            // },5000)
    },
    get:function(token,user,url,callback){
     $.ajax({
                url:url,
                type:"GET",
                headers: { 
                  "token" : token,
                  "user": user
                },
                contentType:"application/json; charset=utf-8",
                success:callback,
                error:callback,
              })
    },
    post:function(token,user,url,data,callback){
      $.ajax({
                url:url,
                type:"POST",
                headers: { 
                  "token" : token,
                  "user" : user
                },
                contentType:"application/json; charset=utf-8",
                data:data,
                dataType:"json",
                success:callback,
                error:callback,
              })
    },
    delete:function(token,user,url,callback){
        $.ajax({
                url:url,
                type:"DELETE",
                headers: { 
                  "token" : token,
                  "user": user
                },
                contentType:"application/json; charset=utf-8",
                success:callback,
                error:callback,
              })
    }

})

$('#logout').click(function(){
    $.post("","","/logout",null,function(result){
        console.log("result:",result);
        if(result.status == "success"){

            console.log("debug, logout success.")
            window.location.href="/"
        }else{
          $.alert("alert",result)
        }
    })
})
    var patterns = new Object();

    //匹配ip地址
    patterns.ip = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])(\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])){3}$/;

    //匹配邮件地址
    patterns.email = /^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$/;

    //匹配日期格式2008-01-31，但不匹配2008-13-00
    patterns.date = /^\d{4}-(0?[1-9]|1[0-2])-(0?[1-9]|[1-2]\d|3[0-1])$/;
    
    /*匹配时间格式00:15:39，但不匹配24:60:00，下面使用RegExp对象的构造方法
    来创建RegExp对象实例，注意正则表达式模式文本中的“\”要写成“\\”*/
    patterns.time = new RegExp("^([0-1]\\d|2[0-3]):[0-5]\\d:[0-5]\\d$");
        
    /*verify – 校验一个字符串是否符合某种模式
     *str – 要进行校验的字符串
     *pat – 与patterns中的某个正则表达式模式对应的属性名称
     */
    function verify(str,pat)
    {        
        thePat = patterns[pat];
        if(thePat.test(str))
        {
            return true;
        }
        else
        {
            return false;
        }
    }