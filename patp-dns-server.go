package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

var prefixes = []string{
	"doz", "mar", "bin", "wan", "sam", "lit", "sig", "hid", "fid", "lis", "sog", "dir", "wac", "sab", "wis", "sib",
	"rig", "sol", "dop", "mod", "fog", "lid", "hop", "dar", "dor", "lor", "hod", "fol", "rin", "tog", "sil", "mir",
	"hol", "pas", "lac", "rov", "liv", "dal", "sat", "lib", "tab", "han", "tic", "pid", "tor", "bol", "fos", "dot",
	"los", "dil", "for", "pil", "ram", "tir", "win", "tad", "bic", "dif", "roc", "wid", "bis", "das", "mid", "lop",
	"ril", "nar", "dap", "mol", "san", "loc", "nov", "sit", "nid", "tip", "sic", "rop", "wit", "nat", "pan", "min",
	"rit", "pod", "mot", "tam", "tol", "sav", "pos", "nap", "nop", "som", "fin", "fon", "ban", "mor", "wor", "sip",
	"ron", "nor", "bot", "wic", "soc", "wat", "dol", "mag", "pic", "dav", "bid", "bal", "tim", "tas", "mal", "lig",
	"siv", "tag", "pad", "sal", "div", "dac", "tan", "sid", "fab", "tar", "mon", "ran", "nis", "wol", "mis", "pal",
	"las", "dis", "map", "rab", "tob", "rol", "lat", "lon", "nod", "nav", "fig", "nom", "nib", "pag", "sop", "ral",
	"bil", "had", "doc", "rid", "moc", "pac", "rav", "rip", "fal", "tod", "til", "tin", "hap", "mic", "fan", "pat",
	"tac", "lab", "mog", "sim", "son", "pin", "lom", "ric", "tap", "fir", "has", "bos", "bat", "poc", "hac", "tid",
	"hav", "sap", "lin", "dib", "hos", "dab", "bit", "bar", "rac", "par", "lod", "dos", "bor", "toc", "hil", "mac",
	"tom", "dig", "fil", "fas", "mit", "hob", "har", "mig", "hin", "rad", "mas", "hal", "rag", "lag", "fad", "top",
	"mop", "hab", "nil", "nos", "mil", "fop", "fam", "dat", "nol", "din", "hat", "nac", "ris", "fot", "rib", "hoc",
	"nim", "lar", "fit", "wal", "rap", "sar", "nal", "mos", "lan", "don", "dan", "lad", "dov", "riv", "bac", "pol",
	"lap", "tal", "pit", "nam", "bon", "ros", "ton", "fod", "pon", "sov", "noc", "sor", "lav", "mat", "mip", "fip",
}

var suffixes = []string{
	"zod", "nec", "bud", "wes", "sev", "per", "sut", "let", "ful", "pen", "syt", "dur", "wep", "ser", "wyl", "sun",
	"ryp", "syx", "dyr", "nup", "heb", "peg", "lup", "dep", "dys", "put", "lug", "hec", "ryt", "tyv", "syd", "nex",
	"lun", "mep", "lut", "sep", "pes", "del", "sul", "ped", "tem", "led", "tul", "met", "wen", "byn", "hex", "feb",
	"pyl", "dul", "het", "mev", "rut", "tyl", "wyd", "tep", "bes", "dex", "sef", "wyc", "bur", "der", "nep", "pur",
	"rys", "reb", "den", "nut", "sub", "pet", "rul", "syn", "reg", "tyd", "sup", "sem", "wyn", "rec", "meg", "net",
	"sec", "mul", "nym", "tev", "web", "sum", "mut", "nyx", "rex", "teb", "fus", "hep", "ben", "mus", "wyx", "sym",
	"sel", "ruc", "dec", "wex", "syr", "wet", "dyl", "myn", "mes", "det", "bet", "bel", "tux", "tug", "myr", "pel",
	"syp", "ter", "meb", "set", "dut", "deg", "tex", "sur", "fel", "tud", "nux", "rux", "ren", "wyt", "nub", "med",
	"lyt", "dus", "neb", "rum", "tyn", "seg", "lyx", "pun", "res", "red", "fun", "rev", "ref", "mec", "ted", "rus",
	"bex", "leb", "dux", "ryn", "num", "pyx", "ryg", "ryx", "fep", "tyr", "tus", "tyc", "leg", "nem", "fer", "mer",
	"ten", "lus", "nus", "syl", "tec", "mex", "pub", "rym", "tuc", "fyl", "lep", "deb", "ber", "mug", "hut", "tun",
	"byl", "sud", "pem", "dev", "lur", "def", "bus", "bep", "run", "mel", "pex", "dyt", "byt", "typ", "lev", "myl",
	"wed", "duc", "fur", "fex", "nul", "luc", "len", "ner", "lex", "rup", "ned", "lec", "ryd", "lyd", "fen", "wel",
	"nyd", "hus", "rel", "rud", "nes", "hes", "fet", "des", "ret", "dun", "ler", "nyr", "seb", "hul", "ryl", "lud",
	"rem", "lys", "fyn", "wer", "ryc", "sug", "nys", "nyl", "lyn", "dyn", "dem", "lux", "fed", "sed", "bec", "mun",
	"lyr", "tes", "mud", "nyt", "byr", "sen", "weg", "fyr", "mur", "tel", "rep", "teg", "pec", "nel", "nev", "fes",
}

