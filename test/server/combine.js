const IO = require('../../node-approx/io');
const io = new IO();
const buff = new Map();

let mediaTypeCount = 0;
let dataAndStatusCount = 0;

io.read((message) => {
  let acc;
  if (buff.has(message.id)) {
    acc = buff.get(message.id);
  } else {
    acc = {
      mediaType: null,
      chunks: [],
    };
  }

  if (message.role === 'media-type') {
    mediaTypeCount++;
    acc.mediaType = message.data.toString();
  } else if (message.role === 'file-content') {
    dataAndStatusCount++;
    const {seq, status, isEnd, data} = message;
    acc.chunks.push({seq, status, isEnd, data});
  }

  io.logger.info(`mediaTypeCount: ${mediaTypeCount}, dataAndStatusCount: ${dataAndStatusCount}\n`);

  if (acc.mediaType && acc.chunks.length) {
    io.logger.info(`combine responds with ${acc.chunks.length} chunks\n`);
    acc.chunks = acc.chunks.filter(({seq, status, data, isEnd}) => {
      io.send({id, seq, role: 'response-chunk', status, mediaType, encoding: 'base64', isEnd, data});
      return false;
    });
  }

  if (message.isEnd && (message.role === 'file-content')) {
    io.logger.info(`combine deletes accumulator for id ${message.id}\n`);
    buff.delete(message.id);
  } else {
    io.logger.info(`combine updates accumulator for id ${message.id}\n`);
    buff.set(message.id, acc);
  }
});
