# ShodSplit

ShodSplit is a blazing-fast subdomain extraction tool written in Go, using Shodan's DNS API.  
It fetches subdomains, filters them, and sorts them by TLD for cleaner and more structured recon results.

## ✨ Features

- Uses Shodan DNS API to fetch subdomains
- Sorts subdomains alphabetically by TLD
- Simple CLI usage with clean output
- Fast and lightweight (pure Go)

## 🚀 Usage

```bash
go run main.go -k YOUR_SHODAN_API_KEY -d target.com
```
<img width="1080" height="431" alt="image" src="https://github.com/user-attachments/assets/1f26086a-1f19-4871-a467-ba4f1a49c5f3" />
