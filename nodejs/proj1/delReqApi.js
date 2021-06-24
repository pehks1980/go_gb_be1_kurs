
const http = require('http');

var token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJhc3Nob2xsZSIsImV4cCI6MTYyNDA3MzI0NiwiaXNzIjoid2VibGlua19hY2Nlc3MifQ.2x8Ke-DCgL_8mJQw5rahYt_htgAEYv7Aj0Dld3lHDbI';

var $shorturl = 'new_asdfg.drrrr';

const options = {
  method: 'DELETE',
  host: '127.0.0.1',
  port:'8000',
  path: '/links/'+$shorturl,
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer '+token,
  }
};

const request = http.request(options, (res) => {
  if (res.statusCode !== 200) {
    console.error(`Did not get an OK from the server. Code: ${res.statusCode}`);
    res.resume();
    //return;
  }
    if (res.statusCode == 200){
        console.log('DELETE OK');
        return;
    }

  //console.log(res.statusCode)
  let data = '';

  res.on('data', (chunk) => {
    data += chunk;
  });


  res.on('close', () => {
      if (res.statusCode !== 200) {
        console.log('DELETE MESSAGE:');
        console.log(JSON.parse(data));
      }
  });
});

request.end();
