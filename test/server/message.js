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
    const [id, role, seqStr, isEndStr, statusStr, mediaType, encoding] = comps;
    const data = encoding === 'base64'
      ? Buffer.from(dataBuff.toString('utf8'), 'base64')
      : dataBuff;
    const seq = Number.parseInt(seqStr);
    const isEnd = Message._parseBool(isEndStr);
    const status = Number.parseInt(statusStr);
    return new Message({id, role, seq, isEnd, status, mediaType, encoding, data}); 
  }
  static _parseBool(strBool, fallback) {
    if (typeof strBool === 'boolean') return strBool;
    strBool = strBool.toUpperCase();
    let value;
    switch (strBool) {
      case 'TRUE':
      case 'T':
      case 'YES':
      case 'Y':
      case '1':
        value = true;
        break;
      case 'FALSE':
      case 'F':
      case 'NO':
      case 'N':
      case '0':
        value = false;
        break;
      default:
        value = fallback;
        break;
    }
    return value;
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
    const buffer = this.cache.nextMessageBuffer();
    this.callback(Message.parse({buffer}));
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
  nextMessageBuffer() {
    const extractedLength = this._extractMessageLength();
    const buffer = this._sliceBuffer(0, extractedLength);
    this._sliceAndKeep(extractedLength);
    return buffer;
  }
  _extractMessageLength() {
    if (this.chunks.length === 0) {
      throw new Error("Can't get message length from empty buffer cache");
    }
    const msgLenBuffs = this._extractMessageLengthChunks();
    if (msgLenBuffs.length === 0) {
      throw new Error("No message length found in buffer cache");
    }
    const msgLenBuffer = Buffer.concat(msgLenBuffs);
    const lengthString = msgLenBuffer.toString('utf8');
    const length = Number.parseInt(lengthString);
    if (length <= 0) {
      throw new Error(`Expecting message length > 0, but found: ${length}`);
    }
    return length;
  }
  _extractMessageLengthChunks() {
    let colonChunkIdx;
    let colonIdx = -1;
    for (let i = 0; i < this.chunks.length; i++) {
      const chunk = this.chunks[i];
      colonIdx = chunk.indexOf(':');
      if (colonIdx > -1) {
        colonChunkIdx = i;
        break;
      }
    }
    if (colonIdx === -1) return [];
    const newChunks = this.chunks.slice(colonChunkIdx);
    const msgLenChunks = this.chunks.slice(0, colonChunkIdx + 1);
    const subject = newChunks[0];
    const msgLenBuff = subject.slice(0, colonIdx);
    const restBuff = subject.slice(colonIdx + 1);
    newChunks[0] = restBuff;
    this.chunks = newChunks;
    msgLenChunks[msgLenChunks.length - 1] = msgLenBuff;
    return msgLenChunks;
  }
  _sliceBuffer(fromIdx, toIdx) {
    if (!this.buffer) {
      this.buffer = Buffer.concat(this.chunks);
      this.chunks = [this.buffer];
    }
    return this.buffer.slice(fromIdx, toIdx);
  }
  _sliceAndKeep(fromIdx, toIdx) {
    let chunk = this._sliceBuffer(fromIdx, toIdx);
    this.chunks = [chunk];
    this.buffer = null; // Important! Otherwise _extractMessageLengthChunks is not working correctly.
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
