# nanote-server
A simple, stateless, straight-to-the-point media server  
## Capabilities
- Recursively digest media collections with thousands of files  
- Extract metadata and album art, either embedded in the file or included as `cover.jpg` in the file's directory  
- Serve metadata library, audio data, and album art
- HTTP basic authentication  
- Multiple users/libraries behind separate passwords  
- Cache metadata library, with an API to trigger rebuilding  
## Installation
1. [Ensure you have Go 1.16 or later installed](https://go.dev/dl/)  
2. Clone or [download](https://github.com/cyfinfaza/nanote-server/archive/refs/heads/master.zip) this repository  
3. Run `go build` to build the binary for your platform
4. Modify the sample `config.yml` with your configuration
5. Run the binary you built, ensuring the config file is in the running directory
## API
### Authentication not required
- `GET /users`: Returns a list of users
- `GET /userImg/<user>`: Returns the user's profile image or `404`
### Authentication required
For the following endpoints, an HTTP basic authentication header is required. If authentication fails, the server will return a `401 Unauthorized` response.
- `GET /test`: Returns `200`
- `GET /coverImage/<filepath>`: Returns embedded album art for a file in the mediaRoot or `404`
- `GET /content/<filepath>`: Returns content of a file in the mediaRoot or `404`
- `GET /library`: Builds (if cache is disabled) and returns a list of all song metadata (title, album, artist, genre, year, media URL, cover URL)
- `GET /rebuildLibrary`: Rebuilds (if cache is enabled) the user's library and returns `200`, `503`, or `400`