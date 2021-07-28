
//// local data 'storage'
var token = '';
var tokenPayload = ''
var shorturl = '';
var puturl = '';
var putredirs = '';
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
        // субфукция получает респонз асинхронно
        // body - уже json
        if (error) {
            callback(error)
        };

        if (!error ){
            callback(body)
        };
    });
};

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
        // субфукция получает респонз асинхронно
        // body - уже json
        //console.log(response,body)
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
function delAPI1(callback){

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
            callback(error,response,body)
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
};
/// post item
function postAPI1(callback){
    var redirsint = parseFloat(putredirs);
    var options = {
      method: 'POST',
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

    var url = apiurl + '/links';

    console.log('http api post (create) call=',url,'data=', options.json);

    request(url, options, function(error, response, body) {
        // субфукция получает респонз асинхронно
        if (error) {
            callback(error)
        };

        if (!error ){
            callback(body)
        };
    });
};
//// get list of items
function getAPI1(callback) {

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
        // субфукция получает респонз асинхронно
        if (error) {
            callback(error)
        };

        if (!error ){
            callback(body)
        };

    });

};
//// get link item
function getAPI2(callback) {
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
        // субфукция получает респонз асинхронно
        if (error) {
            callback(error)
        };

        if (!error) {
            callback(body)
        };

    });

};

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
        // субфукция получает респонз асинхронно
        if (error) {
            callback(error)
        };

        if (!error) {
            callback(body)
        };

    });

};


//// get handlers
//// index page
app.get('/', function(req, res) {
    var mascots = [
        {name: 'Sammy', organization: "DigitalOcean", birth_year: 2012},
        {name: 'Tux', organization: "Linux", birth_year: 1996},
        {name: 'Moby Dock', organization: "Docker", birth_year: 2013}
    ];
    var tagline = "No programming concept is complete without a cute animal mascot.";

    if (token !== '') {
        getUserAPI(
            function (resp,body) {

                if (resp.statusCode == 200) {
                    user.name = body.name
                    user.role = body.role
                    user.balance = body.balance
                    console.log('api getUserAPI res haha', body);

                    console.log('api getUserAPI res haha', body);
                    //res.redirect(nodejsurl+'/list');
                    res.render('page/index', {
                        mascots: mascots,
                        tagline: tagline,
                        user: user,
                    });
                }

            });
    } else {
        res.render('page/index', {
            mascots: mascots,
            tagline: tagline,
            user: user,
        });
    }
});

app.get('/admin', function(req, res) {
    console.log('api haha\n$$$$$$$$$$$$$$$$$$$$$$$$$$\n')
    if (checktoken(token) == true) {

        console.log("user=",user)

        getAllUsersAPI(function(resp, mc) {
            console.log("mc=", mc);

            if (resp.statusCode !== 200) {
                // error
                res.render('page/unathorized', {user: user});
            } else if (mc.data === null) {
                // empty list
                res.render('page/admin', {mc: mc.data,
                    user: user,
                    api: apiurl,
                    nodejs: nodejsurl
                });
            } else {
                // not empty list
                // strip datetime to short format
                /*
                for (const x of mc.data) {
                    let res = x.datetime.split(".");
                    let res1 = res[0].split("T");
                    x.datetime = res1[1] + ' ' + res1[0];
                    console.log(x.datetime);
                }

                 */
                res.render('page/admin', {mc: mc.data,
                    user: user,
                    api: apiurl,
                    nodejs: nodejsurl
                });
            }

        });

    } else {
        // no token
        res.render('page/unathorized', {user: null});
    };
});
//// list items page
app.get('/list', function(req, res) {

    if (checktoken(token) == true) {

        console.log("user=",user)

        getAPI1(function(mc /* mc api response is passed using callback */) {
            console.log("mc=", mc);


            if ('errors' in mc || 'Error' in mc) {
                // error
                res.render('page/unathorized', {user: user});
            } else if (mc.data === null) {
                // empty list
                res.render('page/list', {mc: mc.data,
                    user: user,
                    api: apiurl,
                    nodejs: nodejsurl
                });
            } else {
                // not empty list
                // strip datetime to short format
                for (const x of mc.data) {
                    let res = x.datetime.split(".");
                    let res1 = res[0].split("T");
                    x.datetime = res1[1] + ' ' + res1[0];
                    console.log(x.datetime);
                }
                res.render('page/list', {mc: mc.data,
                    user: user,
                    api: apiurl,
                    nodejs: nodejsurl
                });
            };

        });

    } else {
        // no token
        res.render('page/unathorized', {user: null});
    };

});

