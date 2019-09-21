const IO = require('../node-approx/io');
const io = new IO();

io.read((message) => {
  const {id, data} = message;
  io.send({id, role: 'response', status: 200, mediaType: 'application/json', encoding: 'utf8', data});
});
