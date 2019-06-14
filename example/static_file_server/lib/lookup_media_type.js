const fs = require('fs');
const axmsg = require('./axmsg');
const errors = new axmsg.Errors('lookup_media_type');
const reader = new axmsg.Reader(process.stdin);
const writer = new axmsg.Writer(process.stdout);

readMediaTypeForExtension((err, mediaTypeForExt) => {
  if (err) {
    errors.logFatal(err);
  } else {
    listenForInput(mediaTypeForExt);
  }
});

function listenForInput(mediaTypeForExt) {
  reader.on((action) => {
    // input example:
    // {"axmsg":1,"id":21,"role":"http-request","data":{"method":"GET","url":"http://localhost:3000/index.html"}}
    // {"axmsg":1,"id":21,"role":"http-request","data":{"method":"GET","url":"http://localhost:3000/images/logo.png"}}
    if (action.axmsg === 1 && action.role === 'http-request' && action.data.method === 'GET') {
      let extension = filenameExtension(action.data.url);
      extension = extension === '' ? '.html' : extension;
      const mediaType = mediaTypeForExt[extension];
      writer.write({
        id: action.id,
        seq: action.seq,
        role: 'media-type',
        cmd: 'set-content-type',
      }, {
        extension,
        mediaType,
      });
    }
  });
}

function filenameExtension(urlString) {
  const url = new URL(urlString);
  const comps = url.pathname.split('/');
  const filename = comps[comps.length-1];
  return filename.slice(filename.lastIndexOf('.'));
}

function readMediaTypeForExtension(callback) {
  fs.readFile('media-type-for-extension.csv', 'utf8', (err, data) => {
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