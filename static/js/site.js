jQuery.extend({
    alert: function (dom, result) {
        if (result.status != "success") {
            html = '<div class="alert alert-danger">' +
                ' <a class="close" data-dismiss="alert">×</a> ' +
                '<strong>'+result.status+'!</strong>  ' + result.message
            '</div>'
            $("#"+dom).html(html)
        }
    },
    alertFail: function (dom, message) {
            html = '<div class="alert alert-danger">' +
                ' <a class="close" data-dismiss="alert">×</a> ' +
                '<strong>Fail!</strong>  ' + message
            '</div>'
            $("#"+dom).html(html)
    }
})