import Message from './message';
const stream = require('stream');

export const letterNames = 'ALPHA,BETA,GAMMA,DELTA,EPSILON,ZETA,ETA,THETA,IOTA,KAPPA,LAMBDA,MU,NU,XI,OMICRON,PI,RHO,SIGMA,TAU'.split(',');

export function chunkCacheTestEnvelope(content) {
  const data = Buffer.from(content);
  const lengthBuff = Buffer.from(`${data.length}:`);
  return [data, Buffer.concat([lengthBuff, data])];
}

export function chunkCacheTestChopEnvelope(buffer, count) {
  const acc = [];
  const pieceLen = buffer.length / count;
  for (let i = 0; i < count; i++) {
    acc.push(buffer.slice(i * pieceLen, (i + 1) * pieceLen));
  }
  return acc;
}

export class MockWriter extends stream.Writable {
  constructor(callback, options) {
    super(options);
    this.callback = callback;
  }
  _write(chunk, enc, next) {
    this.callback(chunk);
    next();
  }
}

export class MockReader extends stream.Readable {
  constructor(data, options) {
    super(options);
    this.data = data;
  }
  _read(size) {
    this.push(this.data.slice(0, size));
    this.data = this.data.slice(size);
  }
}

export function messageReaderTestMessages() {
  return letterNames.map((ln, i) => {
    const id = (i + 1).toString(16);
    const role = 'test-message';
    const status = 200;
    const mediaType = 'text/plain';
    const encoding = 'utf8';
    const data = Buffer.from(`TEST-DATA-${ln}`);
    return Message.from({id, role, status, mediaType, encoding, data});
  });
}

