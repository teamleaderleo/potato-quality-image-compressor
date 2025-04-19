import grpc from 'k6/net/grpc';
import { check } from 'k6';
import encoding from 'k6/encoding';

const client = new grpc.Client();
client.load(['../proto'], 'compression_service.proto');

const imageData = encoding.b64encode(open('assets/test.png', 'b'));

export const options = {
  vus: 10,
  iterations: 10,
//   duration: '30s',
};

export default function () {
  client.connect('localhost:9000', { plaintext: true });

  const response = client.invoke('compression.ImageCompressionService/CompressImage', {
    filename: 'test.png',
    format: 'webp',
    strategy: 'scale',
    quality: 80,
    imageData: imageData,
  });

  console.log(JSON.stringify(response.message, null, 2));


  check(response, {
    'no gRPC error': (r) => !r.error,
    'got image back': (r) => r.message && r.message.imageData && r.message.imageData.length > 0,
    'compression ratio present': (r) => r.message.compressionRatio > 0,
  });

  client.close();
}
