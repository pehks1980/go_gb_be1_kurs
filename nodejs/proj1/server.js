//// local data 'storage'
var token = '';
var tokenPayload = ''
var user = {
        name : '',
        role : '',
        balance: ''
    };
// select api address
const apiurl = 'http://127.0.0.1:8000'; //local
//const apiurl = 'https://web-link19801.herokuapp.com'; // heroku
// nodejs (this) server address:port
const nodejsurl = 'http://127.0.0.1:8090';
// must be the same as nodejsurl -):
const srvIP = '127.0.0.1';
const srvPort = '8090';
function jwtdecode(token) {
    let decoded = jwt_decode(token);
    console.log('token payload= ',decoded);
    return decoded;
}
//// check if we have logged in
function checktoken(token){
    if (token === undefined){
        console.log('token is undefined please check if api is running at all!!!')
        return false;
    };

    if (token !== '')  {
        return true;
    };

    return false;
};
//// express setup
//// load the things we need
const request = require('request');
var ejs = require('ejs');
var express = require('express');
var app = express();
// set the view engine to ejs
app.set('view engine', 'ejs');
// set static files catalog relation './views' (fs side) : '/' (web side)
app.use(express.static('./views'));
// allow latest json app use
app.use(express.json());
//post form data parser
const bodyParser = require('body-parser');
const jwt_decode = require('jwt-decode');
app.use(bodyParser.urlencoded({ extended: true }));
//// api cb funcs
//// get token for uid
function getUserAPI(callback, uid){

    var options = {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer '+token,
        },
        json : true,
    };

    var url = apiurl + '/user/'

    if (uid !== undefined ){
        url = apiurl + '/user/'+uid
    }

    console.log('http api user get call=',url);

    request(url, options, function(error, response, body) {
        // субфукция получает респонз асинхронно
            callback(response,body)
    });
}
//// get token for uid
function getAllUsersAPI(callback){

    var options = {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer '+token,
        },
        json : true,
    };

    var url = apiurl + '/users/all';

    console.log('http api user getall call=',url);

    request(url, options, function(error, response, body) {
            callback(response, body)
    });
};
//// get token for uid
function regAPI(callback, username, passwd, email){

    var options = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      json : {'name':username,'passwd':passwd,'email':email},
    };

    var url = apiurl + '/user/register';

    console.log('http api auth call=',url,'username=', options.json);

    request(url, options, function(error, response, body) {
        callback(response,body)
    });
}

function authAPI1(callback, username, password){

    var options = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        json : {'name':username,'passwd':password},
    };

    var url = apiurl + '/user/auth';

    console.log('http api auth call=',url,'username', options.json);

    request(url, options, function(error, response, body) {
        callback(response,body)
    });
}
//// delete user
function delUserAPI(callback,uid){

    var options = {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer '+token,
        },
        json : true,
    };

    var url = apiurl + '/user/' + uid;

    console.log('http api user del call=',url);

    request(url, options, function(error, response, body) {
        // results gets to callback
        callback(response,body)
    });
}
/// put user
function putUserAPI(callback,user,uid){

    var options = {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer '+token,
        },
        json : user,
    };

    var url = apiurl + '/user/' + uid;

    console.log('http api user put call=',url,'data=', options.json);

    request(url, options, function(error, response, body) {
        // субфукция получает респонз асинхронно
        callback(response,body)
    });
}
//// delete item
function delLinkItemAPI(callback, shorturl){

    var options = {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer '+token,
      },
      json : true,

    };

    var url = apiurl + '/links/' + shorturl;

    console.log('http api del call=',url);

    request(url, options, function(error, response, body) {
        // results gets to callback
            callback(response,body)
    });
}
/// put item
function putAPI1(callback){

    var redirsint = parseFloat(putredirs);
    var options = {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer '+token,
      },
      json : {
          'url': puturl,
          'shorturl' : shorturl,
          'redirs' : redirsint,
          'active' : 1,
      },

    };

    var url = apiurl + '/links/' + shorturl;

    console.log('http api put call=',url,'data=', options.json);

    request(url, options, function(error, response, body) {
        // субфукция получает респонз асинхронно
        if (error) {
            callback(error)
        };

        if (!error ){
            callback(body)
        };
    });
}
/// post item
function postAPI1(callback, item){

    var options = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer '+token,
      },
      json : {
          'url': item.puturl,
          'shorturl' : item.shorturl,
          'redirs' : parseFloat(item.putredirs),
          'active' : 1,
      },
    };

    var url = apiurl + '/links';

    console.log('http api post (create) call=',url,'data=', options.json);

    request(url, options, function(error, response, body) {
        callback(response,body)
    });
}
//// get list of items
function getItemsListAPI(callback) {

    var options = {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer '+token,
      },
      json : true,
    };

    var url = apiurl + '/links/all';

    console.log('http call get api list=',url,'token=',token);

    request(url, options, function(error, response, body) {
        callback(response, body)
    });

}
//// get link item
function getLinkItemAPI(callback, shorturl) {
    var options = {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token,
        },
        json: true,
    };

    var url = apiurl + '/shortstat/' + shorturl;

    console.log('http call get api link stat=', url);

    request(url, options, function (error, response, body) {
        callback(response,body)
    });
}

