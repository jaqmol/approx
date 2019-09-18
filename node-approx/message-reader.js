import ChunkCache from './chunk-cache';
import Message from './message';

export default class MessageReader {
  constructor({stream}) {
    this.cache = new ChunkCache();
    this.nextMessageLength = -1;
    this.stream = stream;
  }
  _processChunk(chunk) {
    this.cache.push(chunk);
    let buffer = this.cache.nextMessageBuffer();
    while (buffer) {
      this.callback(Message.parse({buffer}));
      buffer = this.cache.nextMessageBuffer();
    }
  }
  read(callback) {
    if (this.callback) throw new Error('Read must be called once');
    this.callback = callback;
    this.stream.on('data', c => this._processChunk(c));
  }
}
