const sharp = require('sharp');
const fs = require('fs');
const path = require('path');
const os = require('os');

const imagePath = path.resolve(__dirname, 'test.png');
const outputDir = path.resolve(__dirname, 'out');
const iterations = 50;
const concurrency = os.cpus().length;

if (!fs.existsSync(outputDir)) {
  fs.mkdirSync(outputDir);
}

const inputBuffer = fs.readFileSync(imagePath);

async function compress(index) {
  const start = Date.now();

  const outputBuffer = await sharp(inputBuffer)
    .webp({ quality: 80 })
    .toBuffer();

  const duration = Date.now() - start;

  const outPath = path.join(outputDir, `out-${index}.webp`);
  fs.writeFileSync(outPath, outputBuffer);

  return {
    index,
    duration,
    originalSize: inputBuffer.length,
    compressedSize: outputBuffer.length,
  };
}

async function runBatch() {
  console.log(`Compressing ${iterations} images using sharp (${concurrency} parallel workers)`);

  const results = [];
  const queue = Array.from({ length: iterations }, (_, i) => i);

  while (queue.length > 0) {
    const chunk = queue.splice(0, concurrency);
    const promises = chunk.map(compress);
    const res = await Promise.all(promises);
    results.push(...res);
  }

  const totalTime = results.reduce((sum, r) => sum + r.duration, 0);
  const min = Math.min(...results.map(r => r.duration));
  const max = Math.max(...results.map(r => r.duration));
  const avg = totalTime / results.length;

  console.log(`\nCompleted ${results.length} compressions`);
  console.log(`Avg: ${avg.toFixed(2)}ms | Min: ${min}ms | Max: ${max}ms`);
  console.log(`Original size: ${results[0].originalSize} bytes`);
  console.log(`Avg compressed size: ${Math.round(results.reduce((s, r) => s + r.compressedSize, 0) / results.length)} bytes`);
  console.log(`Avg ratio: ${(results.reduce((s, r) => s + r.compressedSize / r.originalSize, 0) / results.length).toFixed(2)}x`);
}

runBatch().catch(console.error);