app.get('/listupd', (req, res) => {

    if (checktoken(token) == true) {

        getAPI1(function (mc /* mc api response is passed using callback */) {
            console.log("mc=", mc);
            if ('errors'  in mc || 'Error' in mc) {
                // no token
                res.render('page/unathorized', {user: user,});
            }

            if (mc.data !== null) {
                // not empty list mc.data
                // strip datetime to short format
                for (const x of mc.data) {
                    let res = x.datetime.split(".");
                    let res1 = res[0].split("T");
                    x.datetime = res1[1] + ' ' + res1[0];
                    console.log(x.datetime);
                };

                ejs.renderFile('views/part/table.ejs', {mc: mc.data, api: apiurl}, {}, function (err, str) {
                    // str => Rendered HTML string
                    //console.log(str,err)
                    result = {"table": str};
                    res.send(result);
                });

            } else {
                // empty list
                //ejs.renderFile('views/part/table.ejs', {mc: mc.data, api: apiurl}, {}, function (err, str) {
                    // str => Rendered HTML string
                    //console.log(str,err)
                    result = {"table": ""};
                    res.send(result);
                //});

            }





        });
    } else {
        res.render('page/unathorized', {user: user,});
    }
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
    const click = {clickTime: new Date()};
    console.log(`delete user ${click.clickTime}`, req.body.uid);
    //store key link for api call
    let uid = req.body.uid;

    delUserAPI(function(resp,body /* api del call*/) {
        console.log('api user delete resp=', resp.statusCode, 'body=',body);
        if (resp.statusCode !== 200 ){
            console.log('click user was deleted / ');
            //we cant redir to POST
            //res.redirect(nodejsurl+'/');
        } else {
            //res.setHeader('Content-Type', 'application/json');
            console.log('click redir /list');
            ///res.redirect('/');
        }

    },uid);
    console.log('end del user post ')
});
//// pressed edit link button form (from list)
app.get('/upduser', (req, res) => {

    let uid = req.query.uid
    console.log(`edit user ${uid}`);
    if (checktoken(token) == true) {

        getUserAPI(
            function (resp,mc) {
                if (resp.statusCode === 200) {
                    console.log("mc=", mc);
                    if (mc !== undefined) {
                        // add uid to form to pass it back in post put user api
                        mc.uid = uid
                        res.render('page/useredit', {title: 'Edit User', user: user, instance: mc});
                    } else {
                        // no edit data - silently go to /
                        res.redirect('/');
                    }
                }
            },uid);
    } else {
        res.render('page/unathorized', {user: user});
    }

});
//// edit user post handler button click
app.post('/upduser', (req, res) => {
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

    putUserAPI(function(resp,body /* update user*/) {
        console.log('api res',resp.statusCode, body);
        console.log('click put accepted');// go to list page after link update
        res.redirect(nodejsurl+'/admin');
    },user,uid);

});



//// pressed check link button form (from list)
app.get('/check', (req, res) => {

    shorturl = req.query.shorturl
    console.log(`edit link ${shorturl}`);
    if (checktoken(token) == true) {
        getShortOpenAPI(function (mc) {
            // get json from api put it as instance //called 'get shortstat/{link} api'
            console.log("mc=", mc);
            if (mc.url !== undefined) {
                res.redirect(mc.url)
            } else {
                // no edit data - silently go to /
                res.redirect('/');
            }
        });
    }else{
        res.render('page/unathorized', {user: user});
    }

});

