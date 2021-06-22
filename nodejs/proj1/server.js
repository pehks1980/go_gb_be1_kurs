
//// local data 'storage'
var token = '';  
var shorturl = '';
var puturl = '';
var putredirs = '';
var username = '';

// api address
const apiurl = 'http://127.0.0.1:8000';
// nodejs (this) server address:port
const nodejsurl = 'http://127.0.0.1:8090';
// must be the same as nodejsapi -):
const srvIP = '127.0.0.1';
const srvPort = '8090';

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
app.use(bodyParser.urlencoded({ extended: true }));

//// api cb funcs
//// get token for uid
function authAPI1(callback){

    var options = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      json : {'uid':username},
    };
    
    var url = apiurl + '/user/auth';
    
    console.log('http api auth call=',url,'username', options.json);
    
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
        // субфукция получает респонз асинхронно 
        if (error) {
            callback(error)
        };
        
        if (!error ){
            callback(body)
        };
    });
};
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
        'Authorization': 'Bearer '+token,
      },
      json : true,
    };

    var url = apiurl + '/shortstat/' + shorturl;

    console.log('http call get api link stat=',url);

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

//// get handlers
//// index page
app.get('/', function(req, res) {
    var mascots = [
        { name: 'Sammy', organization: "DigitalOcean", birth_year: 2012},
        { name: 'Tux', organization: "Linux", birth_year: 1996},
        { name: 'Moby Dock', organization: "Docker", birth_year: 2013}
    ];
    var tagline = "No programming concept is complete without a cute animal mascot.";

    res.render('page/index', {
        mascots: mascots,
        tagline: tagline,
        username: username,
    });
});

//// list items page
app.get('/list', function(req, res) {

    if (checktoken(token) == true) {

        getAPI1(function(mc /* mc api response is passed using callback */) {
            console.log("mc=", mc.data);
            //strip datetime to short format
            if (mc.data !== null){

                for (const x of mc.data) { 
                    let res = x.datetime.split(".");
                    x.datetime = res[0];
                    //console.log(x.datetime); 
                };
            ;}
            res.render('page/list', {mc: mc.data, username: username, api : apiurl, nodejs: nodejsurl});
        });

    } else {
        res.render('page/unathorized', {username: ''});    
    };

});
//// login form
app.get('/login', (req, res) => {
    res.render('page/login', {title : 'login to API', username: username,});
});
//// pressed edit link button form (from list)
app.get('/edit', (req, res) => {

    shorturl = req.query.shorturl
    console.log(`edit link ${shorturl}`);
    getAPI2(function(mc){
         // get json from api put it as instance //called 'get shortstat/{link} api'
        console.log("mc=", mc.data);
        res.render('page/edit', {title : 'Edit Link', username: username, instance: mc.data[0]});
    });

});
//// add new link form
app.get('/add', (req, res) => {
    res.render('page/add', {title : 'add new link', username: username,});
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
     postAPI1(function(res /* post to api new link*/) {
         console.log('api res',res);
     });

    console.log('click post accepted');
    // after add go to list page
    res.redirect(nodejsurl+'/list');
});

//// edit link post handler button click
app.post('/edit', (req, res) => {
    const click = {clickTime: new Date()};
    console.log(`put ${click.clickTime}`, req.body);
    // get back params from form
    shorturl = req.body.shorturl;
    puturl = req.body.url;
    putredirs = req.body.redirs;

    putAPI1(function(res /* update link*/) {
         console.log('api res',res);
     });
    console.log('click put accepted');
    // go to list page after link update
    res.redirect(nodejsurl+'/list');
});
    
//// login form post reply
app.post('/login', (req, res) => {
    // Login Code Here
    username = req.body.username;
    console.log(`auth ${req.body.username}`)
    authAPI1(function(res /* get pair of tokens api */) {
         console.log('api auth res',res.accessToken);
         // store access jwt token
         if (res.accessToken !== null){
             // if api is not responding we get undefined res!!!
             // so we store token only if api is ok and we get authorization token
             token = res.accessToken;
         }
    });
    console.log ('got authorization jwt token=', token)
    //redir to /
    res.redirect(nodejsurl+'/');
});
//// delete item post button click
app.post('/delete', (req, res) => {
    const click = {clickTime: new Date()};
    console.log(`delete ${click.clickTime}`, req.body.shorturl);
    //store key link for api call
    shorturl = req.body.shorturl;

    delAPI1(function(res /* api del call*/) {
         console.log('api res',res);
    });

    console.log('click del accepted');
    res.redirect(nodejsurl+'/list');
});

//// start node server.js
app.listen(srvPort,srvIP);
console.log('server node.js started http://'+ srvIP +':' + srvPort);