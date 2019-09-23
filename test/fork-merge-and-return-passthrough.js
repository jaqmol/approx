const IO = require('../node-approx/io');
const io = new IO();

io.read((message) => {
  let req = JSON.parse(message.data);
  let data = JSON.stringify(req);
  io.send({
    id: message.id,
    role: 'passthrought',
    mediaType: 'application/json',
    encoding: 'utf8',
    data,
  });
});