//// pressed edit link button form (from list)
app.get('/edit', (req, res) => {

    shorturl = req.query.shorturl
    console.log(`edit link ${shorturl}`);
    if (checktoken(token) == true) {
        getAPI2(function (mc) {
            // get json from api put it as instance //called 'get shortstat/{link} api'
            console.log("mc=", mc.data);
            if (mc.data !== undefined) {
                res.render('page/edit', {title: 'Edit Link', user: user, instance: mc.data[0]});
            } else {
                // no edit data - silently go to /
                res.redirect('/');
            }
        });
    }else{
        res.render('page/unathorized', {user: user});
    }

});
//// add new link form
app.get('/add', (req, res) => {
    res.render('page/add', {title : 'add new link', user: user,});
});

//// POST handlers
//// add new link post form
app.post('/add', (req, res) => {
    const click = {clickTime: new Date()};
    console.log(`create ${click.clickTime}`, req.body);
    //must be unique link = text + random part
    shorturl = req.body.shorturl;
    let r = Math.random().toString(36).substring(7);
    //add unique part
    shorturl = shorturl + r;
    // the rest params
    puturl = req.body.url;
    putredirs = req.body.redirs;

    //console.log(shorturl,puturl,putredirs)
     postAPI1(function(res1 /* post to api new link*/) {
         console.log('api res added=',res1.statusCode);
         // after add go to list page
         getUserAPI(
             function (resp,body) {
                 user.name = body.name
                 user.role = body.role
                 user.balance = body.balance
                 console.log('api getUserAPI res',resp.statusCode);
                 console.log('click add accepted');
                 return res.redirect(nodejsurl+'/list');
             });

     });


});

//// edit link post handler button click
app.post('/edit', (req, res) => {
    const click = {clickTime: new Date()};
    console.log(`put ${click.clickTime}`, req.body);
    // get back params from form
    shorturl = req.body.shorturl;
    puturl = req.body.url;
    putredirs = req.body.redirs;

    putAPI1(function(res1 /* update link*/) {
         console.log('api res',res1);
         console.log('click put accepted');// go to list page after link update
         res.redirect(nodejsurl+'/list');
     });


});

//// login form post reply
app.post('/login', (req, res) => {
    // Login Code Here
    let username = req.body.username;
    let password = req.body.password;
    console.log(`auth ${req.body.username}`)
    authAPI1(function(resp,body) {
        console.log('api auth res',resp.statusCode, 'body=',body);
        if (resp.statusCode === 200){
            console.log('api auth res', body.accessToken);
            // store access jwt token
            if (body.accessToken !== null) {
                // if api is not responding we get undefined res!!!
                // so we store token only if api is ok and we get authorization token
                token = body.accessToken;
                tokenPayload = jwtdecode(token)
                let uid = tokenPayload.uid
                getUserAPI(
                    function (res1) {
                        user.name = res1.name
                        user.role = res1.role
                        user.balance = res1.balance
                        console.log('api getUserAPI res', res1);
                        console.log('redirect to ', nodejsurl + '/', 'user=', user);
                        res.redirect(nodejsurl + '/');
                    });


                /*
                getAllUsersAPI(
                    function (res) {
                        console.log('api getUsersAPI res',res);

                    });

                 */
            }
        }
        else {
            //not sucessful
            res.render('page/unathorized', {user: user});
        }
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
    const click = {clickTime: new Date()};
    console.log(`delete ${click.clickTime}`, req.body.shorturl);
    //store key link for api call
    shorturl = req.body.shorturl;

    delAPI1(function(err,res1,body /* api del call*/) {
        console.log('api delete err=',err, 'res1=', res1.statusCode, 'body=',body);
        console.log('click del accepted');
        if (res1.statusCode !== 200 ){
            console.log('click redir / ');
            //we cant redir to POST
            //res.redirect(nodejsurl+'/');
        } else {
            res.setHeader('Content-Type', 'application/json');
            console.log('click redir /list');
            ///res.redirect('/');
        }

    });


});

//// start node server.js
app.listen(srvPort,srvIP);
console.log('server node.js started http://'+ srvIP +':' + srvPort, ' (API URL:', apiurl, ')');
