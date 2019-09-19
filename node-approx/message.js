export default class Message {
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
      case '1':
        value = true;
        break;
      case 'FALSE':
      case 'F':
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
    ];
    const str = comps.join(',');
    return Buffer.from(`${str};`, 'utf8');
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
