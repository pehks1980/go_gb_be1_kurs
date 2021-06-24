const http = require('http');

var token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJhc3Nob2xsZSIsImV4cCI6MTYyNDA3MzI0NiwiaXNzIjoid2VibGlua19hY2Nlc3MifQ.2x8Ke-DCgL_8mJQw5rahYt_htgAEYv7Aj0Dld3lHDbI';

const options = {
  method: 'POST',
  host: '127.0.0.1',
  port:'8000',
  path: '/links',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer '+token,
  }
};

const request = http.request(options, (res) => {
  if (res.statusCode !== 201) {
    console.error(`Did not get a Created from the server. Code: ${res.statusCode}`);
    res.resume();
    //return;
  }

  let data = '';

  res.on('data', (chunk) => {
    data += chunk;
  });

  res.on('close', () => {
    console.log(`Added new Link ${res.statusCode}`);
    console.log(JSON.parse(data));
  });
});

const requestData = {
    'url': '1111www.mail.ruguaddda',
      'shorturl': 'ndddew9945asdfg.drrrrrh1111',
      'datetime': '0001-01-01T00:00:00Z',
      'active': 1,
      'redirs': 1
    };


request.write(JSON.stringify(requestData));

request.end();
