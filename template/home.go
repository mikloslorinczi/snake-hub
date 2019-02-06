package template

import "html/template"

// Home is the index page
var Home = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
<title>Snake-hub</title>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<meta http-equiv="X-UA-Compatible" content="ie=edge">
<link href="data:image/x-icon;base64,AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAgv+AAARcxwA9Wv8ATGb8AAKoWgACJhUAByCwACE3tQCaqPUAOE7PAAgs/AAAAAAAAAAAAAAAAAAAAAAAAFEVAAAAAAAAURUAAAdwAABRFVUAqocAAFERFQAkqAAGUREVAJs6AAZVVVUACSAAAAZmAGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABVVQAAAAAAAFEVAAAAAAAAURUAAAAAAABRFQAAAAAAAFEVAAAAAAAAURUACAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAQAA" rel="icon" type="image/x-icon" />
</head>
<body>
<div class="container"></div>
<script type="text/javascript">
'use strict';

const uId = Math.random().toString(36).substring(2) + (new Date()).getTime().toString(36);

var parser = document.createElement('a');
parser.href = window.location.href;
if (parser.protocol == 'http:') {
    parser.protocol = "ws:"
	} else if (parser.protocol == 'https:') {
    parser.protocol = "wss:"
}
if (parser.pathname.substr(parser.pathname.length - 1) != '/') {
	parser.pathname += '/';
}
parser.pathname += 'hub'
console.log(
    parser.protocol,
    parser.host,
    parser.hostname,
    parser.port,
    parser.pathname,
    parser.hash,
    parser.search,
    parser.origin,
);

console.log('Here comes the HRF MOFO TOTORO');

console.log(parser.href);

const conn = new WebSocket(parser.href + "?snakesecret=pythonforpresident&clientid=" + uId);

const container = document.querySelector('.container');

let initialRender = true;

const keyMap = {
  w: 'up',
  a: 'left',
  s: 'down',
  d: 'right',
}

// Termbox colors to css named colors
const colorList = [
	'',
	'Black', 
	'Red',
	'Green', 
	'Yellow',
	'Blue', 
	'Fuchsia',
	'Aqua',
	'White', 
];

document.addEventListener('keydown', ({ key }) => {
	// console.log(key, 'kacsa')
	if (['w','a','s','d'].includes(key)) {
		conn.send(JSON.stringify({
			clientID: uId,
			type: 'control',
			data: keyMap[key],
		}));
	}
});

function waitUsers (data) {
// document.body.innerHTML = "";
// const childs = data.textbox.map(e => {
//     let msg = document.createElement('p');
//     msg.innerText = e;
//     return msg;
// })
// // childs.forEach(e => document.body.appendChild(e));
// childs.forEach(e => console.log('WaitUsers', e));
};

function render({ snakes, level, foods, ...rest }) {
	if (initialRender) {
		console.log(rest)
		for (let i = 0; i < level.height; i++) {
			for (let j = 0; j < level.width; j++) {
				let cell = document.createElement('div')
				cell.className = "cell";
				cell.id = j + "_" + i;
				container.appendChild(cell);
			}
		}
		container.style.gridTemplateColumns = "repeat(" +level.width + "," + 100 / level.width + "%)";
		container.style.gridAutoRows = 100 / level.width + "%";
		initialRender = false;
	} else {
		document.querySelectorAll('.cell').forEach(cell => {
			const classListArray = Array.from(cell.classList)
			if (classListArray.includes('haveBeenEaten') && !classListArray.includes('snake')) {
				cell.classList.remove('food', 'haveBeenEaten')
			}
			if (classListArray.includes('head') && classListArray.includes('food')) {
				cell.classList.add('haveBeenEaten')
			}
			cell.classList.remove('snake', 'head');
			
			cell.innerHTML = '';
			cell.style.background = 'white';
		})
		foods.forEach(({ pos, type }) => {
			const cell = document.getElementById(pos.x + "_" + pos.y)
			cell.innerHTML = "<p>" + String.fromCodePoint(type.leftrune) + "</p>";
			cell.classList.add('food');
		})
		snakes.forEach(snake => {
			for (let i = 0; i < snake.body.length; i++) {
				const block = snake.body[i];
				const cell = document.getElementById(block.x + "_" + block.y);
				cell.classList.add('snake');
				cell.style.background = colorList[snake.bgcolor];
				cell.style.color = colorList[snake.color];
				if (i === 0) {
					cell.innerHTML = "<p>" + String.fromCodePoint(snake.headrune) + "</p>";
					cell.classList.add('head');
				} else {
					cell.innerHTML = "<p>" + String.fromCodePoint(snake.leftrune) + String.fromCodePoint(snake.rightrune) + "</p>";
				}
			}
		})
	}
}

function handleMsg(wsMsg) {
	const msg = JSON.parse(wsMsg);
	switch (msg.type) {
		case 'login':
			console.log('Login :', msg.data);
		break;
		case 'stateUpdate':
			const data = JSON.parse(msg.data);
			const renderMethod =  data.scene === 'wait' ? waitUsers : render;
			renderMethod(data);
		break;
	}
}


conn.onclose = evt => {
	console.log('Connection closed. ', evt.data);
};

conn.onmessage = evt => {
	var messages = evt.data.split('\n');
	messages.forEach(msg => {
		msg ? handleMsg(msg) : null
	});
};

conn.onerror = error => {
    console.log(error)
}

conn.onopen = () => {
	conn.send(JSON.stringify({
		clientID: uId,
		type: 'login',
		data: 'pythonforpresident'
	}));
}
</script>
<style>
body {
  background-color: beige;
  margin: 0;
}

.container {
  margin: 0 auto;
  display: grid;
  box-sizing: border-box;
  width: 100vw;
  height: 100vw;
  max-width: 800px;
  max-height: 800px;
  box-sizing: border-box;
  border-bottom: 1px solid red;
}

div.food {
  transform: scale(1.5);
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  mix-blend-mode: difference;
  text-shadow: -.1px 2px 0px #ec7070,
                 -.2px 4px 0px #78dae7,
                 -.3px 6px 0px #7eec83;
}

div.food.head {
  transform: scale(2);
  mix-blend-mode: exclusion;
  text-shadow: none;
}

div.food.snake {
  width: 125%;
  height: 125%;
  text-shadow: none;
}

p {
  margin: 0;
  font-size: 5px;
  font-size: 1.8vw;
  align-self: center;
  text-align: center;
  justify-self: center;
  width: 100%
}

.cell {
  background: rgb(0, 0, 0);
  transition: background-color .5s ease-in-out, border-radius .3s ease-out, opacity 1s ease-out, transform .3s linear;
  display: flex;
  mix-blend-mode: color-burn
}

.snake {
  mix-blend-mode: normal;
  display: flex;
  border-radius: 25%;
  transform: scale(1.15);
  transition: all .3s ease-in-out;
  z-index: 5;
  opacity: .3;
  transition: opacity 5s ease, transform .3s linear;
}

.head {
  z-index: 6;
  mix-blend-mode: normal;
  transform: scale(1.5);
  transition: none;
  opacity: 1;
  animation: .5s morph ease-in infinite;
  -moz-box-shadow: 0 0 5PX #000000;
  -webkit-box-shadow: 0 0 5PX #000000;
  box-shadow: 0 0 5PX #000000;

}

@keyframes morph {
  0%, 100% {
    border-radius: 5% 7% 7% 5% / 7% 15% 7% 15%;
  }
  50% {
    border-radius: 45%;
  }
}     
</style>
</body>
</html>
`))
