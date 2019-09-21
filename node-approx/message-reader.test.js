import MessageReader from './message-reader';
import { 
  messageReaderTestMessages,
  MockReader,
} from './test-utils';

test('MessageReader test', (done) => {
  const messages = messageReaderTestMessages();
  const envelopes = messages.map((m) => m.envelope());
  const input = new MockReader(Buffer.concat(envelopes));
  const reader = new MessageReader({stream: input});
  const readMessages = [];

  const conclude = () => {
    for (let i = 0; i < readMessages.length; i++) {
      const readMsg = readMessages[i];
      const origMsg = messages[i];
      const readBuff = readMsg.envelope();
      const origBuff = origMsg.envelope();
      expect(readBuff.equals(origBuff)).toBeTruthy();
    }
    done();
  };

  reader.read((msg) => {
    readMessages.push(msg);
    if (readMessages.length === messages.length) {
      conclude();
    }
  });
});
