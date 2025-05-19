# potato-quality-image-compressor

## Update:
It turns out big binary blobs aren't really that great to send over gRPC compared to plain old HTTP/REST. And if you end up using libvips, like... TypeScript is fine. HTTP is fine. npm Sharp is fine.

The generated code business is so annoying, too.

Oh, but Prometheus/Grafana are kinda cool. I'll maybe take some screenshots later, but honestly, the shell scripts basically told everything. Like... the CPU usage on my machine is just... well, higher is honestly better; it's not exactly a closed environment. Higher CPU usage, higher throughput. That's it.

---
I had to install libwebp to make the webp.Encode and other stuff compileable.

I will be revisiting this in the future for other compression algorithms.

This takes in an http request and sends back the image, or a .zip file.
