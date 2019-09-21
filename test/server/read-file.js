const fs = require('fs');
const path = require('path');
const IO = require('../../node-approx/io');
const io = new IO();

io.read((message) => {
  const req = JSON.parse(message.data);
  let filePath = req.url.path;
  if (filePath === '/') {
    filePath = '/index.html';
  }

  let seq = 0;
  readFileB64(filePath, (err, data, isEnd) => {
    if (err) {
      if (err.code === 'ENOENT') {
        io.send({
          id: message.id, 
          role: 'file-content',
          status: 404,
          isEnd,
        });
      } else {
        const errData = err.toString();
        io.send({
          id: message.id, 
          role: 'file-content',
          status: 500,
          encoding: 'utf8',
          isEnd,
          data: errData,
        });
      }
    } else {
      io.send({
        id: message.id,
        seq, 
        role: 'file-content',
        status: 200,
        encoding: 'base64',
        isEnd,
        data,
      });
    }
    seq++;
  });
});

function readFileB64(filePath, callback) {
  const fullPath = path.join(process.cwd(), process.env.BASE_PATH, filePath);
  const reader = fs.createReadStream(fullPath);
  let lastData = null;
  reader
    .on('error', (err) => {
      callback(err, null, true);
    }) 
    .on('data', (chunk) => {
      if (lastData) {
        callback(null, lastData, false);
      }
      lastData = chunk.toString('base64');
    })
    .on('end', () => {
      callback(null, lastData, true);
    });
}