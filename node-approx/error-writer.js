const { Transform } = require('stream');
const { EvntEnd } = require('./config');

class ErrorWriter extends Transform {
  constructor(outputStream) {
    super({objectMode: true});
    if (outputStream) {
      this.pipe(outputStream);
    }
  }

  _transform(error, enc, callback) {
    const obj = {
      columnNumber: error.columnNumber,
      fileName: error.fileName, 
      lineNumber: error.lineNumber, 
      message: error.message,
      name: error.name,
      stack: error.stack 
    }

    try {
      const str = `${JSON.stringify(obj)}${EvntEnd}`;
      const chunk = Buffer.from(str);
      callback(null, chunk);
    } catch(err) {
      callback(err);
    }
  }
}

module.exports = ErrorWriter;