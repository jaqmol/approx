const fs = require('fs');
const path = require('path');

let func = readAndGet;

async function readAndGet(extension) {
  const mediaTypeForExt = await readCsv();
  func = createGet(mediaTypeForExt);
  return func(extension);
}

function createGet(mediaTypeForExt) {
  return async (extension) => {
    return mediaTypeForExt[extension];
  };
}

function readCsv() {
  return new Promise((resolve, reject) => {
    const csvFilePath = path.join(__dirname, 'media-types.csv');
    fs.readFile(csvFilePath, 'utf8', (err, data) => {
      if (err) {
        reject(err);
      } else {
        const acc = {};
        const lines = data.split('\n');
        for (let l of lines) {
          const cell = l.split(',');
          acc[cell[0]] = cell[1];
        }
        resolve(acc);
      }
    });
  });
}

module.exports = (extension) => func(extension);
