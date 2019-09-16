export class Message {
  constructor({id, role, seq, isEnd, status, mediaType, encoding, data}) {
    this.header = {
      id: id || '',
      role: role || '',
      seq: typeof seq === 'number' ? seq : 0,
      isEnd: typeof isEnd === 'boolean' ? isEnd : true,
      status: typeof status === 'number' ? status : 0,
      mediaType: mediaType || '',
      encoding: encoding || '',
    };
    if (data) {
      this.data = data instanceof Buffer ? data : Buffer.from(data);
    } else {
      this.data = Buffer.from('');
    }
  }
  static from({id, role, seq, isEnd, status, mediaType, encoding, data}) {
    return new Message({id, role, seq, isEnd, status, mediaType, encoding, data});
  }
  static parse({buffer}) {
    const semicolonIndex = buffer.indexOf(';');
    if (semicolonIndex === -1) throw new Error(`Cannot parse message: ${buffer.toString('utf8')}`);
    const msgBuff = buffer.slice(0, semicolonIndex);
    const dataBuff = buffer.slice(semicolonIndex + 1);
    const comps = msgBuff.toString('utf8').split(',');
    const [id, role, seq, isEnd, status, mediaType, encoding] = comps;
    let data;
    if (encoding === 'base64') {
      data = Buffer.from(dataBuff.toString('utf8'), 'base64');
    } else {
      data = dataBuff;
    }
    return new Message({id, role, seq, isEnd, status, mediaType, encoding, data}); 
  }
  _headerBuffer() {
    const comps = [
      this.header.id,
      this.header.role,
      this.header.seq,
      this.header.isEnd,
      this.header.status,
      this.header.mediaType,
      this.header.encoding,
      ';',
    ];
    const str = comps.join(',');
    return Buffer.from(str, 'utf8');
  }
  toBuffer() {
    const headerBuff = this._headerBuffer();
    const length = headerBuff.length + this.data.length;
    const lengthBuff = Buffer.from(`${length}:`);
    return Buffer.concat(
      [lengthBuff, headerBuff, this.data], 
      lengthBuff.length + length,
    );
  }
}

export class MessageReader {
  constructor({stream}) {
    this.cache = new ChunkCache();
    this.nextMessageLength = -1;
    stream.on('data', c => this._processChunk(c));
  }
  _processChunk(chunk) {
    this.cache.push(chunk);
    this._processNextMessageLength();
    this._processMessage();
  }
  _processNextMessageLength() {
    if (this.nextMessageLength) return;
    this.nextMessageLength = this.cache.messageLength(chunk);
  }
  _processMessage() {
    if (this.nextMessageLength > this.cache.length) return;
    const msgBuff = this.cache.sliceBuffer(0, this.nextMessageLength);
    this.cache.sliceInPlace(this.nextMessageLength);
    this.nextMessageLength = -1;
    this.callback(Message.parse({buffer: msgBuff}));
  }
  read(callback) {
    this.callback = callback;
  }
}

export class ChunkCache {
  constructor() {
    this.chunks = [];
    this.buffer = null;
  }
  get length() {
    return this.chunks.reduce((a, c) => a + c.length, 0);
  }
  push(chunk) {
    this.chunks.push(chunk);
    this.buffer = null;
  }
  messageLength() {
    if (this.chunks.length === 0) {
      throw new Error("Can't get message length from empty buffer cache");
    }
    const buff = this.cache[0];
    const colonIndex = buff.indexOf(':');
    if (colonIndex === -1) return -1;
    const lengthString = buff.slice(0, colonIndex).toString('utf8');
    const length = Number.parseInt(lengthString);
    if (length <= 0) {
      throw new Error(`Expect message length to be positive, but found: ${length}`);
    }
    this.cache[0] = buff.slice(colonIndex + 1);
    return length;
  }
  sliceBuffer(fromIdx, toIdx) {
    if (!this.buffer) {
      this.buffer = Buffer.concat(this.chunks);
      this.chunks = [this.buffer];
    }
    return this.buffer.slice(fromIdx, toIdx);
  }
  sliceInPlace(fromIdx, toIdx) {
    let chunk = this.sliceBuffer(fromIdx, toIdx);
    chunk = Uint8Array.prototype.slice.call(chunk);
    this.chunks = [chunk];
    this.buffer = chunk;
  }
}

export class MessageWriter {
  constructor({reader, stream}) {
    if (reader) {
      this._writeMessageData = (data) => {
        if (!stream.write(data)) {
          reader.pause();
          stream.once('drain', () => {
            reader.resume();
          });
        }
      };
    } else {
      this._writeMessageData = (data) => {
        stream.write(data);
      };
    }
  }
  write({message}) {
    const data = message.toBuffer();
    this._writeMessageData(data);
  }
}
