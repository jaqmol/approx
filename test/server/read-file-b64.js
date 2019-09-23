const fs = require('fs');
const path = require('path');

function readFileB64(basePath, relativePath) {
  return new Promise((resolve, reject) => {
    const fullPath = path.join(process.cwd(), basePath, relativePath);
    fs.readFile(fullPath, (err, buff) => {
      if (err) {
        reject(err);
      } else {
        resolve(buff.toString('base64'));
      }
    });
  });
}

module.exports = readFileB64;
