class ChunkCache {
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
    const [len, ok] = this._extractMessageLength();
    if (!ok) return null;
    const buffer = this._sliceBuffer(0, len);
    this._sliceAndKeep(len);
    return buffer;
  }
  _extractMessageLength() {
    if (this.chunks.length === 0) {
      return [null, false];
      // throw new Error("Can't get message length from empty buffer cache");
    }
    const msgLenBuffs = this._extractMessageLengthChunks();
    if (msgLenBuffs.length === 0) {
      return [null, false];
      // throw new Error("No message length found in buffer cache");
    }
    const msgLenBuffer = Buffer.concat(msgLenBuffs);
    const lengthString = msgLenBuffer.toString('utf8');
    const length = Number.parseInt(lengthString);
    if (length <= 0) {
      return [null, false];
      // throw new Error(`Expecting message length > 0, but found: ${length}`);
    }
    return [length, true];
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

module.exports = ChunkCache;
