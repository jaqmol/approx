const readline = require('readline');

class Errors {
  constructor(source) {
    this.source = source;
  }
  log(id, errorObject) {
    const errActn = {
      axmsg: 1,
      id,
      role: 'error',
      data: {
        source: this.source,
        code: 10001,
        message: errorObject.message,
        fileName: errorObject.fileName,
        lineNumber: errorObject.lineNumber,
      },
    };
    const errStr = JSON.stringify(errActn) + '\n';
    process.std.write(errStr, 'utf8');
  }
  logFatal(id, errorObject) {
    this.log(id, errorObject);
    process.exit(1);
  }
}

class Reader {
  constructor(readableStream) {
    this.readable = readableStream
  }
  on(callback) {
    const lineReader = readline.createInterface({
      input: this.readable,
      crlfDelay: Infinity,
    });
    lineReader.on('line', (line) => {
      const action = JSON.parse(line);
      callback(action);
    });
  }
}

class Writer {
  constructor(writableStream) {
    this.writable = writableStream;
  }
  write(params, data) {
    const action = Object.assign({axmsg: 1}, params);
    action.data = data;
    const output = JSON.stringify(action) + '\n';
    this.writable.write(output, 'utf8');
  }
}

module.exports = {
  Errors,
  Reader,
  Writer,
};