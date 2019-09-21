const MessageReader = require('./message-reader');
const MessageWriter = require('./message-writer');
const Message = require('./message');

class IO {
  constructor() {
    const reader = new MessageReader({stream: process.stdin});
    this.reader = reader;
    this.writer = new MessageWriter({reader, stream: process.stdout});
    const logWriter = new MessageWriter({reader, stream: process.stderr});
    this.logger = new Logger({writer: logWriter});
  }
  
  read(callback) {
    this.reader.read(callback);
  }

  send({id, role, seq, isEnd, status, mediaType, encoding, data}) {
    const message = new Message({id, role, seq, isEnd, status, mediaType, encoding, data});
    this.writer.write({message});
  }
}

class Logger {
  constructor({writer}) {
    this.writer = writer;
  }
  info(...args) {
    this.log('info', ...args);
  }
  
  fatal(...args) {
    this.log('fatal', ...args);
  }
  log(role, ...args) {
    const data = args.map(a => `${a}`).join(' ');
    const message = new Message({role, data});
    this.writer.write({message});
  }
}

module.exports = IO;