var prefixToIndex = make(map[string]uint8)
var suffixToIndex = make(map[string]uint8)

func init() {
	for i, p := range prefixes {
		prefixToIndex[p] = uint8(i)
	}
	for i, s := range suffixes {
		suffixToIndex[s] = uint8(i)
	}
}

func generatePatp(ipStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", fmt.Errorf("invalid IP address")
	}

	if v4 := ip.To4(); v4 != nil {
		// IPv4: 4 bytes, 2 heads (4 syllables)
		bytes := v4
		head0 := prefixes[bytes[3]] + suffixes[bytes[2]]
		head1 := prefixes[bytes[1]] + suffixes[bytes[0]]
		return "~" + head0 + "-" + head1, nil
	} else if v16 := ip.To16(); v16 != nil {
		// IPv6: 16 bytes, 8 heads (16 syllables)
		bytes := v16
		var heads [8]string
		for i := 0; i < 8; i++ {
			preIdx := bytes[15-2*i]
			sufIdx := bytes[14-2*i]
			heads[i] = prefixes[preIdx] + suffixes[sufIdx]
		}
		return "~" + strings.Join(heads[:], "-"), nil
	}
	return "", fmt.Errorf("unsupported IP version")
}

func resolvePatp(name string) (string, error) {
	name = strings.TrimPrefix(name, "~")
	parts := strings.Split(name, "-")
	numHeads := len(parts)
	if numHeads == 2 {
		// IPv4: 2 heads
		if len(parts[0]) != 6 || len(parts[1]) != 6 {
			return "", fmt.Errorf("invalid head lengths: expected 6 letters each")
		}
		clan := parts[0]
		titl := parts[1]
		pre0 := clan[:3]
		suf1 := clan[3:]
		pre2 := titl[:3]
		suf3 := titl[3:]

		b0, ok := prefixToIndex[pre0]
		if !ok {
			return "", fmt.Errorf("invalid prefix: %s", pre0)
		}
		b1, ok := suffixToIndex[suf1]
		if !ok {
			return "", fmt.Errorf("invalid suffix: %s", suf1)
		}
		b2, ok := prefixToIndex[pre2]
		if !ok {
			return "", fmt.Errorf("invalid prefix: %s", pre2)
		}
		b3, ok := suffixToIndex[suf3]
		if !ok {
			return "", fmt.Errorf("invalid suffix: %s", suf3)
		}

		ipBytes := []byte{byte(b3), byte(b2), byte(b1), byte(b0)}
		return net.IP(ipBytes).String(), nil
	} else if numHeads == 8 {
		// IPv6: 8 heads
		var ipBytes [16]byte
		for i := 0; i < 8; i++ {
			head := parts[i]
			if len(head) != 6 {
				return "", fmt.Errorf("invalid head length: %s", head)
			}
			pre := head[:3]
			suf := head[3:]
			preIdx, ok := prefixToIndex[pre]
			if !ok {
				return "", fmt.Errorf("invalid prefix: %s", pre)
			}
			sufIdx, ok := suffixToIndex[suf]
			if !ok {
				return "", fmt.Errorf("invalid suffix: %s", suf)
			}
			ipBytes[15-2*i] = preIdx
			ipBytes[14-2*i] = sufIdx
		}
		return net.IP(ipBytes[:]).String(), nil
	}
	return "", fmt.Errorf("invalid patp format: expected 2 or 8 hyphen-separated heads")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <program> <command> [args]")
		fmt.Println("Commands:")
		fmt.Println("  generate <ip>   Generate patp (IPv4) or extended patp (IPv6, 8 heads) from IP")
		fmt.Println("  resolve <patp/extended-patp>  Resolve to IP")
		fmt.Println("  server [-listen=:53] [-upstream=8.8.8.8:53]  Run DNS server (use port 5353 for testing, requires sudo for 53)")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "generate":
		generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
		generateCmd.Parse(os.Args[2:])
		if generateCmd.NArg() != 1 {
			fmt.Println("Usage: <program> generate <ip>")
			os.Exit(1)
		}
		patp, err := generatePatp(generateCmd.Arg(0))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(patp)
	case "resolve":
		resolveCmd := flag.NewFlagSet("resolve", flag.ExitOnError)
		resolveCmd.Parse(os.Args[2:])
		if resolveCmd.NArg() != 1 {
			fmt.Println("Usage: <program> resolve <patp/extended-patp>")
			os.Exit(1)
		}
		ip, err := resolvePatp(resolveCmd.Arg(0))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(ip)
	case "server":
		listen := flag.String("listen", ":53", "address to listen on (e.g., :5353 for testing)")
		upstream := flag.String("upstream", "8.8.8.8:53", "upstream DNS server")
		flag.Parse()

		dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
			m := new(dns.Msg)
			m.SetReply(r)
			defer w.WriteMsg(m)

			if len(r.Question) != 1 {
				return
			}
			q := r.Question[0]
			name := strings.ToLower(strings.TrimSuffix(string(q.Name), "."))

			if !strings.HasSuffix(name, ".urbit") {
				// Forward to upstream
				c := new(dns.Client)
				c.Net = "udp"
				c.Timeout = 5 * time.Second
				resp, _, err := c.Exchange(r, *upstream)
				if err != nil {
					m.Rcode = dns.RcodeServerFailure
					return
				}
				m = resp
				return
			}

			// Local Urbit resolution
			if q.Qtype != dns.TypeA && q.Qtype != dns.TypeAAAA {
				m.Rcode = dns.RcodeRefused
				return
			}

			patpName := "~" + strings.TrimSuffix(name, ".urbit")
			ipStr, err := resolvePatp(patpName)
			if err != nil {
				m.Rcode = dns.RcodeNameError
				return
			}

			ip := net.ParseIP(ipStr)
			if ip == nil {
				m.Rcode = dns.RcodeServerFailure
				return
			}

			isIPv4 := ip.To4() != nil
			isIPv6 := ip.To16() != nil && !isIPv4

			if q.Qtype == dns.TypeA && isIPv4 {
				rr := &dns.A{
					Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
					A:   ip.To4(),
				}
				m.Answer = append(m.Answer, rr)
			} else if q.Qtype == dns.TypeAAAA && isIPv6 {
				rr := &dns.AAAA{
					Hdr:  dns.RR_Header{Name: q.Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 3600},
					AAAA: ip.To16(),
				}
				m.Answer = append(m.Answer, rr)
			} else {
				m.Rcode = dns.RcodeNameError
				return
			}
		})

		udpServer := &dns.Server{Addr: *listen, Net: "udp"}
		tcpServer := &dns.Server{Addr: *listen, Net: "tcp"}

		go func() {
			if err := tcpServer.ListenAndServe(); err != nil {
				fmt.Printf("TCP server failed: %v\n", err)
			}
		}()

		if err := udpServer.ListenAndServe(); err != nil {
			fmt.Printf("UDP server failed: %v\n", err)
		}

		select {}
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
