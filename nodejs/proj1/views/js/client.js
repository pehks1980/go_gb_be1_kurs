
window.onload = doStuff;

function doStuff() {
    
    console.log('Client-side code running');

    var buttons = document.querySelectorAll(".MyDelButton");
    var i = 0, length = buttons.length;
    for (i; i < length; i++) {
        if (document.addEventListener) {
            buttons[i].addEventListener("click", function() {
                // use keyword this to target clicked button
                console.log(`button with key ${this.id} was clicked`);
                var id = this.id
                // send key to node js
                fetch('/delete', {method: 'POST',
                                  headers: {
                                            'Accept': 'application/json',
                                            'Content-Type': 'application/json'
                                            },
                                  body: JSON.stringify({'shorturl': id})
                                 })
                    .then(function(response) {
                      if(response.ok) {
                        console.log('click was recorded');
                        window.location='http://127.0.0.1:8080/list';
                        return;
                      }
                      throw new Error('Request failed.');
                    })
                    .catch(function(error) {
                      console.log(error);
                    });
                
                });
            } else {
                buttons[i].attachEvent("onclick", function() {
                // use buttons[i] to target clicked button
                });
            };
    };
    
    //edit buttons
    var buttonsEd = document.querySelectorAll(".MyEdButton");
    var i = 0, length = buttonsEd.length;
    for (i; i < length; i++) {
        if (document.addEventListener) {
            buttonsEd[i].addEventListener("click", function() {
                // use keyword this to target clicked button
                console.log(`button edit with key ${this.id} was clicked`);
                var id = this.id
                // send key to node js
                window.location='http://127.0.0.1:8080/edit?shorturl='+id;
                
                });
            } else {
                buttonsEd[i].attachEvent("onclick", function() {
                // use buttons[i] to target clicked button
                });
            };
    };
    // add link button handler
    var addbutton = document.getElementById('addbutton');
    addbutton.addEventListener('click', function(e) {
      console.log('button add was clicked');

      window.location='http://127.0.0.1:8080/add';

    });

    setInterval(function() {
  fetch('/clicks', {method: 'GET'})
    .then(function(response) {
      if(response.ok) return response.json();
      throw new Error('Request failed.');
    })
    .then(function(data) {
      document.getElementById('counter').innerHTML = `Button was clicked ${data.times} times`;
    })
    .catch(function(error) {
      console.log(error);
    });
}, 1000);

};