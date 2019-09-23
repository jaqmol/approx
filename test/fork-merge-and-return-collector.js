const IO = require('../node-approx/io');
const io = new IO();

const collector = {};

io.read((message) => {
  let acc = collector[message.id];
  if (!acc) {
    acc = [];
  }
  const strData = message.data.toString();
  acc.push(JSON.parse(strData));
  if (acc.length === 2) {
    send(message.id, acc);  
    delete collector[message.id];
  } else {
    collector[message.id] = acc;
  }
});

function send(id, allMsgs) {
  const dataStr = JSON.stringify(allMsgs, null, 2);
  const data = Buffer.from(dataStr).toString('base64');
  io.send({
    id,
    role: 'response',
    status: 200,
    mediaType: 'application/json',
    encoding: 'base64',
    data,
  });
}