function getShortOpenAPI(callback) {
    var options = {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token,
        },
        json: true,
    };

    var url = apiurl + '/shortopen/' + shorturl;

    console.log('http call get api link open=', url);

    request(url, options, function (error, response, body) {
        callback(response,body)
    });
}
//// get handlers
//// index page
app.get('/', function(req, res) {
    var mascots = [
        {name: 'Sammy', organization: "DigitalOcean", birth_year: 2012},
        {name: 'Tux', organization: "Linux", birth_year: 1996},
        {name: 'Moby Dock', organization: "Docker", birth_year: 2013}
    ];
    var tagline = "No programming concept is complete without a cute animal mascot.";

    if (token == '') {
        res.render('page/index', {
            mascots: mascots,
            tagline: tagline,
            user: user,
        });
        return;
    }

    getUserAPI(
    function (resp, body) {

            if (resp == undefined) {
                res.render('page/index', {
                    mascots: mascots,
                    tagline: tagline,
                    user: user,
                });
                return;
            }

            if (resp.statusCode != 200) {
                res.render('page/index', {
                    mascots: mascots,
                    tagline: tagline,
                    user: user,
                });
                return;
            }

            user.name = body.name
            user.role = body.role
            user.balance = body.balance
            console.log('api getUserAPI res ', body);

            res.render('page/index', {
                mascots: mascots,
                tagline: tagline,
                user: user,
            });

        });


});
//// admin page
app.get('/admin', function(req, res) {
    console.log('api haha\n$$$$$$$$$$$$$$$$$$$$$$$$$$\n')
    if (!checktoken(token)) {
        // no token
        res.render('page/unathorized', {user: null});
        return;
    }

    console.log("user=",user)

    getAllUsersAPI(
        function(resp, mc) {

            if ( resp == undefined ){
                res.render('page/unathorized', {user: user});
                return;
            }

            console.log("mc=", mc);

            if (resp.statusCode !== 200) {
                // error
                res.render('page/unathorized', {user: user});
                return;
            }

            res.render('page/admin', {mc: mc.data,
                user: user,
                api: apiurl,
                nodejs: nodejsurl
            });

    });

});
//// list items page
app.get('/list', function(req, res) {

        if (!checktoken(token)) {
            res.render('page/unathorized', {user: user});
            return;
        }

        getItemsListAPI(
            function(resp, mc) {

                if (resp == undefined){
                    res.render('page/unathorized', {user: user});
                    return;
                }

                //console.log("mc=", mc);

                if (resp.statusCode !== 200) {
                    // error
                    res.render('page/unathorized', {user: user});
                    return;
                }

                if ( mc.data === undefined || mc.data === null) {
                    // empty list
                    mc.data = null;
                    res.render('page/list', {
                        mc: mc.data,
                        user: user,
                        api: apiurl,
                        nodejs: nodejsurl
                    });
                    return;
                }


                // not empty list
                // strip datetime to short format
                for (const x of mc.data) {
                    let res = x.datetime.split(".");
                    let res1 = res[0].split("T");
                    x.datetime = res1[1] + ' ' + res1[0];
                    //console.log(x.datetime);
                }

                res.render('page/list', {mc: mc.data,
                    user: user,
                    api: apiurl,
                    nodejs: nodejsurl
                });
            });

});
//// list update
app.get('/listupd', (req, res) => {

    if (!checktoken(token)) {
        res.render('page/unathorized', {user: user});
        return;
    }

    getItemsListAPI(
        function(resp, mc) {

            if (resp == undefined){
                res.render('page/unathorized', {user: user});
                return;
            }

            //console.log("mc=", mc);

            if (resp.statusCode !== 200) {
                // error
                res.render('page/unathorized', {user: user});
                return;
            }

            if (mc.data === null) {
                // empty list
                result = {"table": ""};
                res.send(result);
                return;
            }

            // not empty list
            // strip datetime to short format
            for (const x of mc.data) {
                let res = x.datetime.split(".");
                let res1 = res[0].split("T");
                x.datetime = res1[1] + ' ' + res1[0];
                console.log(x.datetime);
            }

            ejs.renderFile('views/part/table.ejs', {mc: mc.data, api: apiurl}, {}, function (err, str) {
                // str => Rendered HTML string
                //console.log(str,err)
                result = {"table": str};
                res.send(result);
            });

        });

});
//// login form
app.get('/login', (req, res) => {
    res.render('page/login', {title : 'login to API', user: user,});
});
//// user register  form
app.get('/register', (req, res) => {
    res.render('page/register', {title : 'Register new user', user: user,});
});
//// user delete item post button click
app.post('/deluser', (req, res) => {

    if (!checktoken(token)) {
        res.render('page/unathorized', {user: user});
        return;
    }

    const click = {clickTime: new Date()};
    console.log(`delete user ${click.clickTime}`, req.body.uid);
    //store key link for api call
    let uid = req.body.uid;

    delUserAPI(function(resp,body /* api del call*/) {
        console.log('api user delete resp=', resp.statusCode, 'body=',body);

        if (resp == undefined){
            res.render('page/unathorized', {user: user});
            return;
        }

        if (resp.statusCode !== 200 ) {
            res.render('page/unathorized', {user: user});
            return;
        }

        res.redirect('/admin')

    },uid);

});
//// pressed edit link button form (from list)
app.get('/upduser', (req, res) => {

    if (!checktoken(token)) {
        res.render('page/unathorized', {user: user});
        return;
    }

    let uid = req.query.uid
    console.log(`edit user ${uid}`);

    getUserAPI(
        function (resp,mc) {

            if (resp == undefined) {
                res.render('page/unathorized', {user: user});
                return;
            }

            if (resp.statusCode != 200) {
                res.render('page/unathorized', {user: user});
                return;
            }

            console.log("mc=", mc);
            if (mc !== undefined) {
                // add uid to form to pass it back in post put user api
                mc.uid = uid
                res.render('page/useredit', {title: 'Edit User', user: user, instance: mc});
                return;
            }

            res.redirect('/');

            },uid);
});
//// edit user post handler button click
app.post('/upduser', (req, res) => {

    if (!checktoken(token)) {
        res.render('page/unathorized', {user: user});
        return;
    }

    const click = {clickTime: new Date()};
    console.log(`put ${click.clickTime}`, req.body);

    // get back params from form
    let user = {
        name:    req.body.name,
        email:   req.body.email,
        role:    req.body.role,
        balance: req.body.balance,
    }
    let uid = req.body.uid

    putUserAPI(
        function(resp,body /* update user*/) {

        if (resp == undefined) {
            res.render('page/unathorized', {user: user});
            return;
        }

        if (resp.statusCode != 200) {
            res.render('page/unathorized', {user: user});
            return;
        }

        console.log('api res',resp.statusCode, body);

        res.redirect(nodejsurl+'/admin');
    },user,uid);

});
//// pressed check link button form (from list)
app.get('/check', (req, res) => {

    if (!checktoken(token)) {
        res.render('page/unathorized', {user: user});
        return;
    }

    shorturl = req.query.shorturl
    console.log(`check open link ${shorturl}`);

    getShortOpenAPI(
        function (resp, mc) {

            if (resp == undefined) {
                res.render('page/unathorized', {user: user});
                return;
            }

            if (resp.statusCode != 200) {
                res.render('page/unathorized', {user: user});
                return;
            }
            // get json from api put it as instance //called 'get shortstat/{link} api'
            console.log("mc=", mc);
            // no edit data - silently go to /
            if (mc.url == undefined) {
                res.redirect('/');
            }

            res.redirect(mc.url);
    });

});
//// pressed edit link button form (from list)
app.get('/edit', (req, res) => {

    if (!checktoken(token)) {
        res.render('page/unathorized', {user: user});
        return;
    }

    let shorturl = req.query.shorturl
    console.log(`edit link ${shorturl}`);

    getLinkItemAPI(function (resp,mc) {

        if (resp === undefined) {
            res.render('page/unathorized', {user: user});
            return;
        }

        if (resp.statusCode !== 200) {
            res.render('page/unathorized', {user: user});
            return;
        }

        // get json from api put it as instance //called 'get shortstat/{link} api'
        console.log("mc=", mc.data);

        if (mc.data == undefined) {
            res.redirect('/');
            return;
        }

        res.render('page/edit', {title: 'Edit Link', user: user, instance: mc.data[0]});
    }, shorturl);

});
//// add new link form
app.get('/add', (req, res) => {
    res.render('page/add', {title : 'add new link', user: user,});
});
//// POST handlers
//// add new link post form
app.post('/add', (req, res) => {

    if (!checktoken(token)) {
        res.render('page/unathorized', {user: user});
        return;
    }

    const click = {clickTime: new Date()};
    console.log(`create ${click.clickTime}`, req.body);
    //must be unique link = text + random part
    let r = Math.random().toString(36).substring(7);
    let item = {
        shorturl: req.body.shorturl + r,
        puturl: req.body.url,
        putredirs : req.body.redirs
    }

    //console.log(shorturl,puturl,putredirs)
     postAPI1(function(resp,body )/* post to api new link*/ {

         if (resp == undefined){
             res.render('page/unathorized', {user: user});
             return;
         }

         if (resp.statusCode !== 201 ) {
             res.render('page/unathorized', {user: user});
             return;
         }

         console.log('api res added=',body, 'status', resp.statusCode);
         // after add go to list page
         getUserAPI(
             function (resp,body) {
                 user.name = body.name
                 user.role = body.role
                 user.balance = body.balance
                 console.log('api getUserAPI res',resp.statusCode);
                 console.log('click add accepted');
                 res.redirect('/list');
             });

     }, item);

});
//// edit link post handler button click
app.post('/edit', (req, res) => {

    if (!checktoken(token)) {
        res.render('page/unathorized', {user: user});
        return;
    }

    const click = {clickTime: new Date()};
    console.log(`put ${click.clickTime}`, req.body);
    // get back params from form
    shorturl = req.body.shorturl;
    puturl = req.body.url;
    putredirs = req.body.redirs;

    putAPI1(function(res1 /* update link*/) {
         console.log('api res',res1);
         console.log('click put accepted');// go to list page after link update
         res.redirect('/list');
    });

});
//// login form post reply
app.post('/login', (req, res) => {
    // Login Code Here
    let username = req.body.username;
    let password = req.body.password;
    console.log(`auth ${req.body.username}`)

    authAPI1(function(resp,body) {

        //console.log('api auth res',resp, 'body=',body);
        if ( resp === undefined){
            console.log('please check if api is running at all!!')
            res.render('page/unathorized', {user: user});
            return;
        }

        if (resp.statusCode !== 200) {
            //not successful
            res.render('page/unathorized', {user: user});
            return;
        }

        if (body.accessToken == null) {
            res.render('page/unathorized', {user: user});
            return;
        }

        // store access jwt token
        // if api is not responding we get undefined res!!!
        // so we store token only if api is ok and we get authorization token
        token = body.accessToken;
        tokenPayload = jwtdecode(token)
        let uid = tokenPayload.uid
        console.log('api auth res', body.accessToken,'uid=',uid);

        getUserAPI(
            function (res1) {
                user.name = res1.name
                user.role = res1.role
                user.balance = res1.balance
                //console.log('api getUserAPI res', res1);
                console.log('redirect to ', nodejsurl + '/', 'user=', user);
                res.redirect(nodejsurl + '/');
            });

    }, username, password);

});
//// reg form new userpost reply
app.post('/register', (req, res) => {
    // Register Code Here
    let username = req.body.username;
    let email = req.body.email;
    let password = req.body.password;
    ///
    console.log(`auth ${req.body.username}`)

    regAPI(function(res1) {
        console.log('api register res',res1);
        //redir to /
        res.redirect(nodejsurl+'/login');
    }, username,password,email);

});
//// delete item post button click
app.post('/delete', (req, res) => {

    if (!checktoken(token)) {
        res.render('page/unathorized', {user: user});
        return;
    }

    const click = {clickTime: new Date()};
    console.log(`delete ${click.clickTime}`, req.body.shorturl);
    //store key link for api call
    let shorturl = req.body.shorturl;

    delLinkItemAPI(function (resp,body) /* api del call*/{

        if (resp == undefined) {
            res.render('page/unathorized', {user: user});
            return;
        }

        if (resp.statusCode != 200) {
            console.log('statusCode ', resp.statusCode);
            res.render('page/unathorized', {user: user});
            return;
        }

        res.redirect('/list')

    },shorturl);

});

// запускаем apm agent
var apm = require('elastic-apm-node').start({
    serviceName: 'weblinknodeserver',
    serverUrl: 'http://127.0.0.1:8200'
})
//// start node server.js
app.listen(srvPort,srvIP);
console.log('server node.js started http://'+ srvIP +':' + srvPort, ' (API URL:', apiurl, ')');
