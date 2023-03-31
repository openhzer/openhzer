if (GetToken()){
    window.location.replace("index.html");
}
var app = new Vue({

})


function CaptchaCheck(){
    $.ajax({
        //请求方式
        type : "GET",
        //请求的媒体类型
        contentType: "application/json;charset=UTF-8",
        //请求地址
        url : "/api/captcha",
        //请求成功
        success : function(result) {
            let {error , data , msg} = result;
            if (error!==0){
                NotifyError(msg);
                return;
            }
            $("#captcha").attr("cid",data.id);
            $("#captchal").attr("src",data.data);
        },
        //请求失败，包含具体的错误信息
        error : function(e){

        }
    });
}
$('#captchal').click(CaptchaCheck);
$('#loginb').click(function () {
        let capobj = $("#captcha");
        let user = $('#username').val();
        let passwd = $('#password').val();
        if (!capobj.val() || !user || !passwd){
            NotifyError("参数不完整，请检查后提交");
            return;
        }
        console.log(capobj);
        $.ajax({
            //请求方式
            type : "POST",
            //请求的媒体类型
            contentType: "application/json;charset=UTF-8",
            //请求地址
            url : "/api/admin/login",
            data:JSON.stringify({
                user_name:user,
                pass_word:passwd,
                captcha:{
                    id:capobj.attr('cid'),
                    value:capobj.val(),
                },
            }),
            //请求成功
            success : function(result) {
                let {error , data , msg} = result;
                if (error!==0){
                    NotifyError(msg);
                    return;
                }
                SetToken(data.token);
                NotifySuccess(msg);
                window.location.replace("index.html");
            },
            //请求失败，包含具体的错误信息
            error : function(e){

            }
        });
});

CaptchaCheck();

