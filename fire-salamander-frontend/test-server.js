// Simple HTTP server to test connectivity
const http = require('http');

const server = http.createServer((req, res) => {
  res.writeHead(200, {'Content-Type': 'text/html'});
  res.end(`
    <html>
      <body>
        <h1>Fire Salamander - Test Server</h1>
        <p>If you see this, the server is working!</p>
        <p>Now try accessing the Next.js app at:</p>
        <ul>
          <li><a href="http://127.0.0.1:3000">http://127.0.0.1:3000</a></li>
          <li><a href="http://localhost:3000">http://localhost:3000</a></li>
        </ul>
      </body>
    </html>
  `);
});

server.listen(8080, '127.0.0.1', () => {
  console.log('Test server running at http://127.0.0.1:8080/');
});