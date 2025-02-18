# CLI Tool

This CLI tool provides various helper functions to perform common tasks such as JSON encoding/decoding, Base64 encoding/decoding, string compression/decompression, and UUID generation.

## Installation

To install the CLI tool, use the following command:

```sh
go install github.com/sondy91/quic@latest
```

## Usage

Below are the commands available in this CLI tool:

### JSON Encode/Decode

- **Encode a string to JSON:**
  ```sh
  quic json encode "your string"
  ```

- **Decode a JSON string:**
  ```sh
  quic json decode '{"key":"value"}'
  ```

### Base64 Encode/Decode

- **Encode a string to Base64:**
  ```sh
  quic base64 encode "your string"
  ```

- **Decode a Base64 string:**
  ```sh
  quic base64 decode "eW91ciBzdHJpbmc="
  ```

### Compress/Decompress

- **Compress a string:**
  ```sh
  quic compress "your string"
  ```

- **Decompress a string:**
  ```sh
  quic decompress "compressed string"
  ```

### Generate UUIDs

- **Generate a specified number of UUIDs:**
  ```sh
  quic uuid generate 5
  ```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

## License

This project is licensed under the MIT License.
