const fs = require('fs');
const axmsg = require('./axmsg');
const errors = new axmsg.Errors('read_file');
const reader = new axmsg.Reader(process.stdin);
const writer = new axmsg.Writer(process.stdout);

reader.on((action) => {
  // input example:
  // {"axmsg":1,"id":21,"role":"http-request","data":{"method":"GET","filePath":"../static/files/index.html"}}
  // {"axmsg":1,"id":21,"role":"http-request","data":{"method":"GET","filePath":"./test.json"}}
  if (action.axmsg === 1 && action.role === 'http-request' && action.data.method === 'GET') {
    readFileB64(action.data.filePath, (err, bodyB64) => {
      if (err) {
        if (err.code === 'ENOENT') {
          writer.write({
            id: action.id,
            role: 'response-data',
            cmd: 'respond-with-error',
          }, {
            status: 404,
          });
        } else {
          errors.log(action.id, err);
        }
      } else {
        writer.write({
          id: action.id,
          role: 'response-data',
          cmd: 'respond-with-data',
        }, {
          status: 200,
          bodyB64,
        });
      }
    });
  }
});

function readFileB64(filePath, callback) {
  fs.readFile(filePath, (err, data) => {
    if (err) {
      callback(err);
      return;
    }
    const buf = Buffer.from(data);
    callback(null, buf.toString('base64'));
  });
}