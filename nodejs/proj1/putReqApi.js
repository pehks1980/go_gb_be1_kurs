
const http = require('http');

var token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJhc3Nob2xsZSIsImV4cCI6MTYyNDA3MzI0NiwiaXNzIjoid2VibGlua19hY2Nlc3MifQ.2x8Ke-DCgL_8mJQw5rahYt_htgAEYv7Aj0Dld3lHDbI';

var $shorturl = 'ndddew9945asdfg.drrrr';

const options = {
  method: 'PUT',
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

  let data = '';

  res.on('data', (chunk) => {
    data += chunk;
  });

  res.on('close', () => {
    console.log('Updated data');
    console.log(JSON.parse(data));
  });
});

const requestData = {
    'url': '1111www.mail.ruguaddda',
      'shorturl': 'asdfg.drrrr',
      'datetime': '0001-01-01T00:00:00Z',
      'active': 1,
      'redirs': 120
    };


request.write(JSON.stringify(requestData));

request.end();
