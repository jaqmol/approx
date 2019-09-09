function writer(reader, stream) {
  return (msg) => {
    if (!stream.write(msg, 'utf8')) {
      reader.pause();
      stream.once('drain', () => {
        reader.resume();
      });
    }
  };
}

module.exports = writer;
