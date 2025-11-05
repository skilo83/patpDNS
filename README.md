# patpDNS: Phonetic Names for IPs That Stick

**patpDNS** is an open-source Go app that swaps forgettable IP addresses for slick, memorable patp-style names (inspired by Urbit's phonetic system). Turn `192.168.1.1` into `~sampel-palnet` or a beastly IPv6 into `~ronfyl-mitdeg-nalful-wicpub-siglet-hidwes-fidper-lissut`—easy to recall, type, and share, with zero-loss reversibility baked in.

## Why patpDNS?
Ditch the numeric drudgery of networking for something that *sounds* right:
- **Direct IP-to-Name Encoding**: Syllables map byte-for-byte to IPs—no hashes, no collisions, no bloat.
- **Local-First Freedom**: Your own DNS resolver for LANs, P2P, or private clouds, forwarding everything else upstream.
- **Scales Gracefully**: 2-syllable patps for IPv4, 8-syllable chains for IPv6.

## Key Features
- **CLI Magic**:
  - `generate <ip>`: IP → patp name.
  - `resolve <patp>`: Name → IP (round-trip perfect).
- **DNS Server Mode**:
  - Handles `*.patp` queries (A/AAAA records) locally.
  - Forwards non-patp domains to upstream (e.g., 8.8.8.8).
  - UDP/TCP listener—port 53 (sudo) or 5353 for dev.
- **IPv4 & IPv6 Native**: Fits the bits exactly.
- **Minimalist**: Stdlib + `miekg/dns`; build and ship anywhere.

## Quick Start
1. **Setup**: `go mod init patpdns && go get github.com/miekg/dns`
2. **Name an IP**:  
   ```
   go run main.go generate 192.168.0.1
   # ~sampel-palnet
   ```
3. **Unpack a Name**:  
   ```
   go run main.go resolve ~sampel-palnet
   # 192.168.0.1
   ```
4. **Spin Up DNS**:  
   ```
   sudo go run main.go server -listen=:53
   ```
   Query: `dig @localhost sampel-palnet.patp` → IP magic.

## Use Cases
- **Home Networks**: Name your smart fridge `~dozbud-hidlet`.
- **Dev & Testing**: Mock domains without `/etc/hosts` hacks.
- **P2P Vibes**: Discovery names for meshes or ad-hoc nets.
- **Phonetic Nerdery**: Because why not make tech poetic?

Open-source, forkable, and fun.
