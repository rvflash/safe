function viewDefaultData() {
    return {
        error: "",
        last_upd: "",
        note: "",
        pass: "",
        strength: 0,
        url: "",
        user: "",
        warn: "",
    };
}

Vue.component('view-vault', {
    delimiters: ['${', '}'],
    data: function () {
        return viewDefaultData();
    },
    props: {
        tag: {
            type: String,
            required: true
        },
        vault: {
            type: String,
            required: true
        },
    },
    created: function () {
        this.select();
    },
    watch: {
        tag: function () {
            this.select();
        },
        vault: function () {
            this.select();
        },
    },
    computed: {
        progress: function () {
            if (this.strength === 0) {
                return '';
            }
            return 'w-' + (25 * this.strength);
        },
    },
    filters: {
        flag: function (s) {
            let comma = s.indexOf(',');
            if (comma < 0) {
                return s;
            }
           return s.substring(0, comma);
        },
    },
    methods: {
        copy: function (kind) {
            // Creates a temporary invisible input to contain the value to copy.
            let text = document.createElement('input');
            text.style.position = 'absolute';
            text.style.left = '-1000px';
            text.style.top = '-1000px';
            if (kind === 'name') {
                text.value = this.user;
            } else if  (kind === 'pass') {
                if (this.pass === "") {
                    // Retrieve the pass.
                    this.pwd();
                }
                text.value = this.pass;
            }
            document.body.appendChild(text);
            text.select();
            document.execCommand("copy");
            document.body.removeChild(text);
        },
        remove: function () {
            let ok = window.confirm("Are you sure you want to delete this vault?");
            if (!ok) {
                return;
            }
            let req = new XMLHttpRequest()
                , tag = encodeURIComponent(this.tag)
                , vault = encodeURIComponent(this.vault)
            ;

            // Deletes the vault.
            req.open('DELETE', `/tags/${tag}/vaults/${vault}`, false);
            req.send(null);
            if (req.status === 200) {
                this.$emit('vault', this.tag, '', 'create');
            } else {
                let res = JSON.parse(req.responseText);
                this.error = res.error;
            }
        },
        edit: function () {
            this.$emit('vault', this.tag, this.vault, 'update');
        },
        select: function() {
            let self = this
                , req = new XMLHttpRequest()
                , tag = encodeURIComponent(self.tag)
                , vault = encodeURIComponent(self.vault)
            ;
            self.resetData();

            // Retrieves the vault.
            req.open('GET', `/tags/${tag}/vaults/${vault}`, true);
            req.onload = function() {
                let res = JSON.parse(req.responseText);
                if (this.status === 200) {
                    self.user = res.user;
                    self.strength = res.strength;
                    self.url = res.url;
                    self.note = res.note;
                    self.last_upd = res.last_upd;
                    self.warn = res.safe;
                } else if (this.status === 404) {
                    self.$emit('vault', self.tag, '', 'create');
                } else {
                    self.error = res.error;
                }
            };
            req.send(null);
        },
        pwd: function() {
            let  req = new XMLHttpRequest()
                , tag = encodeURIComponent(this.tag)
                , vault = encodeURIComponent(this.vault)
            ;
            req.open('GET', `/tags/${tag}/vaults/${vault}?pwd`, false);
            req.send(null);
            if (req.status === 200) {
                let res = JSON.parse(req.responseText);
                this.pass = atob(res);
            } else {
                alert("Failed to retrieve the password");
            }
        },
        resetData: function () {
            Object.assign(this.$data, viewDefaultData());
        }
    },
    template: `<div class="vault">
    <div class="d-flex pt-2 pb-4">
        <h2 class="m-0">\${ vault }</h2>
        <div class="ml-auto">
            <button type="button" class="btn btn-primary btn-sm" v-on:click="edit">Edit</button>
            <button type="button" class="btn btn-outline-danger btn-sm" v-on:click="remove">Delete</button>
        </div>
    </div>
    <div class="alert alert-danger" role="alert" v-if="error.length">\${ error }</div>
    <div v-else>
        <dl class="row">
            <dt class="col-sm-3">Username</dt>
            <dd class="col-sm-9">
                \${ user }
                <button type="button" class="btn btn-link btn-sm" v-bind:disabled="user === ''" v-on:click="copy('name')">Copy</button>
            </dd>
            <dt class="col-sm-3">Password</dt>
            <dd class="col-sm-9">
                ••••••••
                <span class="badge badge-danger" v-if="warn" v-bind:title="warn">\${ warn | flag }</span>
                <button type="button" class="btn btn-link btn-sm" v-bind:disabled="user === ''" v-on:click="copy('pass')">Copy</button>
            </dd>
        </dl>
        <div class="progress">
            <div v-bind:class="['progress-bar bg-success', progress ]" role="progressbar" v-bind:aria-valuenow="strength" aria-valuemin="0" aria-valuemax="4"></div>
        </div>
        <hr v-if="url && note">
        <dl class="row">
            <dt class="col-sm-3" v-show="url">URL</dt>
            <dd class="col-sm-9" v-show="url"><a v-bind:href="encodeURI(url)" rel="noreferrer" target="_blank">\${ url }</a></dd>
            <dt class="col-sm-3" v-show="note">Note</dt>
            <dd class="col-sm-9" v-show="note">\${ note }</dd>
        </dl>
        <hr>
        <p><small class="text-muted">Last updated \${ last_upd }</small></p>
    </div>
</div>`
});