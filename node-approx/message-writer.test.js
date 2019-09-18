import MessageWriter from './message-writer';
import {
  messageReaderTestMessages,
  MockWriter,
} from './test-utils'

test('MessageWriter test', (done) => {
  const messages = messageReaderTestMessages();
  const writtenEnvelopes = [];

  const conclude = () => {
    for (let i = 0; i < writtenEnvelopes.length; i++) {
      const writtenBuff = writtenEnvelopes[i];
      const origBuff = messages[i].toBuffer();
      expect(writtenBuff.equals(origBuff)).toBeTruthy();
    }
    done();
  };

  const output = new MockWriter((env) => {
    writtenEnvelopes.push(env);
    if (writtenEnvelopes.length === messages.length) {
      conclude();
    }
  });
  const writer = new MessageWriter({stream: output});
  
  for (const msg of messages) {
    writer.write({message: msg});
  }
});