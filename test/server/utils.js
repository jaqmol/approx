function Writer({reader, stream}) {
  return (msg) => {
    if (!stream.write(msg, 'utf8')) {
      reader.pause();
      stream.once('drain', () => {
        reader.resume();
      });
    }
  };
}

const Header = {
  parse: (msgStr) => {
    if (!msgStr) return [];
    const semicolonIdx = msgStr.indexOf(';');
    const headerStr = msgStr.slice(0, semicolonIdx);
    const comps = headerStr.split(',');
    return [
      {
        id: comps[0],
        role: comps[1],
        isEnd: JSON.parse(comps[2]),
        mediaType: comps[3],
        encoding: comps[4],
      },
      msgStr.slice(semicolonIdx + 1),
    ];
  },
  stringify: ({id, role, seq, isEnd, status, mediaType, encoding}) => {
    id = id || '';
    role = role || '';
    seq = typeof seq === 'number' ? seq : 0;
    isEnd = typeof isEnd === 'boolean' ? isEnd : true;
    status = typeof status === 'number' ? status : 0;
    mediaType = mediaType || '';
    encoding = encoding || '';
    return `${id},${role},${seq},${isEnd},${status},${mediaType},${encoding};`;
  },
};

module.exports = {
  Writer,
  Header,
};
