// Code generated by "esc -o static/static.go -pkg static static/ui"; DO NOT EDIT.

package static

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/static/ui/confirm.ui": {
		local:   "static/ui/confirm.ui",
		size:    3295,
		modtime: 1544925482,
		compressed: `
H4sIAAAAAAAC/9RXT2/bPgy991Pop+sPTtL0soPtYn/QHjZgl+5s0DLjcFFET6LTZp9+cLytyaLUcVdg
6CVIovdEPvKRhtPrh7VVG/SB2GX6cjLTCp3hilyd6S93N8kbfZ1fpP8libpFhx4EK3VPslS1hQrV1WQ+
n1yqJMkvUnKCfgEG8wulUo/fWvIYlKUy07Ws/tePga4m85me7nBcfkUjylgIIdO3svpAYLnWiqpMG3YL
8mvdIZVKG88NetkqB2vMtAFXLNi0Qec3YAOm01+AON5joO9QWjwTv+YKrM7vfDsIvSdX8X3RcCAhdjo3
2FUjYZc04NHJ0AUVBvG8LbraFj3lvMiybbBYUgcXZluCHw5lUZ4qg1mSrfrvHduCwSXbCv30J2C6h+jR
atd8BzbZ/cz0puQH/fuOoya/2zsd39gYhz2hE+irv0EvZMCeQwwNGHK1zudRdFwemC5QAR5hT0dUaSvC
7lDvczXHeBa23EoRZNs1FF11knjQ1qGED+avMOAMWv0nOZZNiVYr8eCChZ3LMr3FoPP3uytOJRe7q/vU
+WEKY/gbCrRzeWyMzmzFWKpHg7TBUFS4gNbK0A3ptK//0f8NmBW5ejgiPjTgqvGZLsja8azHBTd7SlQ0
+4O18RKu5NVfOfLzx9fnpidXw2u20+VL2CkmPy79lOxRq7fXPIoyOD8RsUdCj2fmeF4+9bbfHxdbDjyC
zrL4Sz23lmCpdjoPAl7OJa3B1+QKiwvR+Xw2kuapXj6HJ9w8g1WyCK/HEE+vqrce1ZZbFVqP16eM8w+8
P8Yog7M+bP1DjXuHjwfpdO815EcAAAD//yoQrQnfDAAA
`,
	},

	"/static/ui/sign.ui": {
		local:   "static/ui/sign.ui",
		size:    10587,
		modtime: 1545056643,
		compressed: `
H4sIAAAAAAAC/+xaX2/bthd9z6e4Pz79ttVObK/bUNgq1nXtBgzYsHXPAkVdS1xoUiOvnGiffpCc1E4s
W6KkeGmTtzbkIXn/nMMjS/PX1ysFa7ROGr1gk/EFA9TCxFInC/bnh3ej79jr4Gz+v9EI3qNGywljuJKU
QqJ4jDAbT6fjCYxGwdlcakK75AKDM4C5xb9zadGBktGCJXT5FdtuNBtPL9h5Nc9Ef6EgEIo7t2Dv6fKt
5MokDGS8YE4mmpXTAOaZNRlaKkDzFS6Y4DpcGpE7FrzjyuH8/HZC/XyLTv7DI4Ut519JHZurMDNOkjSa
BQLL+EZGjzJuUVPTAjE6sqYIy2yFGwgLPti8cWcqMgxTWU6Pq1zsAUQqVbz5dwlXXGBqVIz2/GbC+c6M
zWyoqqO5GlX/XbB1ZK7ZxzX2qvBmZ9Q/+XWYFbeJ1KHCJbFgcuEBsTJJfTFkMk9EZIjMqi3IWIma+KY1
1mhJCq7aAF3GhdQJC6a1s+uLxUW5Ucgt8p2q1NYtJzL6bvW6VrAOp3hhcgodFSWTUMcHgXeatOnAW7qH
Lo9Wkth9ZN1RIlQMyHLtFKeS2wtWoGPBH9USh05Wt9ZaOllpQx1BW+bRF2pRoFyjC2Nc8lxRQyVKWm+y
t/f3jItLqZPmLfE64zr2P+pSKuWP2mrn5FhQtae/I2GHw68P/VDYXp2+idkLsg34oh5SE+xeoPvEaVLn
zg09lCi0kPau8t5N4jvLfC+p95G9XyrxOo0wNcjKPjblqjJejrj1ktF2lWlRnakX+MhNwJdekTupE4Wh
khrDlYnb5J0TWRnlhG5/cHf4VlaMplGMTjBYc5Xjgv2Eao1lU8EVllRYTGffwssbf3xfLA/vNtzt4N0v
fa+HiyGuh2futdQ3Pwb9xp3LUssdvjqFMzl5701O23s/arLFxu1mHzPLHqdBTbnrCpWdkVWgUkkqOjSD
1Dd5CkXKLQu+9OqJ7XN0SHhNtXT4VSOUhbsyNgYyYHOFQCmugCvld9QspzDLbWbKBrhd87Pk2PRZ34e3
SLOBBP57SIyJYatG4FKTqxhSvkbgBAq5I5h8AyWnuCC0DriOIUKI5XIpRa6opEKSo3Pw/1gmktwLcMUq
Msq9ACQx/mL8WTb27LSNXT72bX8oEUYvpV2FXvfIAK3bzaN7PUUdToof43uxvm/6erm73g7PTwR+2PTS
jgq8Or74IdIeJ24v8nYncGcHeITIB8ns27s7jrAVox+yp7vAvQxbb9PW2bi9NXmkEESKVTGfTHdP/9vu
flLK3MLU+cny7ziiIqseMHYNGhm4kQooTG5BpEYKHD+Znp4N3dOfsAX9+rQW9Ge9NG+43bGhaK2xYVQq
98O6T3SOJxiWhGBBtWuzfWz9JtX/jeowQvSoWNkJ2upn3H5XzV4VhdGEmrzK2LqAn9A1dPAiPnYZ36Mu
O4TvmZkhsuN3W1YjzUseo1wz7fq8Vu9/J/YiYAMJjxLxIZSqC7+GsA8vn+1Dm1w88i8zfD5saHzYb/4w
426MO4Pbgfn5zieQ/wYAAP//ESw79FspAAA=
`,
	},

	"/static/ui/sign.ui~": {
		local:   "static/ui/sign.ui~",
		size:    10636,
		modtime: 1545056643,
		compressed: `
H4sIAAAAAAAC/+xa3XLbNhO991Psh6uvbSRbUtN2MhIzTdOknelMO216zQHBFYkaAlhgKZt9+g4pO6It
SiRIWXVi3yUmDn5295w9EDl/fb1SsEbrpNELNhlfMEAtTCx1smB/fng3+o69Ds7m/xuN4D1qtJwwhitJ
KSSKxwiz8XQ6nsBoFJzNpSa0Sy4wOAOYW/w7lxYdKBktWEKXX7HtQrPx9IKdV+NM9BcKAqG4cwv2ni7f
Sq5MwkDGC+Zkolk5DGCeWZOhpQI0X+GCCa7DpRG5Y8E7rhzOz28HNI9PuQtjXPJcEQs+2LwVYNHJf3ik
sOMCV1LH5irMjJMkjWaBwDIgI6NHGbeoqW2CGB1ZU4RleMMNpNtOqcgwTGU5PK6CtwMQqVTx5t8lXHGB
qVEx2vObAee1EZvRUKVTczWq/rtg68hcs49z7KTtTe2pf7aaMCtuE6lDhUtiweTCA2JlkvpiyGSeiMgQ
mVVXkLESNfFNaazRkhRcdQG6jAupExZMG0c3J4uLcqGQW+S1rDTmLScy+m72+mawCad4YXIKHRUlk1DH
e4F3irRtw1t9CF0erSSx+8imrUSoGJDl2ilOJbcXrEDHgj+qKfbtrGmutXSy0oYmgnaMoy/UokC5xpqO
HcxESetN9Hb+nnFxKXXSviReZ1zH/ltdSqX8UVvtnBw6VOPu70jY/uM3H33fsb0qfXNmL8j2wBfNkIbD
7hx0lzht6ty7oI8lCh2kva+895P43jI/SOp9ZO+XSrxOI0wtstJkrlTl1Bxx6yWj3TLTITtTL/CBTsCX
Xid3UicKQyU1hisTd4k7J7Iyygnd7sP641tZMZpGMTrBYM1Vjgv2E6o1lkUFV1hSYTGdfQsvbwz1fbHc
v9rxuoN3vQxtDxfHaA/P3Ouob34M+o07l6WWO3x1Cmdy8tqbnLb2ftRki43bzT5Glj1Og1resXtCZW9k
dVCpJBU9ikHqmziFIuWWBV961cT2Hh0SXlMjHX7VCGXiroyNgQzYXCFQiivgSvltNcspzHKbmbIAbuf8
LDk2fdb341uk2ZEE/ntIjIlhq0bgUpOrGFK+RuAECrkjmHwDJae4ILQOuI4hQojlcilFrqikQpKjc/D/
WCaS3AtwxSoyyr0AJDH+YvxZFvbstIVdXvu2P5QIo5fSrkKvPnKE0u3n0b1uUfuD4sf4QawfGr5B7m6w
w/MTgR82tVRTgVeHJ99H2sPEHUTe/gTu7QAPEHkvmX1rt+YIOzH6IWu6D9zLsA02bb2N21uTRwpBpFgl
88lU9/S/re4npcwdTJ2fLP+OIyqy6oJRN2hk4EYqoDC5BZEaKXD8ZGp6duya/oQt6NentaA/66V5w23N
hqK1xoZRqdwP6z7ROZ5gWBKCBdWq7fax85tU/zeqxxGiR8XKXtBOP+MOazU7WRRGE2rySmPnBH5CbWhv
Iz7UjO9Rl+3DD4zMMaLj1y2rJ+1THqJcO+2GvFYf3hMHEbCFhAeJ+BBK1Ydfx7APL5/tQ5dYPPIvM3w+
bGi97Ld/mHH3jLWH2wfz89o3k/8GAAD//8BbM5+MKQAA
`,
	},

	"/static/ui/tag.ui": {
		local:   "static/ui/tag.ui",
		size:    6599,
		modtime: 1544925387,
		compressed: `
H4sIAAAAAAAC/+xZTW/jNhC951ewvBaOZeVSFJIW2La7KFD00u1ZGFFjiQ1NquTYjvvrC8nbtb2RRNHy
BgnQWxLNI+fjzRsOkrx72ii2Q+uk0Slf3UecoRamlLpK+Z+fPix+4O+yu+S7xYJ9RI0WCEu2l1SzSkGJ
7OE+ju9XbLHI7hKpCe0aBGZ3jCUW/95Ki44pWaS8osfv+emih/s44svOzhR/oSAmFDiX8o/0+LMEZSrO
ZJlygoq3VowljTUNWjowDRtMuQCdr43YOp59AOUwWf5n0G9v0cl/oFA40X5jSlA8+2S3XtO91KXZ541x
kqTRPBPYZmJh9KIBi5p8B5S4hq2ifC9Lqnn2EEV+hCNrDnlbifx4yTRf6dBgXsvWvOzy7L9JIY3lTdRS
lcefW7QCgbVRJdrlZ4PlmcXRmnVM0aAW3a8p3xXmiX854xkj3p99DWdCb3XBVlLnCtfEs1UUALGyqkMx
ZJpARGGIzGYqyFiJmuBIvh1akgLUFKBrQEhd8Szute4vFoj2ohwswllVeuu2JTL6snrXVrAPp+BgtpQ7
OrT0RF0OAi9I6nP4i/TkArRAxb8G9nlSoOKMLGinoOuXlB/Q8eyn7oghx/rO2kknu37r6+eJaQyFWhQo
d+jyz1LkOyFZHnP37O8NiEepK/+N+NSALsM9XUulwlEnbY7Ggur1/kLA5jLKbYuNpFmM+qM74n9GvRZG
rW7BqL7w+0MfCjtIOY8xB0G8LdQT7LNAn7eNb9pfzedbDZmvJnIcTQVOnsohevJbJwsv0/KjaerD1qBk
pXnmCGyQQJ2/k+LoCuT4eylcY3/HPSOoOrsfX0K0gnP9xubgL5rs4TQGW2/46xxcNbhp0Ddc/NXLFv9X
vTbvwZ7Kj9Yamxdg+bcVJP/+NYhE56DCvN2Xedb5OwoPW5XCV6ZbZGeYsuO0nUXdK18fswRshMeDXB6p
ojCaUFNQGScXcLLEza/+jIE93O2+p8plz/Mh+MzE3CI5YS+E7ov/yLGO83fdnDVn3uyY3X+eHhztw28h
VNe01/RUDacpvrVMveFnR/xmN+WQRdP7zPIvypcxnn08fUiWZ/9u+TcAAP//rydwzscZAAA=
`,
	},

	"/static/ui/vault.ui": {
		local:   "static/ui/vault.ui",
		size:    25295,
		modtime: 1544925368,
		compressed: `
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
`,
	},

	"/": {
		isDir: true,
		local: "",
	},

	"/static": {
		isDir: true,
		local: "static",
	},

	"/static/ui": {
		isDir: true,
		local: "static/ui",
	},
}
