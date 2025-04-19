import http from 'k6/http';
import { check } from 'k6';

// Load the image during init stage
const imageData = open('assets/test.png', 'b');

export const options = {
  vus: 16,
  iterations: 16,
//   duration: '30s',
};

export default function () {
  const url = 'http://localhost:8080/compress';

  const payload = {
    image: http.file(imageData, 'test.png'),
    quality: '80',
    format: 'webp',
    algorithm: 'scale',
  };

  const res = http.post(url, payload);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'received image': (r) => r.headers['Content-Type'].includes('image/'),
  });
}
