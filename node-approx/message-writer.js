export default class MessageWriter {
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
