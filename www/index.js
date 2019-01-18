const container = document.querySelector('.container');
let initialRender = true;
const keyMap = {
    w: 'up',
    d: 'right',
    s: 'down',
    a: 'left',
}

document.addEventListener('keydown', ({ key }) => {
    // console.log(e)
    conn.send(JSON.stringify({
        clientID: `${new Date().getTime()}`,
        type: 'control',
        data: keyMap[key],
    }))
})


const waitUsers = data => {
    // document.body.innerHTML = "";
    // const childs = data.textbox.map(e => {
    //     let msg = document.createElement('p');
    //     msg.innerText = e;
    //     return msg;
    // })
    // childs.forEach(e => document.body.appendChild(e));
};

const render = ({ snakes, level, foods, ...rest }) => {
    if (initialRender) {
        console.log(rest)
        for (let i = 0; i < level.height; i++) {
            for (let j = 0; j < level.width; j++) {
                let cell = document.createElement('div')
                cell.className = "cell";
                cell.id = `${j}_${i}`
                container.appendChild(cell);
            }
        }
        container.style.gridTemplateColumns = `repeat(${level.width}, ${100 / level.width}%)`;
        container.style.gridAutoRows = `${100 / level.width}%`;
        initialRender = false;
    } else {
        document.querySelectorAll('.cell').forEach(cell => {
            cell.classList = ['cell'];
            cell.innerHTML = '';
        })
        foods.forEach(({ pos, type }) => {
            const cell = document.getElementById(`${pos.x}_${pos.y}`)
            cell.innerHTML = `<p>${String.fromCodePoint(type.leftrune)}</p>`
            cell.classList.add('food');
        })
        snakes.forEach(snake => {
            snake.body.forEach(block => {
                const cell = document.getElementById(`${block.x}_${block.y}`)
                cell.classList.add('snake');
            })
        })
    }
    // console.log(data);
}

const appendLog = log => {
    // console.log(JSON.parse(JSON.parse(log).data));
    const data = JSON.parse(JSON.parse(log).data);
    // console.log(data.scene )
    const renderMethod =  data.scene === 'wait' ? waitUsers : render;
    renderMethod(data);
}

conn = new WebSocket(`ws://localhost:4545/hub?snakesecret=pythonforpresident&clientid=${new Date().getTime()}`);

conn.onclose = function (evt) {
    var item = document.createElement("div");
    appendLog(item );
};

conn.onmessage = function (evt) {
    var messages = evt.data.split('\n');
    messages.forEach(msg => {
        msg ? appendLog(msg) : null
    });
};

conn.onerror = error => {
    console.log(error)
}

conn.onopen = () => {
    conn.send(JSON.stringify({
        clientID: `${new Date().getTime()}`,
        type: 'login',
        data: 'pythonforpresident'
    }))
}
