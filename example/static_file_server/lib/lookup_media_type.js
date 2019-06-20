const fs = require('fs');
const path = require('path');
const axmsg = require('./axmsg');
const axenvs = require('./axenvs');
const errors = new axmsg.Errors('lookup_media_type');
const envs = new axenvs.Envs('lookup_media_type');

readMediaTypeForExtension((err, mediaTypeForExt) => {
  if (err) {
    errors.logFatal(err);
  } else {
    openInputOutputAndListen(mediaTypeForExt);
  }
});

function openInputOutputAndListen(mediaTypeForExt) {
  const [ins, outs] = envs.insOuts();
  let fatalErrors = false;
  if (ins.length !== 1) {
    errors.log(null, new Error(`Exactly 1 input expected, but got: ${ins.length}`));
    fatalErrors = true;
  }
  if (outs.length !== 1) {
    errors.log(null, new Error(`Exactly 1 output expected, but got: ${out.length}`));
    fatalErrors = true;
  }
  if (fatalErrors) {
    errors.logFatal(null, new Error("Too many fatal errors, can't run"));
  }
  startListening(ins[0], outs[0], mediaTypeForExt);
}

function startListening(reader, writer, mediaTypeForExt) {
  reader.on((action) => {
    // input examples:
    // {"axmsg":1,"id":21,"role":"http-request","data":{"method":"GET","filePath":"../static/files/index.html"}}
    // {"axmsg":1,"id":21,"role":"http-request","data":{"method":"GET","filePath":"../static/files/images/logo.png"}}
    if (action.axmsg === 1 && action.role === 'http-request' && action.data.method === 'GET') {
      const rawExt = path.extname(action.data.filePath);
      const extension = rawExt === '' ? '.html' : rawExt;
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