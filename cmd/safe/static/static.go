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

	"/static/ui/confirm.ui~": {
		local:   "static/ui/confirm.ui~",
		size:    3411,
		modtime: 1544925482,
		compressed: `
H4sIAAAAAAAC/9RXwW7bOBC95yu4c13Isp3LHiQFu1skBVqgl/QsjKixzJomVXKkRP36QlaT2DEdWWmA
IhfDNt8bzpt5M7KTq/utFi05r6xJYTGbgyAjbalMlcLX2+voH7jKLpK/okjckCGHTKW4U7wWlcaSxOVs
uZwtRBRlF4kyTG6FkrILIRJH3xvlyAutihQq3vwNTxddzpZziHc4W3wjyUJq9D6FG958UKhtBUKVKUhr
VsptoUcKkdTO1uS4Ewa3lIJEk6+sbDxk16g9JfEDIIx35NUPLDSdid/aEjVkt64Zhd4pU9q7vLZesbIG
Mkl9NSJrohodGR4LUJJnZ7u8r20+UM67mbua8rXq4WytLtCNX6WJXyqDXCtdDu97tkZJa6tLcvEvQLyH
GNBi13yDOtp9TKEt7D08xjhq8n97p9MbG+JYp8gwDtVvybGSqM8h+hqlMhVkyyA6LA9lf1GOjnBPR1Bp
w2zNod7Xag7xNHa24dxz1zeUTHmSeNDWsYQP5i+XaCRpeE4OZVOQBsEOjde4c1kKHXnI/t+FOJVcKFb/
CtlhClP4rfJq5/LQGJ3ZiqlUR5JUSz4vaYWN5rEISTzU/+j7GuVGmWr8Rrqv0ZTTM10praeznhbc/CVR
wewP1sZbuNJufsuRXz69Pze9uBres50Wb2GnkPyw9FOyJ63eQfMkyuj8BMQeCT2emeN5+TzYfn9cdDHy
CDrL4m/13FqjVpWBzDM6Ppe0RVcpk2taMWTL+USaU9X6NTy29StYhWW22ynE06vqX0eis43wjaOrk+GQ
2amiYfJHI/R49GBdazgqyUsQLeqGUvhIuqX+R5NYXEL8fNTCof/IvE0x5+h+GR+3Q417h08HSbz31+dn
AAAA///TF9qVUw0AAA==
`,
	},

	"/static/ui/sign.ui": {
		local:   "static/ui/sign.ui",
		size:    10583,
		modtime: 1544925503,
		compressed: `
H4sIAAAAAAAC/+xaQXPbNhO9+1fgw+lrG8mW1LSdjMRM0zRpZzrTTpueOUtwRaKGABZYymZ/fYeUHckW
JRIkrTqxb4mBB2Cx+94+ipy/vl4ptkbrpNELPhlfcIZamFjqZMH//PBu9B1/HZzN/zcasfeo0QJhzK4k
pSxRECObjafT8YSNRsHZXGpCuwSBwRljc4t/59KiY0pGC57Q5Vd8u9FsPL3g59U8E/2FgphQ4NyCv6fL
txKUSTiT8YI7mWheTmNsnlmToaWCaVjhggvQ4dKI3PHgHSiH8/PbCfXzLTr5D0QKW86/kjo2V2FmnCRp
NA8ElvGNjB5lYFFT0wIxOrKmCMvbCjcQHnyweePOVGQYprKc7jIFLnXCIurm/RTSsQBFKlW8+XeJViAw
NSpGe34z4XxnxmY2q3KqQY2q/y74OjLX/OMae7l7szPqn7I6zApsInWocEk8mFx4QKxMUl8MmcwTERki
s2oLMlaiJtgU1BotSQGqDdBlIKROeDCtnV2fLBDlRiFYhJ2s1OYtJzL6bva6ZrAOp6AwOYWOirI8UccH
gXeKtOnAW5EIXR6tJPH7yLqjRKg4IwvaKagIs+AFOh78US1x6GR1a62lkxXh6mjd8h59oRYFyjW6MMYl
5IoaMlHSenN7e3/PQFxKnTRvidcZ6Nj/qEuplD9qq7iTY0HVnv6OhB0Ovz70Q2F7VfomZi/INuCLekhN
sHuB7hOnSZ07F/RQotBC2rvKezeJ7yzzvaTeR/Z+qcTrNMLUICv72BRUZdccgfWS0XaZaZGdqRf4SCeA
pVfkTupEYaikxnBl4jb3DkRWRjmh2x/cHb6VFaNpFKMTnK1B5bjgP6FaY1lU7ApLKiyms2/ZyxtXfV8s
D+82XHfwrpe+7eFiiPbwzL2W+ubHoN/AuSy14PDVKZzJyWtvctra+1GTLTZuN/t4s/xxGtQUXFdodVyp
JBUdUir1TbShSMHy4EuvzG6fhkPCa6ot6l81svL6r4yNGRlmc4WMUlwxUMrvqFlOYZbbzJRpvF3zs2TK
9Fmlhzc6s4Fk+nuWGBOzraYwl5pcxSyFNTIgphAcsck3rOQUCELrGOiYRchiuVxKkSsqqZDk6Bz7fywT
Se4Fc8UqMsq9YEhi/MX4syzs2WkLu3x42/7cIYxeSrsKvbrBAKXbzWl7PQsdvhQ/xvdifd/r6+XRevs0
PxH4YVNLOyrw6vjih0h7nLi9yNudwJ193BEiHySzb+3u+LpWjH7Imu4C9zJsvU1bZ+P21uSRQiZSrJL5
ZKp7+t9W95NS5hamzk+Wf8cRFVn1gLFr0MiwG6lghcktE6mRAsdPpqZnQ9f0J2xBvz6tBf1ZL80bsDs2
FK01NoxK5X5Y94nOQYJhSQgeVLs228fW70P934sOI0SPipWdoK1+jO3XavayKIwm1OSVxtYJ/ITa0MFG
fKwZ36MuP4TveTND3I5ft6xGmpc8Rrlm2vV5Od6/J/YiYAMJjxLxIZSqC7+GsA8vn+1Dm7t45N9X+Hye
0Piw3/x5xd0Ydwa3A/Pznc8f/w0AAP//mW8YCVcpAAA=
`,
	},

	"/static/ui/sign.ui~": {
		local:   "static/ui/sign.ui~",
		size:    10735,
		modtime: 1544925503,
		compressed: `
H4sIAAAAAAAC/+xaX3PjtBd976fQT08/YJM2CQvMTuIdlmUXZpiBgeXZcy3f2KKKZKTrtObTM3baTdo4
sWW7obvtWxvpSrp/zrnHluevr1eKrdE6afSCT8YXnKEWJpY6WfA/P7wbfcdfB2fz/41G7D1qtEAYsytJ
KUsUxMhm4+l0PGGjUXA2l5rQLkFgcMbY3OLfubTomJLRgid0+RXfbjQbTy/4eTXPRH+hICYUOLfg7+ny
rQRlEs5kvOBOJpqX0xibZ9ZkaKlgGla44AJ0uDQidzx4B8rh/Px2Qv18i07+A5HClvOvpI7NVZgZJ0ka
zQOBpX8jo0cZWNTUtECMjqwpwjJa4caEBx9s3rgzFRmGqSynu0yBS52wiLp5P4V0zEGRShVv/i6tFQhM
jYrRnt9MON+ZsZnNqpxqUKPq3wVfR+aaf1xjL3dvdkb9U1ZnswKbSB0qXBIPJhceJlYmqa8NmczTIjJE
ZtXWyFiJmmBTUGu0JAWoNoYuAyF1woNp7ez6ZIEoNwrBIuxkpTZvOZHRd7PXNYN1dgoKk1PoqCjLE3V8
0PBOkTYdeEsSocujlSR+37LuKBEqzsiCdgoqwCx4gY4Hf1RLHDpZ3Vpr6WQFuDpYt4yjr6lFgXKNLoxx
CbmihkyUsN5Eb+/3DMSl1EnzlnidgY79j7qUSvlbbRl3csyp2tPfobDD7te7fshtr0rf+OxlsnX4ot6k
xtk9R/eB08TOnQt6KFJoQe1d6b0bxXem+V5U70N7v1TkdRpiaqCVfdsUVCXXHIH1otF2mWmRnamX8ZFO
AEsvz53UicJQSY3hysRt4g5EVkY5odsf3B2+pRWjaRSjE5ytQeW44D+hWmNZVOwKSygsprNv2csbVX2f
LA/vNlx38K6Xvu3hYoj28Iy9lvzmh6DfwLksteDw1SmUyclrb3La2vtRky02ajf7GFn+OAVqCq6raXVc
qSQVHVIq9Y23oUjB8uBLr8xun4ZDwmuqLepfNbIy/FfGxowMs7lCRimuGCjld9QspzDLbWbKNN6u+Vki
ZfrM0sMLndlANP09S4yJ2ZZTmEtNrmKWwhoZEFMIjtjkG1ZiCgShdQx0zCJksVwupcgVlVBIcnSO/T+W
iST3grliFRnlXjAkMf5i/FkW9uy0hV0+vG1fdwijl9KuQq9uMEDpdlPaXs9Ch4Pih/heqO8bvl4arbdO
8yOBHza1tMMCr44vfgi0x4HbC7zdAdxZxx0B8kEw+9bujq5rheiHrOku5l6Crbdo6yzc3po8UshEilUy
n0x1T//b6n5SzNxC1PnR8u84oiKrHjB2BRoZdkMVrDC5ZSI1UuD4ydT0bOia/oQl6NenlaA/66V5A3ZH
hqK1xoZRydwPqz7ROUgwLAHBg2rXZvnY+j7U/150GCJ6VKjsZNrqZWy/VrOXRWE0oSavNLZO4CfUhg42
4mPN+B50+SH7npEZIjp+3bIaabXk8ZuXTjcwk2ndtUub65c2VNBMB30u7fv36l7E0EAORwniIRi0C+6H
kDUvn2VNm1g88u8+fD6baHwJ0fzZx10fdwa3A/Pznc8y/w0AAP//4H1c6O8pAAA=
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

	"/static/ui/tag.ui~": {
		local:   "static/ui/tag.ui~",
		size:    6751,
		modtime: 1544925387,
		compressed: `
H4sIAAAAAAAC/+xZUW/bNhB+z6/g+DootpWXYbBUoNvaDRj2su5ZOFFniQtNauTJjvfrB8ltbTeSKFpu
0AB7S6L7SN7dd9/dIes3T1vFdmidNDrhq/slZ6iFKaQuE/7Xh3fRD/xNerf+LorYe9RogbBge0kVKxUU
yB7u4/h+xaIovVtLTWg3IDC9Y2xt8Z9GWnRMyTzhJT1+z08XPdzHS77o7Ez+NwpiQoFzCX9Pjz9LUKbk
TBYJJyh5a8XYuramRksHpmGLCRegs40RjePpO1AO14tPBv32Fp38F3KFE+23pgDF0w+28ZrupS7MPquN
kySN5qnANhKR0VENFjX5DihwA42ibC8Lqnj6sFz6EY6sOWRtJrLjJdPeSocas0q25kUXZ/9NCmksbqKS
qjj+3KIVCKyMKtAuPhosziyO1qxjigYVdb8mfJebJ/75jGeMeHv2NZwJvdkFW0qdKdwQT1fLAIiVZRWK
IVMHInJDZLZTQcZK1ARH8u3QkhSgpgBdDULqkqdxr3V/skC0F2VgEc6y0pu3hsjoy+xdm8E+nIKDaShz
dGjpiboYBF6Q1Pfgz9KTCdACFf8S2PeSHBVnZEE7BV29JPyAjqc/dUcMPazvrJ10squ3vnqeGMZQqEWB
cocu+yhFvhPWi2Psnv29BvEodem/EZ9q0EX4SzdSqXDUSZuXY071vv5CwOYyyjX5VtIsRv3ZHfE/o74V
Rq1uwag+9/tdH3I7SDmPPgdBvCXU4+wzR5+Xja/bX83nWzWZLzpyvJwKnNyVQ/Tk904WXqbkR8PUh61A
yVLz1BHYIIE6n5Pi5RXI8XkpXGP/wD0jKDu7H19CtIJj/cr64C+a7OHUBtvX8G+zcVXgpkFfcfJXL5v8
3/TGvAV7Sj9aa2yWg+VfV5D8+9cgEp2DErN2X+Zp995ReNiqFL4y3SI6w5Qdp+0s6l45fcwSsBEeD3J5
JIvCaEJNQWmcnMDJEjc/+zMa9nC1+0aVy5rnQ/CZgblFcMImhO7LpCOByMq8IXTDRudmn4rEaIoKdIKz
HagGE/4rqh22UyRbxXwxcuPCf+WYEvjVYM76Na+nzdYFjzaM6sPXENBryn56qIbDFN9aPl/xOBS/2g0+
ZAH2jn/+Bf7Sx7OPpw/rxdm/gf4LAAD//136Pz5fGgAA
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

	"/static/ui/vault.ui~": {
		local:   "static/ui/vault.ui~",
		size:    25447,
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
