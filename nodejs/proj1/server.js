const request = require('request');

var token = '';  

var shorturl = '';

var puturl = '';
var putredirs = '';

var username = '';

//// check if we have logged in
function checktoken(token){
    if (token !== '')  {
        return true
    }
    return false
}

//// express setup
//// load the things we need
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
    
    console.log('http api auth call=',username);
    
    var url = 'http://127.0.0.1:8000/user/auth';
    
    request(url, options, function(error, response, body) {
        // субфукция получает респонз асинхронно 
        // когда она получит и куда кидать инфу непонятно
        if (error) {
            callback(error)
        };
        
        if (!error ){//&& response.statusCode == 200) {
            callback(body)
        };
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
        
    console.log('http api del call=',shorturl);
    
    var url = 'http://127.0.0.1:8000/links/'+shorturl;
    
    request(url, options, function(error, response, body) {
        // субфукция получает респонз асинхронно 
        // когда она получит и куда кидать инфу непонятно
        if (error) {
            callback(error)
        };
        
        if (!error ){//&& response.statusCode == 200) {
            callback(body)
        };
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
        
    console.log('http api put call=',shorturl);
    
    var url = 'http://127.0.0.1:8000/links/'+shorturl;
    
    request(url, options, function(error, response, body) {
        // субфукция получает респонз асинхронно 
        // когда она получит и куда кидать инфу непонятно
        if (error) {
            callback(error)
        };
        
        if (!error ){//&& response.statusCode == 200) {
            callback(body)
        };
    });
}
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
        
    console.log('http api post (create) call=');
    
    var url = 'http://127.0.0.1:8000/links';
    
    request(url, options, function(error, response, body) {
        // субфукция получает респонз асинхронно 
        // когда она получит и куда кидать инфу непонятно
        if (error) {
            callback(error)
        };
        
        if (!error ){//&& response.statusCode == 200) {
            callback(body)
        };
    });
}
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
    
    var js = 'haha';
    
    console.log('httpcall token=',token);
    
    var url = 'http://127.0.0.1:8000/links/all';
    
    request(url, options, function(error, response, body) {
        // субфукция получает респонз асинхронно 
        // когда она получит и куда кидать инфу непонятно
        if (error) {
            callback(error)
        };
        
        if (!error ){//&& response.statusCode == 200) {
            callback(body)

        };

    });

}
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
    
    console.log('httpcall token=',token);
    
    var url = 'http://127.0.0.1:8000/shortstat/'+shorturl;
    
    request(url, options, function(error, response, body) {
        // субфукция получает респонз асинхронно 
        // когда она получит и куда кидать инфу непонятно
        if (error) {
            callback(error)
        };
        
        if (!error ){//&& response.statusCode == 200) {
            callback(body)

        };

    });

}


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
        
        getAPI1(function(mc /* js is passed using callback */) {
            var quotation = [
            {
                  "dimension" : 0,
                  "currency": "RMB",
                  "quantity": "100",
                  "price": "3",
                  "factory": "rx"},

            ];
            console.log("mc=", mc.data);
            //console.log("quot=",quotation);
            if (mc.data !== null){
                
                for (const x of mc.data) { 
                    let res = x.datetime.split(".");
                    x.datetime = res[0];
                    console.log(x.datetime); 
                };
            ;}
            res.render('page/list', {quotation: quotation, mc: mc.data, username: username,});
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
         // get json from req.body. put it as instance
        console.log("mc=", mc.data);

        res.render('page/edit', {title : 'Edit Link', username: username, instance: mc.data[0]});
    });        
  
});

app.get('/add', (req, res) => {
  res.render('page/add', {title : 'add new link', username: username,});
});

app.post('/add', (req, res) => {
  const click = {clickTime: new Date()};

  //console.log(req.body);  
  console.log(`create ${click.clickTime}`, req.body);
  //console.log(db);
    
  //must be unique link
  shorturl = req.body.shorturl;
  let r = Math.random().toString(36).substring(7);
  //add unique part
  shorturl = shorturl + r;
  puturl = req.body.url;
  putredirs = req.body.redirs;
  
  console.log(shorturl,puturl,putredirs)
    // todo shorturl param
  postAPI1(function(res /* js is passed using callback */) {
       console.log('api res',res);     
   });
    
  console.log('click post accepted');
  res.redirect('http://127.0.0.1:8080/list');
  ///res.sendStatus(200);
});

app.post('/edit', (req, res) => {
  const click = {clickTime: new Date()};

  //console.log(req.body);  
  console.log(`put ${click.clickTime}`, req.body);
  //console.log(db);
  shorturl = req.body.shorturl;
  puturl = req.body.url;
  putredirs = req.body.redirs;
  
  console.log(shorturl,puturl,putredirs)
    // todo shorturl param
  putAPI1(function(res /* js is passed using callback */) {
       console.log('api res',res);     
   });
    
  console.log('click put accepted');
  res.redirect('http://127.0.0.1:8080/list');
  ///res.sendStatus(200);
});
    
//// login form post reply
app.post('/login', (req, res) => {
  // Login Code Here
    username = req.body.username;
    console.log(`auth ${req.body.username}`) 
    authAPI1(function(res /* js is passed using callback */) {
       console.log('api auth res',res.accessToken);
       token = res.accessToken
    });
    //redir to /
    res.redirect('http://127.0.0.1:8080/');
  //res.sendStatus(200);
});
//// delete item
app.post('/delete', (req, res) => {
  const click = {clickTime: new Date()};

  //console.log(req.body);  
  console.log(`delete ${click.clickTime}`, req.body);
  //console.log(db);
   shorturl = req.body.shorturl;
    console.log(shorturl)
    // todo shorturl param
   delAPI1(function(res /* js is passed using callback */) {
       console.log('api res',res);
       
   });
    
  console.log('click del accepted');
  res.redirect('http://127.0.0.1:8080/list');
  ///res.sendStatus(200);
});

// etc add a document to the DB collection recording the click event
app.post('/clicked', (req, res) => {
  const click = {clickTime: new Date()};
  console.log(click);
  //console.log(db);
    
    console.log('click added to db');
    res.sendStatus(201);
});
// get the click data from the database
app.get('/clicks', (req, res) => {
    result = {"times": 120};
    res.send(result);
});

//// start node server.js
app.listen(8080);
console.log('8080 is the magic port');