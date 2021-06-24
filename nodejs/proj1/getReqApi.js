
//module.exports ={
  //      getAPI
    //}

const getAPI = function(){

    const request = require('request');

    var token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJhc3Nob2xsZSIsImV4cCI6MTYyNDA3MzI0NiwiaXNzIjoid2VibGlua19hY2Nlc3MifQ.2x8Ke-DCgL_8mJQw5rahYt_htgAEYv7Aj0Dld3lHDbI';

    var options = {
      method: 'GET',
      hostname: '127.0.0.1',
      port:'8000',
      path: '/links/all',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer '+token,
      },
      json : true,

    };

    var js = 'haha';

    console.log('httpcall');

    var url = 'http://127.0.0.1:8000/links/all';

    request(url, options, function(error, response, body) {
        // субфукция получает респонз асинхронно
        // когда она получит и куда кидать инфу непонятно
        if (error) {
            return  console.log(error)
        };
        js = body
        if (!error && response.statusCode == 200) {
        // do something with JSON, using the 'body' variable
            console.log(body)
          //  js = body
            return js
        };
        console.log(body);
        return js;
    });

    console.log('httpcalled');

};


const f1 = function() {

    var resJson1 = 'json';
    var resCode1 = 'code';
    return { resJson : "fff", resCode : "fffs"};
}


//var a = getAPI();

//console.log('a=',a)

function getAPI1(callback) {

    const request = require('request');

    var token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJhc3Nob2xsZSIsImV4cCI6MTYyNDA3MzI0NiwiaXNzIjoid2VibGlua19hY2Nlc3MifQ.2x8Ke-DCgL_8mJQw5rahYt_htgAEYv7Aj0Dld3lHDbI';

    var options = {
      method: 'GET',
      hostname: '127.0.0.1',
      port:'8000',
      path: '/links/all',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer '+token,
      },
      json : true,

    };

    var js = 'haha';

    console.log('httpcall');

    var url = 'http://127.0.0.1:8000/links/all';

    request(url, options, function(error, response, body) {
        // субфукция получает респонз асинхронно
        // когда она получит и куда кидать инфу непонятно
        if (error) {
            callback(error)
        };

        //js = body
        if (!error ){//&& response.statusCode == 200) {
            callback(body)
        // do something with JSON, using the 'body' variable
            //console.log(body)
          //  js = body
           // return js
        };
        //console.log(body);
        // отправить js в определение getAPI
        //callback(js)
    });


    // may be a heavy db call or http request?
    // do not return any data, use callback mechanism
    //callback(js)
}

//var js = undefined;
getAPI1(function(mc /* js is passed using callback */) {
    console.log("mc=",mc); // a is 5
    return
})

//var result1 = getAPI();
//console.log('res=',mc);
//var json = result1.resJson

//console.log(json)
