import createRaylib from './raylib.js';

const errnoCodes = { 2: 'EACCES', 8: 'EBADF', 20: 'EEXIST', 28: 'EINVAL', 29: 'EIO', 31: 'EISDIR', 44: 'ENOENT', 52: 'ENOSYS', 54: 'ENOTDIR', 55: 'ENOTEMPTY', 63: 'EPERM' };

function goFS(FS) {
  const decoder = new TextDecoder('utf-8');
  const output = { 1: '', 2: '' };
  const console_ = (fd, buf) => {
    output[fd] += decoder.decode(buf);
    const nl = output[fd].lastIndexOf('\n');
    if (nl !== -1) {
      (fd === 2 ? console.error : console.log)(output[fd].substring(0, nl));
      output[fd] = output[fd].substring(nl + 1);
    }
    return buf.length;
  };
  const stat = (s) => ({
    dev: s.dev, ino: s.ino, mode: s.mode, nlink: s.nlink, uid: s.uid, gid: s.gid, rdev: s.rdev, size: s.size,
    blksize: s.blksize, blocks: s.blocks, atimeMs: +s.atime, mtimeMs: +s.mtime, ctimeMs: +s.ctime,
    isDirectory() { return FS.isDir(s.mode); },
  });
  const call = (fn) => (...args) => {
    const callback = args.pop();
    let result;
    try {
      result = fn(...args);
    } catch (e) {
      if (!(e && 'errno' in e)) throw e;
      const err = new Error(errnoCodes[e.errno] || 'EIO');
      err.code = errnoCodes[e.errno] || 'EIO';
      callback(err);
      return;
    }
    callback(null, result);
  };
  const pos = (p) => (p === null || p === undefined) ? undefined : p;
  return {
    constants: { O_WRONLY: 1, O_RDWR: 2, O_CREAT: 64, O_EXCL: 128, O_TRUNC: 512, O_APPEND: 1024, O_DIRECTORY: 65536 },
    writeSync(fd, buf) {
      if (fd === 1 || fd === 2) return console_(fd, buf);
      return FS.write(FS.getStream(fd), buf, 0, buf.length);
    },
    write: call((fd, buf, offset, length, position) => {
      if (fd === 1 || fd === 2) return console_(fd, buf.subarray(offset, offset + length));
      return FS.write(FS.getStream(fd), buf, offset, length, pos(position));
    }),
    read: call((fd, buf, offset, length, position) => FS.read(FS.getStream(fd), buf, offset, length, pos(position))),
    open: call((path, flags, mode) => FS.open(path, flags, mode).fd),
    close: call((fd) => { FS.close(FS.getStream(fd)); }),
    fstat: call((fd) => stat(FS.fstat(fd))),
    stat: call((path) => stat(FS.stat(path))),
    lstat: call((path) => stat(FS.lstat(path))),
    readdir: call((path) => FS.readdir(path).filter((n) => n !== '.' && n !== '..')),
    mkdir: call((path, perm) => { FS.mkdir(path, perm); }),
    rmdir: call((path) => { FS.rmdir(path); }),
    unlink: call((path) => { FS.unlink(path); }),
    rename: call((from, to) => { FS.rename(from, to); }),
    chmod: call((path, mode) => { FS.chmod(path, mode); }),
    fchmod: call((fd, mode) => { FS.fchmod(fd, mode); }),
    chown: call(() => {}),
    fchown: call(() => {}),
    lchown: call(() => {}),
    truncate: call((path, length) => { FS.truncate(path, length); }),
    ftruncate: call((fd, length) => { FS.ftruncate(fd, length); }),
    fsync: call(() => {}),
    utimes: call((path, atime, mtime) => { FS.utime(path, atime * 1000, mtime * 1000); }),
    readlink: call((path) => FS.readlink(path)),
    symlink: call((target, path) => { FS.symlink(target, path); }),
    link: call(() => { throw { errno: 52 }; }),
  };
}

function untar(FS, data) {
  const text = (off, len) => {
    let end = off;
    while (end < off + len && data[end] !== 0) end++;
    return new TextDecoder().decode(data.subarray(off, end));
  };
  let longName = null;
  for (let off = 0; off + 512 <= data.length;) {
    if (data[off] === 0) break;
    const size = parseInt(text(off + 124, 12).trim() || '0', 8);
    const type = String.fromCharCode(data[off + 156] || 48);
    const prefix = text(off + 345, 155);
    let name = longName ?? (prefix ? prefix + '/' : '') + text(off, 100);
    longName = null;
    const bodyOff = off + 512;
    const body = data.subarray(bodyOff, bodyOff + size);
    off = bodyOff + Math.ceil(size / 512) * 512;
    if (type === 'L') {
      longName = text(bodyOff, size);
      continue;
    }
    if (type === 'x') {
      const m = /\d+ path=([^\n]*)\n/.exec(new TextDecoder().decode(body));
      if (m) longName = m[1];
      continue;
    }
    name = '/' + name.replace(/^(\.\/)+/, '').replace(/\/+$/, '');
    if (name === '/') continue;
    if (type === '5') {
      FS.mkdirTree(name);
    } else if (type === '0' || type === '7') {
      FS.mkdirTree(name.substring(0, name.lastIndexOf('/')) || '/');
      FS.writeFile(name, body);
    }
  }
}

export async function run({ wasm = 'app.wasm', canvas = document.getElementById('canvas'), assets = null, args = [], env = {} } = {}) {
  const hooks = {};
  globalThis.rlgo = hooks;
  const rl = await createRaylib({ canvas });

  if (assets) {
    const res = await fetch(assets);
    if (res.ok) {
      let body = res.body;
      if (assets.endsWith('.gz')) body = body.pipeThrough(new DecompressionStream('gzip'));
      untar(rl.FS, new Uint8Array(await new Response(body).arrayBuffer()));
    }
  }

  globalThis.fs = goFS(rl.FS);
  const go = new Go();
  go.argv = ['app', ...args];
  go.env = env;

  let goBuf, goView, cBuf, cView;
  const goMem = () => {
    const b = go._inst.exports.mem.buffer;
    if (b !== goBuf) { goBuf = b; goView = new Uint8Array(b); }
    return goView;
  };
  const cMem = () => {
    const b = rl.wasmExports.memory.buffer;
    if (b !== cBuf) { cBuf = b; cView = new Uint8Array(b); }
    return cView;
  };
  go.importObject.raylib = rl.wasmExports;
  go.importObject.rlgo = {
    copyIn: (dst, src, n) => cMem().set(goMem().subarray(src, src + n), dst),
    copyOut: (dst, src, n) => goMem().set(cMem().subarray(src, src + n), dst),
    cString: (dst, src, n) => { const c = cMem(); c.set(goMem().subarray(src, src + n), dst); c[dst + n] = 0; },
  };

  const source = await fetch(wasm);
  const { instance } = await WebAssembly.instantiateStreaming(source, go.importObject);
  await go.run(instance);
}
