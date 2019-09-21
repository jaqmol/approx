const fs = require('fs');
const path = require('path');
const IO = require('../../node-approx/io');
const io = new IO();

readMediaTypes((err, mediaTypeForExt) => {
  if (err) {
    io.logger.fatal(err.toString());
    process.exit(-1);
  } else {
    listen(mediaTypeForExt);
  }
});

function listen(mediaTypeForExt) {
  io.read((message) => {
    io.logger.info(message.data.toString().slice(0, 50));
    const req = JSON.parse(message.data);
    const rawExt = path.extname(req.url.path);
    const extension = rawExt === '' ? '.html' : rawExt;
    io.send({
      id: message.id, 
      role: 'media-type', 
      mediaType: 'text/plain', 
      data: mediaTypeForExt[extension],
    });
  });
}

function readMediaTypes(callback) {
  const csvFilePath = path.join(__dirname, 'media-types.csv');
  fs.readFile(csvFilePath, 'utf8', (err, data) => {
    if (err) {
      callback(err);
      return;
    }
    const acc = {};
    const lines = data.split('\n');
    for (let l of lines) {
      const cell = l.split(',');
      acc[cell[0]] = cell[1];
    }
    callback(null, acc);
  });
}
