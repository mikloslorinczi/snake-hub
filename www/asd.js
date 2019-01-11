const container = document.body;
for (let i = 0; i < 25; i++) {
    for (let j = 0; j < 25; j++) {
        let cell = document.createElement('div')
        cell.className = "cell";
        cell.id = `${i}_${j}`
        container.appendChild(cell);        
    }
}


const waitUsers = data => { 
    // document.body.innerHTML = "";
    // const childs = data.textbox.map(e => {
    //     let msg = document.createElement('p');
    //     msg.innerText = e;
    //     return msg;
    // })
    // childs.forEach(e => document.body.appendChild(e));
};

const render = data => {
    document.querySelectorAll('.cell').forEach(cell => {
        cell.classList = ['cell'];
    })
    data.snakes.forEach(snake => {
        snake.body.forEach(b => {
            let cell = document.getElementById(`${b.x}_${b.y}`)
            cell.classList.add('snake');
        })
    })
    console.log(data);
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

conn.onopen = () => {
    conn.send(JSON.stringify({
        clientID: `${new Date().getTime()}`,
        type: 'login',
        data: 'pythonforpresident'
    }))
}
