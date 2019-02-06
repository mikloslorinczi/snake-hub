'use strict';

let oldData = {};

let gameState = 'initial';
let initialRender = true;

const loginMenu = document.getElementById('loginMenu');
const messageBox = document.getElementById('messageBox');
const gameScene = document.getElementById('gameScene');

const uId = Math.random().toString(36).substring(2) + (new Date()).getTime().toString(36);

let urlParser = document.createElement('a');
urlParser.href = window.location.href;
urlParser.protocol = (urlParser.protocol == 'https:') ? 'wss:' : 'ws:';
if (urlParser.pathname.substr(urlParser.pathname.length - 1) != '/') { urlParser.pathname += '/' };
urlParser.pathname += 'hub';

const conn = new WebSocket(urlParser.href + '?snakesecret=pythonforpresident&clientid=' + uId);

const keyMap = {
  w: 'up',
  a: 'left',
  s: 'down',
  d: 'right',
};

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

function listenKeys() {
	document.addEventListener('keydown', ({ key }) => {
		if (['w','a','s','d'].includes(key)) {
			conn.send(JSON.stringify({
				clientID: uId,
				type: 'control',
				data: keyMap[key],
			}));
		}
	});
}

function renderMessageBox(data) {
	messageBox.innerHTML = '';
	data.forEach(element => {
		let line = document.createElement('p');
		line.innerHTML = element;
		messageBox.appendChild(line);
	});
};

function drawGameScene(level) {
	for (let y = 0; y < level.height; y++) {
		for (let x = 0; x < level.width; x++) {
			let cell = document.createElement('div')
			cell.id = x + '_' + y;
			gameScene.appendChild(cell);
		}
	}
	gameScene.style.gridTemplateColumns = 'repeat(' +level.width + ',' + 100 / level.width + '%)';
	gameScene.style.gridAutoRows = 100 / level.width + '%';
}

function drawSnakes(snakes, hide = false) {
	snakes.forEach(snake => {
		for (let i = 0; i < snake.body.length; i++) {
			const block = snake.body[i];
			const cell = document.getElementById(block.x + '_' + block.y);
			if (hide) {
				cell.classList = '';
				cell.style.color = '';
				cell.style.background = '';
				cell.innerHTML = '';
			} else {
				cell.classList.add('snake');
				cell.style.background = colorList[snake.bgcolor];
				cell.style.color = colorList[snake.color];
				if (i === 0) {
					cell.innerHTML = '<p>' + String.fromCodePoint(snake.headrune) + '</p>';
					cell.classList.add('head');
				} else {
					cell.innerHTML = '<p>' + String.fromCodePoint(snake.leftrune) + String.fromCodePoint(snake.rightrune) + '</p>';
				};
			};
		};
	});
};

function drawFoods(foods, hide = false) {
	foods.forEach(({ pos, type }) => {
		const cell = document.getElementById(pos.x + '_' + pos.y);
		if (hide) {
			cell.classList = '';
			cell.style.color = '';
			cell.style.background = '';
			cell.innerHTML = '';
		} else {
			cell.innerHTML = '<p>' + String.fromCodePoint(type.leftrune) + '</p>';
			cell.classList.add('food');
		};
	});
};

function renderGameScene(data) {

	if (!initialRender) {
		drawFoods(oldData.foods, true);
		drawSnakes(oldData.snakes, true);
	} else {
		initialRender = false;
	};

	drawFoods(data.foods);
	drawSnakes(data.snakes);
	
	oldData = data;

};

function changeState(data) {
	loginMenu.style.display = 'none';
	messageBox.style.display = 'none';
	gameScene.style.display = 'none';
	switch (data.scene) {
		case 'wait':
			messageBox.style.display = 'block';
		break;
		case 'scores':
			messageBox.style.display = 'block';
		break;
		case 'game':
			if (initialRender) {
				drawGameScene(data.level);
				listenKeys();
			};
			gameScene.style.display = 'grid';
		break;
	};
	gameState = data.scene;
};

function handleMsg(wsMsg) {
	const msg = JSON.parse(wsMsg);
	switch (msg.type) {
		case 'login':
			console.log('Login :', msg.data);
		return;
		case 'stateUpdate':
			const data = JSON.parse(msg.data);
			if (data.scene != gameState) { changeState(data) }
			switch (data.scene) {
				case 'game':
					renderGameScene(data);
					return;
				case 'wait':
					renderMessageBox(data.textbox);
					return;
				case 'scores':
					renderMessageBox(data.textbox);
					return;
			};
		return;
	};
};

function connOpen(event){
	console.log(event);
	conn.send(JSON.stringify({
		clientID: uId,
		type: 'login',
		data: 'pythonforpresident'
	}));
};

conn.onopen = event => { connOpen(event) };

conn.onmessage = event => { handleMsg(event.data) };

conn.onclose = event => { console.log('Connection closed. ', event.data) };

conn.onerror = err => { console.log('WebSocket error: ', err) };
