jQuery.extend({
    alert: function (dom, result) {
        if (result.status != "success") {
            html = '<div class="alert alert-danger">' +
                ' <a class="close" data-dismiss="alert">×</a> ' +
                '<strong>'+result.status+'!</strong>  ' + result.message
            '</div>'
            $("#"+dom).html(html)
             setTimeout(function(){
                $("#"+dom).html("")
            },4000)
        }
    },
    alertFail: function (dom, message) {
            html = '<div class="alert alert-danger">' +
                ' <a class="close" data-dismiss="alert">×</a> ' +
                '<strong>Fail!</strong>  ' + message
            '</div>'
            $("#"+dom).html(html)
            setTimeout(function(){
                $("#"+dom).html("")
            },4000)
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