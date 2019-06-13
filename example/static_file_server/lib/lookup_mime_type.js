const fs = require('fs');
const readline = require('readline');

readMediaTypeForExtension((err, mediaTypeForExt) => {
  if (err) {
    console.error(err);
    process.exit(1);
  } else {
    listenForInput(mediaTypeForExt);
  }
})

function listenForInput(mediaTypeForExt) {
  const rl = readline.createInterface({
    input: process.stdin, // fs.createReadStream('sample.txt'),
    crlfDelay: Infinity
  });

  rl.on('line', (line) => {
    // expecting AX action message:
    // {
    //   axmsg: int,
    //   id: int,
    //   seq: int,
    //   role: string,
    //   cmd: string,
    //   data: {
    //     method: string,
    //     url: string,
    //     header: object,
    //     bodyB64: string,
    //   },
    // }
    // example:
    // {"axmsg":1,"id":21,"role":"http-request","data":{"method":"GET","url":"http://localhost:3000/index.html"}}
    const inActn = JSON.parse(line);
    if (inActn.axmsg === 1 && inActn.role === 'http-request' && inActn.data.method === 'GET') {
      let extension = filenameExtension(inActn.data.url);
      extension = extension === '' ? '.html' : extension;
      const mediaType = mediaTypeForExt[extension];
      const outActn = {
        axmsg: 1,
        id: inActn.id,
        seq: inActn.seq,
        role: 'media-type',
        cmd: 'set-content-type',
        data: {
          extension,
          mediaType,
        },
      };
      const outputStr = JSON.stringify(outActn) + '\n';
      process.stdout.write(outputStr, 'utf8');
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